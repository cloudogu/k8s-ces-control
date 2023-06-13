package account

import (
	"context"
	"crypto/rsa"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/cesapp-lib/keys"
	"github.com/cloudogu/k8s-ces-control/packages/ssl"
)

var (
	testPrivateKeyString = "p12o3i4u0xc98v7"
	testPrivateKeyBytes  = []byte(testPrivateKeyString)
	testCtx              = context.Background()
)

func Test_generateUsernamePassword(t *testing.T) {
	t.Run("should generate username and password with correct length", func(t *testing.T) {
		generatedData := generateUsernamePassword("testservice")

		assert.NotEqual(t, ServiceAccountData{}, generatedData)
		assert.Contains(t, generatedData.Username, "testservice_")
		assert.Len(t, generatedData.Username, len("testservice_")+randUsernameSuffixLength)
		assert.Len(t, generatedData.Password, passwordLength)
	})
}

func Test_encrypt(t *testing.T) {
	t.Run("should encrypt the given string", func(t *testing.T) {
		mockKeyProvider := newMockKeyProvider(t)
		testKeyPair := generateTestKeyPair(t)
		mockKeyProvider.EXPECT().FromPrivateKey(testPrivateKeyBytes).Return(testKeyPair, nil)

		encryptedString, err := encrypt("test-password", testPrivateKeyString, mockKeyProvider)

		assert.NoError(t, err)
		assert.NotEmpty(t, encryptedString)
		decryptedString, _ := testKeyPair.Private().Decrypt(encryptedString)
		assert.Equal(t, "test-password", decryptedString)
	})
	t.Run("should return an error if the certificate could not be loaded", func(t *testing.T) {
		pubKeyBytes, _ := generateTestKeyPair(t).Public().AsBytes()
		mockKeyProvider := newMockKeyProvider(t)
		mockKeyProvider.EXPECT().FromPrivateKey(pubKeyBytes).Return(nil, assert.AnError)

		encryptedString, err := encrypt("test-password", string(pubKeyBytes), mockKeyProvider)

		assert.Error(t, err)
		assert.Empty(t, encryptedString)
	})
	t.Run("should return an error if the encryption fails", func(t *testing.T) {
		pubKeyBytes, _ := generateTestKeyPair(t).Public().AsBytes()
		mockKeyProvider := newMockKeyProvider(t)
		mockKeyProvider.EXPECT().FromPrivateKey(pubKeyBytes).Return(nil, assert.AnError)

		encryptedString, err := encrypt("test-password", string(pubKeyBytes), mockKeyProvider)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, encryptedString)
	})
}

func Test_decrypt(t *testing.T) {
	mockProvider := newMockKeyProvider(t)
	testKeyPair := generateTestKeyPair(t)
	mockProvider.EXPECT().FromPrivateKey(testPrivateKeyBytes).Return(testKeyPair, nil)
	encryptedString, _ := encrypt("test-password", testPrivateKeyString, mockProvider)

	t.Run("should run", func(t *testing.T) {
		mockProvider.EXPECT().FromPrivateKey(testPrivateKeyBytes).Return(testKeyPair, nil)
		decryptedString, err := decrypt(encryptedString, testPrivateKeyString, mockProvider)
		require.NoError(t, err)
		assert.Contains(t, "test-password", decryptedString)
	})

	t.Run("should return error after faulty keypair", func(t *testing.T) {
		mockProvider := newMockKeyProvider(t) // redeclare
		mockProvider.EXPECT().FromPrivateKey(testPrivateKeyBytes).Return(nil, assert.AnError)
		_, err := decrypt(encryptedString, testPrivateKeyString, mockProvider)
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestNewServiceAccountManager(t *testing.T) {
	t.Run("should fail to fetch key provider from registry", func(t *testing.T) {
		// given
		globalConfigMock := newMockRegistryContext(t)
		globalConfigMock.EXPECT().Get("key_provider").Return("", assert.AnError)
		registryMock := newMockConfigRegistry(t)
		registryMock.EXPECT().GlobalConfig().Return(globalConfigMock)

		// when
		actual, err := NewServiceAccountManager("test-service", registryMock)

		// then
		require.Error(t, err)
		assert.Empty(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to create ServiceAccountManager: failed to get key provider from global config")
	})
	t.Run("should fail to fetch key provider from registry", func(t *testing.T) {
		// given
		globalConfigMock := newMockRegistryContext(t)
		globalConfigMock.EXPECT().Get("key_provider").Return("abcde", nil)
		registryMock := newMockConfigRegistry(t)
		registryMock.EXPECT().GlobalConfig().Return(globalConfigMock)

		// when
		actual, err := NewServiceAccountManager("test-service", registryMock)

		// then
		require.Error(t, err)
		assert.Empty(t, actual)
		assert.ErrorContains(t, err, "failed to create ServiceAccountManager: failed to create key provider: could not find provider from type abcde")
	})
	t.Run("should create service account manager", func(t *testing.T) {
		// given
		globalConfigMock := newMockRegistryContext(t)
		globalConfigMock.EXPECT().Get("key_provider").Return("pkcs1v15", nil)
		hostConfigMock := newMockRegistryContext(t)
		registryMock := newMockConfigRegistry(t)
		registryMock.EXPECT().GlobalConfig().Return(globalConfigMock)
		registryMock.EXPECT().HostConfig("k8s-ces-control").Return(hostConfigMock)

		// when
		actual, err := NewServiceAccountManager("test-service", registryMock)

		// then
		require.NoError(t, err)
		assert.Equal(t, "test-service", actual.serviceName)
		assert.NotNil(t, actual.keyProvider)
		assert.Equal(t, hostConfigMock, actual.hostConfiguration)
	})
}

func TestServiceAccountManager_SetServiceName(t *testing.T) {
	t.Run("should set service name", func(t *testing.T) {
		// given
		sut := &ServiceAccountManager{serviceName: "some-name"}

		// when
		sut.SetServiceName("different-name")

		// then
		assert.Equal(t, "different-name", sut.serviceName)
	})
}

func TestServiceAccountManager_SetKeyProvider(t *testing.T) {
	t.Run("should set key provider", func(t *testing.T) {
		// given
		keyProviderMock := newMockKeyProvider(t)
		sut := &ServiceAccountManager{keyProvider: nil}

		// when
		sut.SetKeyProvider(keyProviderMock)

		// then
		assert.Equal(t, keyProviderMock, sut.keyProvider)
	})
}

func TestServiceAccountManager_SetHostConfiguration(t *testing.T) {
	t.Run("should set host config", func(t *testing.T) {
		// given
		hostConfigMock := newMockRegistryContext(t)
		sut := &ServiceAccountManager{hostConfiguration: nil}

		// when
		sut.SetHostConfiguration(hostConfigMock)

		// then
		assert.Equal(t, hostConfigMock, sut.hostConfiguration)
	})
}

func TestServiceAccountManager_GetHostConfiguration(t *testing.T) {
	t.Run("should get host config", func(t *testing.T) {
		// given
		hostConfigMock := newMockRegistryContext(t)
		sut := &ServiceAccountManager{hostConfiguration: hostConfigMock}

		// when
		actual := sut.GetHostConfiguration()

		// then
		assert.Equal(t, hostConfigMock, actual)
	})
}

func TestServiceAccountManager_GetServiceAccountData(t *testing.T) {
	t.Run("should fail to get username", func(t *testing.T) {
		// given
		hostConfigMock := newMockRegistryContext(t)
		hostConfigMock.EXPECT().Get("myService/username").Return("", assert.AnError)
		sut := &ServiceAccountManager{
			serviceName:       "myService",
			hostConfiguration: hostConfigMock,
		}

		// when
		actual, err := sut.GetServiceAccountData(testCtx)

		// then
		require.Error(t, err)
		assert.Empty(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get username for service account 'myService'")
	})
	t.Run("should fail to get password", func(t *testing.T) {
		// given
		hostConfigMock := newMockRegistryContext(t)
		hostConfigMock.EXPECT().Get("myService/username").Return("myUser", nil)
		hostConfigMock.EXPECT().Get("myService/password").Return("", assert.AnError)
		sut := &ServiceAccountManager{
			serviceName:       "myService",
			hostConfiguration: hostConfigMock,
		}

		// when
		actual, err := sut.GetServiceAccountData(testCtx)

		// then
		require.Error(t, err)
		assert.Empty(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get password for service account 'myService'")
	})
	t.Run("should fail to load key pair data", func(t *testing.T) {
		// given
		hostConfigMock := newMockRegistryContext(t)
		hostConfigMock.EXPECT().Get("myService/username").Return("myUser", nil)
		hostConfigMock.EXPECT().Get("myService/password").Return("myPassword", nil)
		globalConfigMock := newMockRegistryContext(t)
		globalConfigMock.EXPECT().Get(ssl.CertificateKeyRegistryKey).Return(testPrivateKeyString, nil)
		keyProviderMock := newMockKeyProvider(t)
		keyProviderMock.EXPECT().FromPrivateKey(testPrivateKeyBytes).Return(nil, assert.AnError)
		sut := &ServiceAccountManager{
			serviceName:         "myService",
			hostConfiguration:   hostConfigMock,
			globalConfiguration: globalConfigMock,
			keyProvider:         keyProviderMock,
		}

		// when
		actual, err := sut.GetServiceAccountData(testCtx)

		// then
		require.Error(t, err)
		assert.Empty(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to decrypt password for service account 'myService': failed to load key pair data")
	})
	t.Run("should fail to decrypt password", func(t *testing.T) {
		// given
		hostConfigMock := newMockRegistryContext(t)
		hostConfigMock.EXPECT().Get("myService/username").Return("myUser", nil)
		hostConfigMock.EXPECT().Get("myService/password").Return("myPassword", nil)
		globalConfigMock := newMockRegistryContext(t)
		globalConfigMock.EXPECT().Get(ssl.CertificateKeyRegistryKey).Return(testPrivateKeyString, nil)
		keyProviderMock := newMockKeyProvider(t)
		dummyKeyPair := &keys.KeyPair{}
		keyProviderMock.EXPECT().FromPrivateKey(testPrivateKeyBytes).Return(dummyKeyPair, nil)
		sut := &ServiceAccountManager{
			serviceName:         "myService",
			hostConfiguration:   hostConfigMock,
			globalConfiguration: globalConfigMock,
			keyProvider:         keyProviderMock,
		}

		// when
		actual, err := sut.GetServiceAccountData(testCtx)

		// then
		require.Error(t, err)
		assert.Empty(t, actual)
		assert.ErrorContains(t, err, "failed to decrypt password for service account 'myService': failed to decode base64 input")
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		hostConfigMock := newMockRegistryContext(t)
		hostConfigMock.EXPECT().Get("myService/username").Return("myUser", nil)
		hostConfigMock.EXPECT().Get("myService/password").Return("bXlQYXNzd29yZAo=", nil)
		globalConfigMock := newMockRegistryContext(t)
		globalConfigMock.EXPECT().Get(ssl.CertificateKeyRegistryKey).Return(testPrivateKeyString, nil)
		dummyKeyProvider := &keys.KeyProvider{Decrypter: func(_ io.Reader, _ *rsa.PrivateKey, cipherText []byte) ([]byte, error) {
			assert.Equal(t, []byte("myPassword\n"), cipherText)
			return []byte("myDecryptedPassword"), nil
		}}
		dummyKeyPair, err := dummyKeyProvider.Generate()
		require.NoError(t, err)

		keyProviderMock := newMockKeyProvider(t)
		keyProviderMock.EXPECT().FromPrivateKey(testPrivateKeyBytes).Return(dummyKeyPair, nil)
		sut := &ServiceAccountManager{
			serviceName:         "myService",
			hostConfiguration:   hostConfigMock,
			globalConfiguration: globalConfigMock,
			keyProvider:         keyProviderMock,
		}

		// when
		actual, err := sut.GetServiceAccountData(testCtx)

		// then
		require.NoError(t, err)
		assert.Equal(t, ServiceAccountData{Username: "myUser", Password: "myDecryptedPassword"}, actual)
	})
}

func TestServiceAccountData_String(t *testing.T) {
	t.Run("should return stringified data", func(t *testing.T) {
		// given
		sut := ServiceAccountData{
			Username: "myUsername",
			Password: "myPassword",
		}

		// when
		actual := sut.String()

		// then
		assert.Equal(t, "username:myUsername\npassword:myPassword", actual)
	})
}

func TestServiceAccountManager_Create(t *testing.T) {
	t.Run("should fail to load key pair data", func(t *testing.T) {
		// given
		hostConfigMock := newMockRegistryContext(t)
		keyProviderMock := newMockKeyProvider(t)
		keyProviderMock.EXPECT().FromPrivateKey(testPrivateKeyBytes).Return(nil, assert.AnError)
		globalConfigMock := newMockRegistryContext(t)
		globalConfigMock.EXPECT().Get(ssl.CertificateKeyRegistryKey).Return(testPrivateKeyString, nil)
		sut := &ServiceAccountManager{
			serviceName:         "myService",
			hostConfiguration:   hostConfigMock,
			globalConfiguration: globalConfigMock,
			keyProvider:         keyProviderMock,
		}

		// when
		actual, err := sut.Create(testCtx)

		// then
		require.Error(t, err)
		assert.Empty(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to encrypt the users password for service account 'myService': failed to load key pair data")
	})
	t.Run("should fail to encrypt password", func(t *testing.T) {
		// given
		hostConfigMock := newMockRegistryContext(t)

		dummyKeyProvider := &keys.KeyProvider{Encrypter: func(_ io.Reader, _ *rsa.PublicKey, _ []byte) ([]byte, error) {
			return nil, assert.AnError
		}}
		dummyKeyPair, err := dummyKeyProvider.Generate()
		require.NoError(t, err)

		keyProviderMock := newMockKeyProvider(t)
		keyProviderMock.EXPECT().FromPrivateKey(testPrivateKeyBytes).Return(dummyKeyPair, nil)
		globalConfigMock := newMockRegistryContext(t)
		globalConfigMock.EXPECT().Get(ssl.CertificateKeyRegistryKey).Return(testPrivateKeyString, nil)
		sut := &ServiceAccountManager{
			serviceName:         "myService",
			hostConfiguration:   hostConfigMock,
			globalConfiguration: globalConfigMock,

			keyProvider: keyProviderMock,
		}

		// when
		actual, err := sut.Create(testCtx)

		// then
		require.Error(t, err)
		assert.Empty(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to encrypt the users password for service account 'myService': failed to encrypt input text")
	})
	t.Run("should fail to set username", func(t *testing.T) {
		// given
		hostConfigMock := newMockRegistryContext(t)
		hostConfigMock.EXPECT().Set("myService/username", mock.AnythingOfType("string")).Return(assert.AnError)

		dummyKeyPair := generateTestKeyPair(t)

		keyProviderMock := newMockKeyProvider(t)
		keyProviderMock.EXPECT().FromPrivateKey(testPrivateKeyBytes).Return(dummyKeyPair, nil)
		globalConfigMock := newMockRegistryContext(t)
		globalConfigMock.EXPECT().Get(ssl.CertificateKeyRegistryKey).Return(testPrivateKeyString, nil)
		sut := &ServiceAccountManager{
			serviceName:         "myService",
			hostConfiguration:   hostConfigMock,
			globalConfiguration: globalConfigMock,

			keyProvider: keyProviderMock,
		}

		// when
		actual, err := sut.Create(testCtx)

		// then
		require.Error(t, err)
		assert.Empty(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to write username to registry for service account 'myService'")
	})
	t.Run("should fail to set password", func(t *testing.T) {
		// given
		hostConfigMock := newMockRegistryContext(t)
		hostConfigMock.EXPECT().Set("myService/username", mock.AnythingOfType("string")).Return(nil)
		hostConfigMock.EXPECT().Set("myService/password", mock.AnythingOfType("string")).Return(assert.AnError)

		dummyKeyPair := generateTestKeyPair(t)

		keyProviderMock := newMockKeyProvider(t)
		keyProviderMock.EXPECT().FromPrivateKey(testPrivateKeyBytes).Return(dummyKeyPair, nil)
		globalConfigMock := newMockRegistryContext(t)
		globalConfigMock.EXPECT().Get(ssl.CertificateKeyRegistryKey).Return(testPrivateKeyString, nil)
		sut := &ServiceAccountManager{
			serviceName:         "myService",
			hostConfiguration:   hostConfigMock,
			globalConfiguration: globalConfigMock,

			keyProvider: keyProviderMock,
		}

		// when
		actual, err := sut.Create(testCtx)

		// then
		require.Error(t, err)
		assert.Empty(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to write password to registry for service account 'myService'")
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		hostConfigMock := newMockRegistryContext(t)
		hostConfigMock.EXPECT().Set("myService/username", mock.AnythingOfType("string")).Return(nil)
		hostConfigMock.EXPECT().Set("myService/password", mock.AnythingOfType("string")).Return(nil)

		dummyKeyProvider := &keys.KeyProvider{Encrypter: func(_ io.Reader, _ *rsa.PublicKey, _ []byte) ([]byte, error) {
			return []byte("encryptedPassword"), nil
		}}
		dummyKeyPair, err := dummyKeyProvider.Generate()
		require.NoError(t, err)

		keyProviderMock := newMockKeyProvider(t)
		keyProviderMock.EXPECT().FromPrivateKey(testPrivateKeyBytes).Return(dummyKeyPair, nil)
		globalConfigMock := newMockRegistryContext(t)
		globalConfigMock.EXPECT().Get(ssl.CertificateKeyRegistryKey).Return(testPrivateKeyString, nil)
		sut := &ServiceAccountManager{
			serviceName:         "myService",
			hostConfiguration:   hostConfigMock,
			globalConfiguration: globalConfigMock,
			keyProvider:         keyProviderMock,
		}

		// when
		actual, err := sut.Create(testCtx)

		// then
		require.NoError(t, err)
		assert.NotEmpty(t, actual)
		assert.Contains(t, actual.Username, "myService")
		assert.NotEmpty(t, actual.Password)
	})
}

func generateTestKeyPair(t *testing.T) *keys.KeyPair {
	t.Helper()

	provider, _ := keys.NewKeyProvider("pkcs1v15")
	keyPair, _ := provider.Generate()
	return keyPair
}
