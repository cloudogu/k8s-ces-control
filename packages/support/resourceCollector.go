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

type resourceCollector struct {
	client          k8sClient
	discoveryClient discoveryInterface
}

type gvkMatcher schema.GroupVersionKind

func (m gvkMatcher) Matches(gvk schema.GroupVersionKind) bool {
	return (gvk.Group == m.Group || m.Group == "*") &&
		(gvk.Version == m.Version || m.Version == "*") &&
		(gvk.Kind == m.Kind || m.Kind == "*")
}

func (rc *resourceCollector) Collect(ctx context.Context, labelSelector *metav1.LabelSelector, excludedGVKs []gvkMatcher) ([]*unstructured.Unstructured, error) {
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
		return nil, fmt.Errorf("failed to delete api resources with label selector %q: %w", selector, errors.Join(errs...))
	}

	return resources, nil
}

func (rc *resourceCollector) listApiResourcesByLabelSelector(ctx context.Context, list *metav1.APIResourceList, selector labels.Selector, excludedGVKs []gvkMatcher) ([]*unstructured.Unstructured, []error) {
	if len(list.APIResources) == 0 {
		return nil, nil
	}

	gv, err := schema.ParseGroupVersion(list.GroupVersion)
	if err != nil {
		log.FromContext(ctx).Error(err, fmt.Sprintf("failed to delete api resources with group version: %s", list.GroupVersion))
		return nil, nil
	}

	var errs []error
	var resources []*unstructured.Unstructured
	for _, resource := range list.APIResources {
		if len(resource.Verbs) != 0 && slices.Contains(resource.Verbs, listVerb) {
			resource.Group = gv.Group
			resource.Version = gv.Version

			resourcesByLabelSelector, listErr := rc.listByLabelSelector(ctx, resource, selector, excludedGVKs)
			resources = append(resources, resourcesByLabelSelector...)
			errs = append(errs, listErr)
		}
	}

	return resources, errs
}

func (rc *resourceCollector) listByLabelSelector(ctx context.Context, resource metav1.APIResource, labelSelector labels.Selector, excludedGVKs []gvkMatcher) ([]*unstructured.Unstructured, error) {
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
