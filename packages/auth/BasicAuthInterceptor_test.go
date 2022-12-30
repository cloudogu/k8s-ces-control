package auth

//
// import (
// 	"context"
// 	"errors"
// 	"github.com/cloudogu/cesappd/account"
// 	"github.com/cloudogu/cesappd/generated/dummy"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/metadata"
// 	"google.golang.org/grpc/test/bufconn"
// 	"log"
// 	"net"
// 	"testing"
// )
//
// var (
// 	listener *bufconn.Listener
// )
//
// const bufferSize1MB = 1024 * 1024
//
// func Test_authorize(t *testing.T) {
// 	emptyGetCredentialsFunc := func(serviceName string) (account.ServiceAccountData, error) {
// 		return account.ServiceAccountData{}, nil
// 	}
//
// 	correctGetCredentialsFunc := func(serviceName string) (account.ServiceAccountData, error) {
// 		return account.ServiceAccountData{Username: "admin", Password: "password"}, nil
// 	}
//
// 	errorGetCredentialsFunc := func(serviceName string) (account.ServiceAccountData, error) {
// 		return account.ServiceAccountData{}, errors.New("Oops sth went wrong")
// 	}
//
// 	wrongGetCredentialsFunc := func(serviceName string) (account.ServiceAccountData, error) {
// 		return account.ServiceAccountData{Username: "nouser", Password: "nopasswrd"}, nil
// 	}
//
// 	t.Run("should return error as no servicename header is present", func(t *testing.T) {
// 		md := metadata.MD{}
// 		err := authorize(&md, "/doguAdministration.DoguAdministration/GetDoguList", emptyGetCredentialsFunc)
// 		require.Error(t, err)
// 		assert.Contains(t, err.Error(), "cannot find \"servicename\" header")
// 	})
//
// 	t.Run("should return error as no authorization header is present", func(t *testing.T) {
// 		md := metadata.MD{}
// 		md.Set("servicename", "test")
// 		err := authorize(&md, "/doguAdministration.DoguAdministration/GetDoguList", emptyGetCredentialsFunc)
// 		require.Error(t, err)
// 		assert.Contains(t, err.Error(), "cannot find \"authorization\" header")
// 	})
//
// 	t.Run("should return error as authorization header is faulty", func(t *testing.T) {
// 		md := metadata.MD{}
// 		md.Set("servicename", "test")
// 		md.Set("authorization", "blub-blub")
// 		err := authorize(&md, "/doguAdministration.DoguAdministration/GetDoguList", emptyGetCredentialsFunc)
// 		require.Error(t, err)
// 		assert.Contains(t, err.Error(), "header wrongly formatted")
// 	})
//
// 	t.Run("should return error as serviceAccountData could not be accessed", func(t *testing.T) {
// 		md := metadata.MD{}
// 		md.Set("servicename", "test")
// 		md.Set("authorization", "user:password")
// 		err := authorize(&md, "/doguAdministration.DoguAdministration/GetDoguList", errorGetCredentialsFunc)
// 		require.Error(t, err)
// 		assert.Contains(t, err.Error(), "could not get serviceAccountData")
// 	})
//
// 	t.Run("should return error as username does not match", func(t *testing.T) {
// 		md := metadata.MD{}
// 		md.Set("servicename", "test")
// 		md.Set("authorization", "user:password")
// 		err := authorize(&md, "/doguAdministration.DoguAdministration/GetDoguList", wrongGetCredentialsFunc)
// 		require.Error(t, err)
// 		assert.Contains(t, err.Error(), "User/Password invalid")
// 	})
//
// 	t.Run("should return error as password does not match", func(t *testing.T) {
// 		md := metadata.MD{}
// 		md.Set("servicename", "test")
// 		md.Set("authorization", "nouser:password")
// 		err := authorize(&md, "/doguAdministration.DoguAdministration/GetDoguList", wrongGetCredentialsFunc)
// 		require.Error(t, err)
// 		assert.Contains(t, err.Error(), "User/Password invalid")
// 	})
//
// 	t.Run("should return no error as auth header and credentials match", func(t *testing.T) {
// 		md := metadata.MD{}
// 		md.Set("servicename", "test")
// 		md.Set("authorization", "admin:password")
// 		err := authorize(&md, "/doguAdministration.DoguAdministration/GetDoguList", correctGetCredentialsFunc)
// 		require.NoError(t, err)
// 	})
// }
//
// func Test_UnaryInterceptor(t *testing.T) {
// 	t.Run("should fail as no servicename is present in the request header", func(t *testing.T) {
// 		setupTestServer()
//
// 		client, ctx, conn := getGrpcClient(t)
// 		defer func() { _ = conn.Close() }()
//
// 		response, err := client.HelloWorld(*ctx, &dummy.HelloWorldRequest{})
// 		require.Error(t, err)
// 		assert.Contains(t, err.Error(), "cannot find \"servicename\" header in metadata")
// 		require.Nil(t, response)
// 	})
// 	t.Run("should fail as no basic-auth is present in the request header", func(t *testing.T) {
// 		setupTestServer()
//
// 		client, ctx, conn := getGrpcClient(t)
// 		defer func() { _ = conn.Close() }()
//
// 		*ctx = metadata.AppendToOutgoingContext(*ctx, "servicename", "testservice")
//
// 		response, err := client.HelloWorld(*ctx, &dummy.HelloWorldRequest{})
// 		require.Error(t, err)
// 		assert.Contains(t, err.Error(), "cannot find \"authorization\" header in metadata")
// 		require.Nil(t, response)
// 	})
// 	t.Run("should fail as basic-auth is wrongly formatted", func(t *testing.T) {
// 		setupTestServer()
//
// 		client, ctx, conn := getGrpcClient(t)
// 		defer func() { _ = conn.Close() }()
//
// 		*ctx = metadata.AppendToOutgoingContext(*ctx, "authorization", "user_password")
// 		*ctx = metadata.AppendToOutgoingContext(*ctx, "servicename", "testservice")
//
// 		response, err := client.HelloWorld(*ctx, &dummy.HelloWorldRequest{})
// 		require.Error(t, err)
// 		assert.Contains(t, err.Error(), "\"authorization\" header wrongly formatted")
// 		require.Nil(t, response)
// 	})
// 	t.Run("should fail as credentials can't be retrieved", func(t *testing.T) {
// 		// If this test loads forever, try to remove your local file /etc/cesapp/configuration.json
// 		setupTestServer()
//
// 		client, ctx, conn := getGrpcClient(t)
// 		defer func() { _ = conn.Close() }()
//
// 		*ctx = metadata.AppendToOutgoingContext(*ctx, "authorization", "user:password")
// 		*ctx = metadata.AppendToOutgoingContext(*ctx, "servicename", "testservice")
//
// 		response, err := client.HelloWorld(*ctx, &dummy.HelloWorldRequest{})
// 		require.Error(t, err)
// 		assert.Contains(t, err.Error(), "could not get serviceAccountData")
// 		require.Nil(t, response)
// 	})
// }
//
// func getGrpcClient(t *testing.T) (dummy.HelloWorldServiceClient, *context.Context, *grpc.ClientConn) {
// 	ctx := context.Background()
// 	conn, clientErr := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
// 	if clientErr != nil {
// 		t.Fatalf("Failed to dial bufnet: %v", clientErr)
// 	}
// 	return dummy.NewHelloWorldServiceClient(conn), &ctx, conn
// }
//
// func bufDialer(context.Context, string) (net.Conn, error) {
// 	return listener.Dial()
// }
//
// type server struct {
// 	dummy.UnimplementedHelloWorldServiceServer
// }
//
// func (s *server) HelloWorld(_ context.Context, _ *dummy.HelloWorldRequest) (*dummy.HelloWorldResponse, error) {
// 	return &dummy.HelloWorldResponse{}, nil
// }
//
// func setupTestServer() {
//
// 	s := grpc.NewServer(grpc.UnaryInterceptor(BasicAuthUnaryInterceptor))
// 	listener = bufconn.Listen(bufferSize1MB)
//
// 	srv := &server{}
//
// 	dummy.RegisterHelloWorldServiceServer(s, srv)
// 	go func() {
// 		if err := s.Serve(listener); err != nil {
// 			log.Fatalf("Server exited with error: %v", err)
// 		}
// 	}()
// }
