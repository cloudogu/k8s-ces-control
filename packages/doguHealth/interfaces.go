package doguHealth

import (
	"github.com/cloudogu/k8s-dogu-lib/v2/client"
)

type doguClient interface {
	client.DoguInterface
}
