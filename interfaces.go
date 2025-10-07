package main

import (
	backupClientV1 "github.com/cloudogu/k8s-backup-lib/api/ecosystem"
	"github.com/cloudogu/k8s-ces-control/packages/doguAdministration"
	debugClientV1 "github.com/cloudogu/k8s-debug-mode-cr-lib/pkg/client/v1"
	ecoSystemV2 "github.com/cloudogu/k8s-dogu-lib/v2/client"
	supClientV1 "github.com/cloudogu/k8s-support-archive-lib/client/v1"
	"k8s.io/client-go/kubernetes"
	appsV1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	coreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

//nolint:unused
//goland:noinspection GoUnusedType
type configMapInterface interface {
	coreV1.ConfigMapInterface
}

//nolint:unused
//goland:noinspection GoUnusedType
type coreV1Interface interface {
	coreV1.CoreV1Interface
}

//nolint:unused
//goland:noinspection GoUnusedType
type appsV1Interface interface {
	appsV1.AppsV1Interface
}

//nolint:unused
//goland:noinspection GoUnusedType
type clusterClient interface {
	ecoSystemV2.EcoSystemV2Interface
	doguAdministration.BlueprintLister
	kubernetes.Interface
	supClientV1.SupportArchiveV1Interface
	debugClientV1.DebugModeV1Interface
	backupClientV1.BackupsGetter
	backupClientV1.RestoresGetter
	backupClientV1.BackupSchedulesGetter
}
