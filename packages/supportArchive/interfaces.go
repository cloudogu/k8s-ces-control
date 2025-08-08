package supportArchive

import (
	"github.com/cloudogu/ces-control-api/generated/maintenance"
	"github.com/cloudogu/k8s-support-archive-lib/client/v1"
	"net/http"
)

type supportArchiveClient interface {
	v1.SupportArchiveInterface
}

//nolint:unused
//goland:noinspection GoUnusedType
type supportArchiveCreateserver interface {
	maintenance.SupportArchive_CreateServer
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
