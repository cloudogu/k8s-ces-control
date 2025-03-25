package support

import (
	"context"
	"errors"
	"fmt"
	"slices"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const listVerb = "list"

type defaultResourceCollector struct {
	client          k8sClient
	discoveryClient discoveryInterface
}

type resourceCollector interface {
	Collect(ctx context.Context, labelSelector *metav1.LabelSelector, excludedGVKs []gvkMatcher) ([]*unstructured.Unstructured, error)
	listApiResourcesByLabelSelector(ctx context.Context, list *metav1.APIResourceList, selector labels.Selector, excludedGVKs []gvkMatcher) ([]*unstructured.Unstructured, []error)
	listByLabelSelector(ctx context.Context, resource metav1.APIResource, labelSelector labels.Selector, excludedGVKs []gvkMatcher) ([]*unstructured.Unstructured, error)
}

type gvkMatcher schema.GroupVersionKind

// Matches checks if the fields of the supplied schema.GroupVersionKind equal those of the gvkMatcher.
// Particular fields can be ignored by using the star-notation (*) in the matcher.
func (m gvkMatcher) Matches(gvk schema.GroupVersionKind) bool {
	return (gvk.Group == m.Group || m.Group == "*") &&
		(gvk.Version == m.Version || m.Version == "*") &&
		(gvk.Kind == m.Kind || m.Kind == "*")
}

// Collect fetches all resources in the cluster that match the supplied label selector,
// skipping resources matched by the excluded GVKs.
func (rc *defaultResourceCollector) Collect(ctx context.Context, labelSelector *metav1.LabelSelector, excludedGVKs []gvkMatcher) ([]*unstructured.Unstructured, error) {
	resourceKindLists, err := rc.discoveryClient.ServerPreferredResources()
	if err != nil {
		return nil, fmt.Errorf("failed to get resource kind lists from server: %w", err)
	}

	selector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to create selector from given label selector %s: %w", labelSelector, err)
	}

	var errs []error
	var resources []*unstructured.Unstructured
	for _, resourceKindList := range resourceKindLists {
		resourcesOfKind, listErrs := rc.listApiResourcesByLabelSelector(ctx, resourceKindList, selector, excludedGVKs)
		resources = append(resources, resourcesOfKind...)
		errs = append(errs, listErrs...)
	}

	if len(errs) != 0 {
		return nil, fmt.Errorf("failed to list api resources with label selector %q: %w", selector, errors.Join(errs...))
	}

	return resources, nil
}

func (rc *defaultResourceCollector) listApiResourcesByLabelSelector(ctx context.Context, list *metav1.APIResourceList, selector labels.Selector, excludedGVKs []gvkMatcher) ([]*unstructured.Unstructured, []error) {
	if len(list.APIResources) == 0 {
		return nil, nil
	}

	gv, err := schema.ParseGroupVersion(list.GroupVersion)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to list api resources with group version %q: %w", list.GroupVersion, err)}
	}

	var errs []error
	var resources []*unstructured.Unstructured
	for _, resource := range list.APIResources {
		if len(resource.Verbs) != 0 && slices.Contains(resource.Verbs, listVerb) {
			resource.Group = gv.Group
			resource.Version = gv.Version

			resourcesByLabelSelector, listErr := rc.listByLabelSelector(ctx, resource, selector, excludedGVKs)
			if listErr != nil {
				errs = append(errs, listErr)
			} else {
				resources = append(resources, resourcesByLabelSelector...)
			}
		}
	}

	return resources, errs
}

func (rc *defaultResourceCollector) listByLabelSelector(ctx context.Context, resource metav1.APIResource, labelSelector labels.Selector, excludedGVKs []gvkMatcher) ([]*unstructured.Unstructured, error) {
	logger := log.FromContext(ctx)

	gvk := groupVersionKind(resource)
	for _, gvkMatcher := range excludedGVKs {
		if gvkMatcher.Matches(gvk) {
			logger.Info(fmt.Sprintf("skipping resource %s as it is excluded", gvk))
			return nil, nil
		}
	}

	listOptions := client.ListOptions{LabelSelector: &client.MatchingLabelsSelector{Selector: labelSelector}}
	objectList := &unstructured.UnstructuredList{}
	objectList.SetGroupVersionKind(gvk)
	err := rc.client.List(ctx, objectList, &listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects in %s: %w", gvk, err)
	}

	return sliceToPointers(objectList.Items), nil
}

func groupVersionKind(resource metav1.APIResource) schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   resource.Group,
		Version: resource.Version,
		Kind:    resource.Kind,
	}
}

func sliceToPointers[T any](raw []T) []*T {
	pointers := make([]*T, len(raw))
	for i, t := range raw {
		pointers[i] = &t
	}
	return pointers
}
