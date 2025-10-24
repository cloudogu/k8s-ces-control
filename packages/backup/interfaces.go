package backup

import (
	"context"

	backupClientV1 "github.com/cloudogu/k8s-backup-lib/api/ecosystem"
	backupV1 "github.com/cloudogu/k8s-backup-lib/api/v1"
	v2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	componentV1 "github.com/cloudogu/k8s-component-lib/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type backupInterface interface {
	backupClientV1.BackupInterface
}

type restoreInterface interface {
	backupClientV1.RestoreInterface
}

type blueprintLister interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v2.BlueprintList, error)
}

type backupScheduleClient interface {
	// Create takes the representation of a backup schedule and creates it.  Returns the server's representation of the backup schedule, and an error, if there is any.
	Create(ctx context.Context, backupSchedule *backupV1.BackupSchedule, opts metav1.CreateOptions) (*backupV1.BackupSchedule, error)

	// Update takes the representation of a backup schedule and updates it. Returns the server's representation of the backup schedule, and an error, if there is any.
	Update(ctx context.Context, backupSchedule *backupV1.BackupSchedule, opts metav1.UpdateOptions) (*backupV1.BackupSchedule, error)

	// Get takes name of the backup schedule, and returns the corresponding backup schedule object, and an error if there is any.
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*backupV1.BackupSchedule, error)
}

type componentClient interface {
	// Get takes name of the component, and returns the corresponding component object, and an error if there is any.
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*componentV1.Component, error)
}
