package doguHealth

import (
	"github.com/cloudogu/k8s-dogu-operator/v2/api/ecoSystem"
)

type doguClient interface {
	ecoSystem.DoguInterface
}
