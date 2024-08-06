package interceptor_test

import (
	"fmt"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/adapters/grpc/interceptor"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestMapDomainErrorToGRPCCodeErr(t *testing.T) {
	assert.Equal(t, codes.OK, interceptor.MapDomainErrorToGRPCCodeErr(nil))
	assert.Equal(t, codes.Internal, interceptor.MapDomainErrorToGRPCCodeErr(fmt.Errorf("%w err", domain.ErrServerInternal)))
	assert.Equal(t, codes.PermissionDenied, interceptor.MapDomainErrorToGRPCCodeErr(fmt.Errorf("%w err", domain.ErrNotAuthorized)))
	assert.Equal(t, codes.Internal, interceptor.MapDomainErrorToGRPCCodeErr(fmt.Errorf("%w err", domain.ErrDublicateKeyViolation)))
	assert.Equal(t, codes.InvalidArgument, interceptor.MapDomainErrorToGRPCCodeErr(fmt.Errorf("%w err", domain.ErrDataNotExists)))
	assert.Equal(t, codes.InvalidArgument, interceptor.MapDomainErrorToGRPCCodeErr(fmt.Errorf("%w err", domain.ErrAuthDataIncorrect)))
	assert.Equal(t, codes.InvalidArgument, interceptor.MapDomainErrorToGRPCCodeErr(fmt.Errorf("%w err", domain.ErrClientDataIncorrect)))
	assert.Equal(t, codes.Internal, interceptor.MapDomainErrorToGRPCCodeErr(fmt.Errorf("new")))
}
