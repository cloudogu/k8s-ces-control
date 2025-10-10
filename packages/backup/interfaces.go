package backup

import (
	backupClientV1 "github.com/cloudogu/k8s-backup-lib/api/ecosystem"
)

type backupInterface interface {
	backupClientV1.BackupInterface
}

type restoreInterface interface {
	backupClientV1.RestoreInterface
}
