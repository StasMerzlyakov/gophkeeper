package handler

import (
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func WrapErr(err error) error {
	action := domain.GetAction(1)
	return fmt.Errorf("%v err - %w", action, err)
}
