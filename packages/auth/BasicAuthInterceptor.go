package auth

import (
	"context"
	"github.com/cloudogu/k8s-ces-control/packages/authHelper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
)

func BasicAuthUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logger := log.FromContext(ctx)
	logger.Info("Interceptor called; FullMethod: ", info.FullMethod)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	getCredentials := authHelper.GetServiceAccountCredentials
	err := authorize(ctx, &md, info.FullMethod, getCredentials)

	if err != nil {
		logger.Error(err, "cannot authorize request")
		return nil, err
	}

	return handler(ctx, req)
}

func authorize(ctx context.Context, metadata *metadata.MD, _ string, getCredentials authHelper.AuthenticationFunc) error {
	log.FromContext(ctx).Info("cesappd: metadata %v", metadata)

	if metadata.Get("servicename") == nil {
		return status.Errorf(codes.InvalidArgument, `cannot find "servicename" header in metadata`)
	}
	serviceName := metadata.Get("servicename")[0]

	// handle basic authHelper without base64
	if metadata.Get("authorization") == nil {
		return status.Errorf(codes.InvalidArgument, `cannot find "authorization" header in metadata`)
	}
	auth := metadata.Get("authorization")[0]
	if len(auth) == 0 || !strings.ContainsRune(auth, ':') {
		return status.Errorf(codes.InvalidArgument, `"authorization" header wrongly formatted`)
	}
	split := strings.Split(auth, ":")

	username := split[0]
	password := split[1]

	// check is user exists
	serviceAccountCredentials, err := getCredentials(ctx, serviceName)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "could not get serviceAccountData: %v", err)
	}

	// check if his password is valid
	if serviceAccountCredentials.Username != username {
		return status.Errorf(codes.Unauthenticated, "User/Password invalid")
	}

	// check if his password is valid
	if serviceAccountCredentials.Password != password {
		return status.Errorf(codes.Unauthenticated, "User/Password invalid")
	}

	return nil
}
