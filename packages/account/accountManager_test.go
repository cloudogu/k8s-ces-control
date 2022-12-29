package account

// TODO :)
// import (
// 	"errors"
// 	"github.com/cloudogu/cesapp-lib/core"
// 	"github.com/cloudogu/cesapp-lib/keys"
// 	"github.com/cloudogu/k8s-ces-control/packages/account/mocks"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"
// 	"testing"
// )
//
// const (
// 	testdataDir = "../testdata"
// )
//
// var (
// 	testError = errors.New("test error")
// )
//
// func TestServiceAccountManager_Create(t *testing.T) {
// 	t.Run("should create an service account successfully", func(t *testing.T) {
// 		mockHostConfiguration := &mocks.ConfigurationRegistry{}
// 		mockHostConfiguration.On("Set", "dummy/username", mock.Anything).Return(nil)
// 		mockHostConfiguration.On("Set", "dummy/password", mock.Anything).Return(nil)
// 		mockKeyProvider := &mocks.KeyProvider{}
// 		testKeyPair := generateTestCertificate()
// 		mockKeyProvider.On("FromPrivateKeyPath", ServerCertificatePath).Return(testKeyPair, nil)
// 		testServiceAccountManager := ServiceAccountManager{
// 			hostConfiguration: mockHostConfiguration,
// 			keyProvider:       mockKeyProvider,
// 			serviceName:       "dummy",
// 		}
//
// 		testAccountData, err := testServiceAccountManager.Create()
//
// 		assert.NoError(t, err)
// 		assert.NotEqual(t, ServiceAccountData{}, testAccountData)
// 		assert.Contains(t, testAccountData.Username, "dummy_")
// 		assert.Len(t, testAccountData.Username, len("dummy_")+randUsernameSuffixLength) // "dummy_" + randUsernameSuffixLength
// 		assert.Len(t, testAccountData.Password, passwordLength)
// 		mockKeyProvider.AssertExpectations(t)
// 		mockHostConfiguration.AssertExpectations(t)
// 	})
// 	t.Run("should return an error if the password encryption fails", func(t *testing.T) {
// 		mockHostConfiguration := &mocks.ConfigurationRegistry{}
// 		mockKeyProvider := &mocks.KeyProvider{}
// 		mockKeyProvider.On("FromPrivateKeyPath", ServerCertificatePath).Return(nil, testError)
// 		testServiceAccountManager := ServiceAccountManager{
// 			hostConfiguration: mockHostConfiguration,
// 			keyProvider:       mockKeyProvider,
// 			serviceName:       "dummy",
// 		}
//
// 		testAccountData, err := testServiceAccountManager.Create()
//
// 		assert.Error(t, err)
// 		assert.Equal(t, ServiceAccountData{}, testAccountData)
// 		mockKeyProvider.AssertExpectations(t)
// 		mockHostConfiguration.AssertExpectations(t)
// 	})
// 	t.Run("should return an error if setting the username fails", func(t *testing.T) {
// 		mockHostConfiguration := &mocks.ConfigurationRegistry{}
// 		mockHostConfiguration.On("Set", "dummy/username", mock.Anything).Return(testError)
// 		mockKeyProvider := &mocks.KeyProvider{}
// 		testKeyPair := generateTestCertificate()
// 		mockKeyProvider.On("FromPrivateKeyPath", ServerCertificatePath).Return(testKeyPair, nil)
// 		testServiceAccountManager := ServiceAccountManager{
// 			hostConfiguration: mockHostConfiguration,
// 			keyProvider:       mockKeyProvider,
// 			serviceName:       "dummy",
// 		}
//
// 		testAccountData, err := testServiceAccountManager.Create()
//
// 		assert.Error(t, err)
// 		assert.Equal(t, ServiceAccountData{}, testAccountData)
// 		mockKeyProvider.AssertExpectations(t)
// 		mockHostConfiguration.AssertExpectations(t)
// 	})
// 	t.Run("should return an error if setting the password fails", func(t *testing.T) {
// 		mockHostConfiguration := &mocks.ConfigurationRegistry{}
// 		mockHostConfiguration.On("Set", "dummy/username", mock.Anything).Return(nil)
// 		mockHostConfiguration.On("Set", "dummy/password", mock.Anything).Return(testError)
// 		mockKeyProvider := &mocks.KeyProvider{}
// 		testKeyPair := generateTestCertificate()
// 		mockKeyProvider.On("FromPrivateKeyPath", ServerCertificatePath).Return(testKeyPair, nil)
// 		testServiceAccountManager := ServiceAccountManager{
// 			hostConfiguration: mockHostConfiguration,
// 			keyProvider:       mockKeyProvider,
// 			serviceName:       "dummy",
// 		}
//
// 		testAccountData, err := testServiceAccountManager.Create()
//
// 		assert.Error(t, err)
// 		assert.Equal(t, ServiceAccountData{}, testAccountData)
// 		mockKeyProvider.AssertExpectations(t)
// 		mockHostConfiguration.AssertExpectations(t)
// 	})
// }
//
// func Test_generateUsernamePassword(t *testing.T) {
// 	t.Run("should generate username and password with correct length", func(t *testing.T) {
// 		generatedData := generateUsernamePassword("testservice")
//
// 		assert.NotEqual(t, ServiceAccountData{}, generatedData)
// 		assert.Contains(t, generatedData.Username, "testservice_")
// 		assert.Len(t, generatedData.Username, len("testservice_")+randUsernameSuffixLength) // "dummy_" + randUsernameSuffixLength
// 		assert.Len(t, generatedData.Password, passwordLength)
// 	})
// }
//
// func Test_encrypt(t *testing.T) {
// 	t.Run("should encrypt the given string", func(t *testing.T) {
// 		mockKeyProvider := &mocks.KeyProvider{}
// 		testKeyPair := generateTestCertificate()
// 		mockKeyProvider.On("FromPrivateKeyPath", "path/to/server.key").Return(testKeyPair, nil)
//
// 		encryptedString, err := encrypt("test-password", "path/to/server.key", mockKeyProvider)
//
// 		assert.NoError(t, err)
// 		assert.NotEmpty(t, encryptedString)
// 		decryptedString, _ := testKeyPair.Private().Decrypt(encryptedString)
// 		assert.Equal(t, "test-password", decryptedString)
// 	})
// 	t.Run("should return an error if the certificate could not be loaded", func(t *testing.T) {
// 		mockKeyProvider := &mocks.KeyProvider{}
// 		mockKeyProvider.On("FromPrivateKeyPath", "path/to/server.key").Return(nil, testError)
//
// 		encryptedString, err := encrypt("test-password", "path/to/server.key", mockKeyProvider)
//
// 		assert.Error(t, err)
// 		assert.Empty(t, encryptedString)
// 	})
// 	t.Run("should return an error if the encryption fails", func(t *testing.T) {
// 		tmpKeyProvider, _ := keys.NewKeyProvider("oaesp")
//
// 		encryptedString, err := encrypt("test-password", testdataDir+"/private_wrong.pub", tmpKeyProvider)
//
// 		assert.Error(t, err)
// 		assert.Empty(t, encryptedString)
// 	})
// }
//
// // func Test_getKeyProvider(t *testing.T) {
// // 	t.Run("should create key provider successfully", func(t *testing.T) {
// // 		ctx := &config.Context{Configuration: &config.Configuration{Keys: config.Keys{Type: "oaesp"}}}
// // 		provider, err := getKeyProvider(ctx)
// //
// // 		require.NoError(t, err)
// // 		assert.NotNil(t, provider)
// // 	})
// // 	t.Run("should return an error if the creation fails", func(t *testing.T) {
// // 		ctx := &config.Context{Configuration: &config.Configuration{Keys: config.Keys{Type: "der-gary"}}}
// // 		provider, err := getKeyProvider(ctx)
// //
// // 		require.Error(t, err)
// // 		assert.Nil(t, provider)
// // 	})
// // }
//
// // func TestNewServiceAccountManager(t *testing.T) {
// // 	t.Run("should create ServiceAccountManager successfully", func(t *testing.T) {
// // 		testContext := &config.Context{
// // 			Configuration: &config.Configuration{
// // 				Keys: config.Keys{Type: "oaesp"},
// // 				Registry: core.Registry{
// // 					Type:      "etcd",
// // 					Endpoints: []string{"my-endpoint"},
// // 				},
// // 			},
// // 		}
// // 		accountManager, err := NewServiceAccountManager("dummy", testContext)
// //
// // 		require.NoError(t, err)
// // 		assert.NotNil(t, accountManager)
// // 		assert.Equal(t, "dummy", accountManager.serviceName)
// // 		assert.NotNil(t, accountManager.keyProvider)
// // 		assert.NotNil(t, accountManager.hostConfiguration)
// // 	})
// // 	t.Run("should return an error if creating the key provider fails", func(t *testing.T) {
// // 		testContext := &config.Context{
// // 			Configuration: &config.Configuration{
// // 				Keys: config.Keys{Type: "der-gary"},
// // 				Registry: core.Registry{
// // 					Type:      "etcd",
// // 					Endpoints: []string{"my-endpoint"},
// // 				},
// // 			},
// // 		}
// // 		accountManager, err := NewServiceAccountManager("dummy", testContext)
// //
// // 		require.Error(t, err)
// // 		assert.Contains(t, err.Error(), "failed to create key provider")
// // 		assert.Empty(t, accountManager)
// // 	})
// // 	t.Run("should return an error if creating the host configuration fails", func(t *testing.T) {
// // 		testContext := &config.Context{
// // 			Configuration: &config.Configuration{
// // 				Keys: config.Keys{Type: "oaesp"},
// // 				Registry: core.Registry{
// // 					Type:      "redis",
// // 					Endpoints: []string{"my-endpoint"},
// // 				},
// // 			},
// // 		}
// // 		accountManager, err := NewServiceAccountManager("dummy", testContext)
// //
// // 		require.Error(t, err)
// // 		assert.Contains(t, err.Error(), "failed to create registry")
// // 		assert.Empty(t, accountManager)
// // 	})
// // }
//
// // func generateTestCertificate() *keys.KeyPair {
// // 	tmpKeyProvider, _ := keys.NewKeyProvider(config.Keys{Type: "oaesp"})
// // 	testKeyPair, _ := tmpKeyProvider.Generate()
// // 	return testKeyPair
// // }
//
// func TestServiceAccountData_String(t *testing.T) {
// 	data := ServiceAccountData{
// 		Username: "dummyusr",
// 		Password: "dummypwd",
// 	}
//
// 	result := data.String()
//
// 	assert.Equal(t, "username:dummyusr\npassword:dummypwd", result)
// }
//
// func TestServiceAccountManager_SetServiceName(t *testing.T) {
// 	accountManager := ServiceAccountManager{}
//
// 	result := accountManager.SetServiceName("dummy")
//
// 	assert.Equal(t, "dummy", result.serviceName)
// }
//
// func TestServiceAccountManager_SetKeyProvider(t *testing.T) {
// 	accountManager := ServiceAccountManager{}
// 	mockProvider := &mocks.KeyProvider{}
//
// 	result := accountManager.SetKeyProvider(mockProvider)
//
// 	assert.Equal(t, mockProvider, result.keyProvider)
// }
//
// func TestServiceAccountManager_SetHostConfiguration(t *testing.T) {
// 	accountManager := ServiceAccountManager{}
// 	mockConfigRegistsy := &mocks.ConfigurationRegistry{}
//
// 	result := accountManager.SetHostConfiguration(mockConfigRegistsy)
//
// 	assert.Equal(t, mockConfigRegistsy, result.hostConfiguration)
// }
//
// // func Test_decrypt(t *testing.T) {
// // 	mockProvider := &mocks.KeyProvider{}
// // 	testKeyPair := generateTestCertificate()
// // 	mockProvider.On("FromPrivateKeyPath", "/etc/ces/cesappd/server.key").Return(testKeyPair, nil)
// // 	encryptedString, _ := encrypt("test-password", "/etc/ces/cesappd/server.key", mockProvider)
// //
// // 	t.Run("should run", func(t *testing.T) {
// // 		mockProvider.On("FromPrivateKeyPath", "/etc/ces/cesappd/server.key").Return(testKeyPair, nil)
// // 		decryptedString, err := decrypt(encryptedString, "/etc/ces/cesappd/server.key", mockProvider)
// // 		require.NoError(t, err)
// // 		assert.Contains(t, "test-password", decryptedString)
// // 	})
// //
// // 	t.Run("should return error after faulty keypair", func(t *testing.T) {
// // 		mockProvider := &mocks.KeyProvider{} // redeclare
// // 		mockProvider.On("FromPrivateKeyPath", "/etc/ces/cesappd/server.key").Return(nil, errors.New("oh no"))
// // 		_, err := decrypt(encryptedString, "/etc/ces/cesappd/server.key", mockProvider)
// // 		require.Error(t, err)
// // 		assert.Contains(t, err.Error(), "oh no")
// // 	})
// // }
//
// func Test_GetHostConfiguration(t *testing.T) {
// 	mockConfigRegistsy := &mocks.ConfigurationRegistry{}
// 	accountManager := ServiceAccountManager{hostConfiguration: mockConfigRegistsy}
// 	assert.Equal(t, mockConfigRegistsy, accountManager.GetHostConfiguration())
// }
//
// func Test_GetServiceAccountData(t *testing.T) {
// 	t.Run("should return error at get username", func(t *testing.T) {
// 		mockConfigRegistsy := &mocks.ConfigurationRegistry{}
// 		accountManager := ServiceAccountManager{serviceName: "service", hostConfiguration: mockConfigRegistsy}
// 		mockConfigRegistsy.On("Get", "service/username").Return("", errors.New("oh no"))
// 		_, err := accountManager.GetServiceAccountData()
// 		require.Error(t, err)
// 	})
//
// 	t.Run("should return error at get decrypted pw", func(t *testing.T) {
// 		mockConfigRegistsy := &mocks.ConfigurationRegistry{}
// 		accountManager := ServiceAccountManager{serviceName: "service", hostConfiguration: mockConfigRegistsy}
// 		mockConfigRegistsy.On("Get", "service/username").Return("jeff", nil)
// 		mockConfigRegistsy.On("Get", "service/password").Return("", errors.New("oh no"))
// 		_, err := accountManager.GetServiceAccountData()
// 		require.Error(t, err)
// 		assert.Contains(t, err.Error(), "oh no")
// 	})
//
// 	// t.Run("should return decryptionError", func(t *testing.T) {
// 	// 	mockKeyProvider := &mocks.KeyProvider{}
// 	// 	testKeyPair := generateTestCertificate()
// 	// 	mockKeyProvider.On("FromPrivateKeyPath", "/etc/ces/cesappd/server.key").Return(testKeyPair, nil)
// 	//
// 	// 	mockConfigRegistsy := &mocks.ConfigurationRegistry{}
// 	// 	mockConfigRegistsy.On("Get", "service/username").Return("jeff", nil)
// 	// 	mockConfigRegistsy.On("Get", "service/password").Return("worngstringthatisnotencrypted", nil)
// 	//
// 	// 	accountManager := ServiceAccountManager{serviceName: "service", hostConfiguration: mockConfigRegistsy, keyProvider: mockKeyProvider}
// 	//
// 	// 	_, err := accountManager.GetServiceAccountData()
// 	// 	require.Error(t, err)
// 	// 	assert.Contains(t, err.Error(), "illegal base64 data")
// 	// })
// 	//
// 	// t.Run("should return serviceAccountData", func(t *testing.T) {
// 	// 	mockKeyProvider := &mocks.KeyProvider{}
// 	// 	testKeyPair := generateTestCertificate()
// 	// 	mockKeyProvider.On("FromPrivateKeyPath", "/etc/ces/cesappd/server.key").Return(testKeyPair, nil)
// 	//
// 	// 	mockConfigRegistsy := &mocks.ConfigurationRegistry{}
// 	// 	mockConfigRegistsy.On("Get", "service/username").Return("jeff", nil)
// 	// 	encryptedString, _ := encrypt("test-password", "/etc/ces/cesappd/server.key", mockKeyProvider)
// 	// 	mockConfigRegistsy.On("Get", "service/password").Return(encryptedString, nil)
// 	//
// 	// 	accountManager := ServiceAccountManager{serviceName: "service", hostConfiguration: mockConfigRegistsy, keyProvider: mockKeyProvider}
// 	//
// 	// 	serviceAccountData, err := accountManager.GetServiceAccountData()
// 	// 	require.NoError(t, err)
// 	// 	assert.NotEmpty(t, serviceAccountData)
// 	// 	assert.Equal(t, ServiceAccountData{Username: "jeff", Password: "test-password"}, serviceAccountData)
// 	// })
// }
