package util

import (
	"github.com/cloudogu/ces-commons-lib/dogu"
)

type doguVersionRegistry interface {
	dogu.VersionRegistry
}

type localDoguDescriptorRepository interface {
	dogu.LocalDoguDescriptorRepository
}
