package interceptor

import (
	"context"
	"fmt"
	"strings"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func EncrichWithRequestIDUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		logger := domain.GetApplicationLogger()
		requestUUID := uuid.New()
		enrichedCtx := domain.EnrichWithRequestIDLogger(ctx, requestUUID, logger)
		return handler(enrichedCtx, req)
	}
}

func EncrichWithRequestIDStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Pass stream context complexity
		// https://stackoverflow.com/questions/63518470/grpc-stream-interceptor-not-passing-context-to-request-method
		// see https://github.com/fru-io/api-common/blob/main/interceptors/state.go#L207

		logger := domain.GetApplicationLogger()
		requestUUID := uuid.New()
		w := newStreamContextWrapper(ss)
		eCtx := domain.EnrichWithRequestIDLogger(w.Context(), requestUUID, logger)
		w.SetContext(eCtx)
		return handler(srv, w)
	}
}

func ErrorCodeUnaryInteceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			log := domain.GetCtxLogger(ctx)
			log.Infow("error", "err", err.Error())
			resCode := MapDomainErrorToGRPCCodeErr(err)
			return resp, status.Error(resCode, "")
		}
		return resp, nil
	}
}

func ErrorCodeStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, ss)
		if err != nil {
			log := domain.GetCtxLogger(ss.Context())
			log.Infow("error", "err", err.Error())
			resCode := MapDomainErrorToGRPCCodeErr(err)
			return status.Error(resCode, "")
		}
		return nil
	}
}

func JWTUnaryInterceptor(tokenSecret []byte, needToBeAuthentificated []string) grpc.UnaryServerInterceptor {
	needAuthentificatedFn := func(method string) bool {
		for _, v := range needToBeAuthentificated {
			if strings.Contains(method, v) {
				return true
			}
		}
		return false
	}
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		method := info.FullMethod
		if needAuthentificatedFn(method) {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, fmt.Errorf("%w - %v", domain.ErrNotAuthorized, "metadata is not provided")
			}
			values := md[domain.AuthorizationMetadataTokenName]
			if len(values) == 0 {
				return nil, fmt.Errorf("%w - %v", domain.ErrNotAuthorized, "authorization token is not provided")
			}
			accessToken := values[0]
			userID, err := domain.ParseJWTToken(tokenSecret, domain.JWTToken(accessToken))
			if err != nil {
				return nil, fmt.Errorf("%w - %v", domain.ErrNotAuthorized, err.Error())
			}
			eCtx := domain.EnrichWithUserID(ctx, userID)
			return handler(eCtx, req)
		} else {
			return handler(ctx, req)
		}
	}
}

func JWTStreamInterceptor(tokenSecret []byte) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		// pass stream context complexity
		// https://stackoverflow.com/questions/63518470/grpc-stream-interceptor-not-passing-context-to-request-method
		// see https://github.com/fru-io/api-common/blob/main/interceptors/state.go#L207
		w := newStreamContextWrapper(ss)

		md, ok := metadata.FromIncomingContext(w.Context())
		if !ok {
			return fmt.Errorf("%w - %v", domain.ErrNotAuthorized, "metadata is not provided")
		}
		values := md[domain.AuthorizationMetadataTokenName]
		if len(values) == 0 {
			return fmt.Errorf("%w - %v", domain.ErrNotAuthorized, "authorization token is not provided")
		}
		accessToken := values[0]
		userID, err := domain.ParseJWTToken(tokenSecret, domain.JWTToken(accessToken))
		if err != nil {
			return fmt.Errorf("%w - %v", domain.ErrNotAuthorized, err.Error())
		}
		eCtx := domain.EnrichWithUserID(w.Context(), userID)
		w.SetContext(eCtx)

		return handler(srv, w)
	}
}

type StreamContextWrapper interface {
	grpc.ServerStream
	SetContext(context.Context)
}

type wrapper struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrapper) Context() context.Context {
	return w.ctx
}

func (w *wrapper) SetContext(ctx context.Context) {
	w.ctx = ctx
}

func newStreamContextWrapper(inner grpc.ServerStream) StreamContextWrapper {
	ctx := inner.Context()
	return &wrapper{
		inner,
		ctx,
	}
}
