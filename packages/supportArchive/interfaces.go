package supportArchive

import (
	"github.com/cloudogu/k8s-support-archive-lib/client/v1"
)

type supportArchiveClient interface {
	v1.SupportArchiveInterface
}
