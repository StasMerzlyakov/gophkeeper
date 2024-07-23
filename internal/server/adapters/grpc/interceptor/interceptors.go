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

// EncrichWithRequestIDInterceptor Добавляет к запросу RequestID и устанавливает в контекст логгер
func EncrichWithRequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		logger := domain.GetApplicationLogger()
		requestUUID := uuid.New()
		enrichedCtx := domain.EnrichWithRequestIDLogger(ctx, requestUUID, logger)
		return handler(enrichedCtx, req)
	}
}

// ErrorCodeInteceptor отвечает за преобразование ошибки выполнения в grpc code
func ErrorCodeInteceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			resCode := MapDomainErrorToGRPCCodeErr(err)
			return resp, status.Error(resCode, "")
		}
		return resp, nil
	}
}

// JWTInterceptor отвечает за разбор JWT токена
func JWTInterceptor(tokenSecret []byte, needToBeAuthentificated []string) grpc.UnaryServerInterceptor {

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
