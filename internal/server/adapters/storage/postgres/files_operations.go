package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (st *storage) GetUserFilesBucket(ctx context.Context) (string, error) {
	userID, err := domain.GetUserID(ctx)
	action := domain.GetAction(1)
	log := domain.GetCtxLogger(ctx)
	if err != nil {
		return "", fmt.Errorf("%w userID is not set", domain.ErrServerInternal)
	}

	var bucket string
	err = st.pPool.QueryRow(ctx, "select bucket from user_info where user_id = $1", userID).Scan(&bucket)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// wrong userID
			return "", fmt.Errorf("%w - user %d is not exists", domain.ErrClientDataIncorrect, userID)
		}
		log.Infow(action, "err", err.Error())
		return "", fmt.Errorf("%w - %s", domain.ErrServerInternal, err.Error())
	}
	return bucket, nil
}
