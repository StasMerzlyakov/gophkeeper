package interceptor

import (
	"errors"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	gp "google.golang.org/grpc/codes"
)

func MapDomainErrorToGRPCCodeErr(err error) gp.Code {

	if err == nil {
		return gp.OK
	}

	if errors.Is(err, domain.ErrServerInternal) {
		return gp.Internal
	}

	if errors.Is(err, domain.ErrNotAuthorized) {
		return gp.PermissionDenied
	}

	if errors.Is(err, domain.ErrDublicateKeyViolation) {
		return gp.Internal
	}

	if errors.Is(err, domain.ErrDataNotExists) {
		return gp.InvalidArgument
	}

	if errors.Is(err, domain.ErrAuthDataIncorrect) {
		return gp.InvalidArgument
	}

	if errors.Is(err, domain.ErrClientDataIncorrect) {
		return gp.InvalidArgument
	}

	return gp.Internal
}
