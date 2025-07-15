package supportArchive

import (
	"github.com/cloudogu/ces-control-api/generated/maintenance"
	"github.com/cloudogu/k8s-support-archive-lib/client/v1"
)

type supportArchiveClient interface {
	v1.SupportArchiveInterface
}

type supportArchiveCreateserver interface {
	maintenance.SupportArchive_CreateServer
}
