package util

import (
	"context"
	common "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

var testCtx = context.Background()

func TestNewDoguGetter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		doguVersionRegistryMock := newMockDoguVersionRegistry(t)
		localDoguDescriptorRepositoryMock := newMockLocalDoguDescriptorRepository(t)

		// when
		getter := NewDoguGetter(doguVersionRegistryMock, localDoguDescriptorRepositoryMock)

		// then
		require.NotNil(t, getter)
		assert.Equal(t, getter.doguRepository, localDoguDescriptorRepositoryMock)
		assert.Equal(t, getter.versionRegistry, doguVersionRegistryMock)
	})
}

func Test_doguGetter_GetCurrent(t *testing.T) {
	ldapVersion, err := core.ParseVersion("1.0.0")
	require.NoError(t, err)
	ldapDoguVersion := common.SimpleNameVersion{Name: "ldap", Version: ldapVersion}
	ldapDogu := &core.Dogu{Name: "ldap", Version: "1.0.0"}

	type args struct {
		ctx            context.Context
		simpleDoguName string
	}
	tests := []struct {
		name    string
		args    args
		want    *core.Dogu
		wantErr bool
		mockFn  func() (common.VersionRegistry, common.LocalDoguDescriptorRepository)
	}{
		{
			name:    "success",
			args:    args{ctx: testCtx, simpleDoguName: "ldap"},
			want:    ldapDogu,
			wantErr: false,
			mockFn: func() (common.VersionRegistry, common.LocalDoguDescriptorRepository) {
				versionRegistryMock := newMockDoguVersionRegistry(t)
				localDoguDescriptorRepositoryMock := newMockLocalDoguDescriptorRepository(t)

				versionRegistryMock.EXPECT().GetCurrent(testCtx, common.SimpleName("ldap")).Return(ldapDoguVersion, nil)
				localDoguDescriptorRepositoryMock.EXPECT().Get(testCtx, ldapDoguVersion).Return(ldapDogu, nil)

				return versionRegistryMock, localDoguDescriptorRepositoryMock
			},
		},
		{
			name:    "should return error on error getting current version",
			args:    args{ctx: testCtx, simpleDoguName: "ldap"},
			want:    nil,
			wantErr: true,
			mockFn: func() (common.VersionRegistry, common.LocalDoguDescriptorRepository) {
				versionRegistryMock := newMockDoguVersionRegistry(t)
				localDoguDescriptorRepositoryMock := newMockLocalDoguDescriptorRepository(t)

				versionRegistryMock.EXPECT().GetCurrent(testCtx, common.SimpleName("ldap")).Return(ldapDoguVersion, assert.AnError)

				return versionRegistryMock, localDoguDescriptorRepositoryMock
			},
		},
		{
			name:    "should return error on error getting dogu",
			args:    args{ctx: testCtx, simpleDoguName: "ldap"},
			want:    nil,
			wantErr: true,
			mockFn: func() (common.VersionRegistry, common.LocalDoguDescriptorRepository) {
				versionRegistryMock := newMockDoguVersionRegistry(t)
				localDoguDescriptorRepositoryMock := newMockLocalDoguDescriptorRepository(t)

				versionRegistryMock.EXPECT().GetCurrent(testCtx, common.SimpleName("ldap")).Return(ldapDoguVersion, nil)
				localDoguDescriptorRepositoryMock.EXPECT().Get(testCtx, ldapDoguVersion).Return(nil, assert.AnError)

				return versionRegistryMock, localDoguDescriptorRepositoryMock
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			versionRegistryMock, localDoguDescriptorRepositoryMock := tt.mockFn()

			r := &doguGetter{
				versionRegistry: versionRegistryMock,
				doguRepository:  localDoguDescriptorRepositoryMock,
			}
			got, err := r.GetCurrent(tt.args.ctx, tt.args.simpleDoguName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCurrent() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_doguGetter_GetCurrentOfAll(t *testing.T) {
	ldapVersion, err := core.ParseVersion("1.0.0")
	require.NoError(t, err)
	ldapDoguVersion := common.SimpleNameVersion{Name: "ldap", Version: ldapVersion}
	ldapDogu := &core.Dogu{Name: "ldap", Version: "1.0.0"}

	casVersion, err := core.ParseVersion("1.0.0")
	require.NoError(t, err)
	casDoguVersion := common.SimpleNameVersion{Name: "cas", Version: casVersion}
	casDogu := &core.Dogu{Name: "cas", Version: "1.0.0"}

	doguVersions := []common.SimpleNameVersion{ldapDoguVersion, casDoguVersion}
	dogusMap := map[common.SimpleNameVersion]*core.Dogu{ldapDoguVersion: ldapDogu, casDoguVersion: casDogu}
	dogus := []*core.Dogu{ldapDogu, casDogu}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*core.Dogu
		wantErr bool
		mockFn  func() (common.VersionRegistry, common.LocalDoguDescriptorRepository)
	}{
		{
			name:    "success",
			args:    args{ctx: testCtx},
			want:    dogus,
			wantErr: false,
			mockFn: func() (common.VersionRegistry, common.LocalDoguDescriptorRepository) {
				versionRegistryMock := newMockDoguVersionRegistry(t)
				localDoguDescriptorRepositoryMock := newMockLocalDoguDescriptorRepository(t)

				versionRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return(doguVersions, nil)
				localDoguDescriptorRepositoryMock.EXPECT().GetAll(testCtx, doguVersions).Return(dogusMap, nil)

				return versionRegistryMock, localDoguDescriptorRepositoryMock
			},
		},
		{
			name:    "should return error on error getting all current versions",
			args:    args{ctx: testCtx},
			want:    nil,
			wantErr: true,
			mockFn: func() (common.VersionRegistry, common.LocalDoguDescriptorRepository) {
				versionRegistryMock := newMockDoguVersionRegistry(t)
				localDoguDescriptorRepositoryMock := newMockLocalDoguDescriptorRepository(t)

				versionRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return(doguVersions, assert.AnError)

				return versionRegistryMock, localDoguDescriptorRepositoryMock
			},
		},
		{
			name:    "should return error on error getting all dogus",
			args:    args{ctx: testCtx},
			want:    nil,
			wantErr: true,
			mockFn: func() (common.VersionRegistry, common.LocalDoguDescriptorRepository) {
				versionRegistryMock := newMockDoguVersionRegistry(t)
				localDoguDescriptorRepositoryMock := newMockLocalDoguDescriptorRepository(t)

				versionRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return(doguVersions, nil)
				localDoguDescriptorRepositoryMock.EXPECT().GetAll(testCtx, doguVersions).Return(dogusMap, assert.AnError)

				return versionRegistryMock, localDoguDescriptorRepositoryMock
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			versionRegistryMock, localDoguDescriptorRepositoryMock := tt.mockFn()

			r := &doguGetter{
				versionRegistry: versionRegistryMock,
				doguRepository:  localDoguDescriptorRepositoryMock,
			}
			got, err := r.GetCurrentOfAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrentOfAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.ElementsMatch(t, got, tt.want) {
				t.Errorf("GetCurrentOfAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}
