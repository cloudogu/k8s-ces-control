package supportArchive

import (
	"net/http"

	"github.com/cloudogu/ces-control-api/generated/maintenance"
	"github.com/cloudogu/k8s-support-archive-lib/client/v1"
)

type supportArchiveClient interface {
	v1.SupportArchiveInterface
}

//nolint:unused
//goland:noinspection GoUnusedType
type supportArchiveDownloadServer interface {
	maintenance.SupportArchive_DownloadSupportArchiveServer
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
