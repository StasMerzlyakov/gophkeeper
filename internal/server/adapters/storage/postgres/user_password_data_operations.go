package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (st *storage) GetUserPasswordDataList(ctx context.Context) ([]domain.EncryptedUserPasswordData, error) {
	userID, err := domain.GetUserID(ctx)
	action := domain.GetAction(0)
	log := domain.GetCtxLogger(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w userID is not set", domain.ErrServerInternal)
	}

	rows, err := st.pPool.Query(ctx, "select hint, content from user_password_data where user_id = $1", userID)

	if err != nil {
		log.Infow(action, "err", err.Error())
		return nil, fmt.Errorf("%w - %s", domain.ErrServerInternal, err.Error())
	}

	defer rows.Close()

	var result []domain.EncryptedUserPasswordData

	for rows.Next() {
		var data domain.EncryptedUserPasswordData
		err = rows.Scan(&data.Hint, &data.Content)
		if err != nil {
			log.Infow(action, "err", err.Error())
			return nil, fmt.Errorf("%w - %s", domain.ErrServerInternal, err.Error())
		}
		result = append(result, data)
	}

	err = rows.Err()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return result, nil
		}
		log.Infow(action, "err", err.Error())
		return nil, fmt.Errorf("%w - %s", domain.ErrServerInternal, err.Error())
	}

	return result, nil
}

func (st *storage) CreateUserPasswordData(ctx context.Context, data *domain.EncryptedUserPasswordData) error {
	userID, err := domain.GetUserID(ctx)
	action := domain.GetAction(0)
	log := domain.GetCtxLogger(ctx)
	if err != nil {
		log.Infow(action, "err", err.Error())
		return fmt.Errorf("%w userID is not set", domain.ErrServerInternal)
	}
	var id int64

	if err := st.pPool.QueryRow(ctx,
		`insert into user_password_data(hint, content, user_id) values ($1, $2, $3) on conflict("hint","user_id") do nothing returning id;
	  	`, data.Hint, data.Content, userID).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// email already registered
			log.Infow(action, "err", fmt.Sprintf("used data with hint %s already registered", data.Hint))
			return fmt.Errorf("%w - used data with hint %s already registereed", domain.ErrClientDataIncorrect, data.Hint)
		}
		log.Infow(action, "err", err.Error())
		return fmt.Errorf("%w - %s", domain.ErrServerInternal, err.Error())
	} else {
		log.Debugw(action, "msg", fmt.Sprintf("used data %v for userID %v registered, id = %v", data.Hint, userID, id))
		return nil
	}
}

func (st *storage) UpdateUserPasswordData(ctx context.Context, data *domain.EncryptedUserPasswordData) error {
	userID, err := domain.GetUserID(ctx)
	action := domain.GetAction(0)
	log := domain.GetCtxLogger(ctx)
	if err != nil {
		log.Infow(action, "err", err.Error())
		return fmt.Errorf("%w userID is not set", domain.ErrServerInternal)
	}
	var id int64

	if err := st.pPool.QueryRow(ctx,
		`update user_password_data set content = $1 where hint = $2 and user_id = $3 returning id;
	  	`, data.Content, data.Hint, userID).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// email already registered
			log.Infow(action, "err", fmt.Sprintf("user data with hint %v for user %v is not exists", data.Hint, userID))
			return fmt.Errorf("%w - user data with hint %v for user %v is not exists", domain.ErrClientDataIncorrect, data.Hint, userID)
		}
		log.Infow(action, "err", err.Error())
		return fmt.Errorf("%w - %s", domain.ErrServerInternal, err.Error())
	} else {
		log.Debugw(action, "msg", fmt.Sprintf("user data %v for userID %v updated", data.Hint, userID))
		return nil
	}
}

func (st *storage) DeleteUserPasswordData(ctx context.Context, hint string) error {
	userID, err := domain.GetUserID(ctx)
	action := domain.GetAction(0)
	log := domain.GetCtxLogger(ctx)
	if err != nil {
		log.Infow(action, "err", err.Error())
		return fmt.Errorf("%w userID is not set", domain.ErrServerInternal)
	}
	var id int64

	if err := st.pPool.QueryRow(ctx,
		`delete from user_password_data where hint = $1 and user_id = $2 returning id;
	  	`, hint, userID).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// email already registered
			log.Infow(action, "err", fmt.Sprintf("user data with hint %v for user %v nit exists", hint, userID))
			return fmt.Errorf("%w - user data with hint %v for user %v is not exists", domain.ErrClientDataIncorrect, hint, userID)
		}
		log.Infow(action, "err", err.Error())
		return fmt.Errorf("%w - %s", domain.ErrServerInternal, err.Error())
	} else {
		log.Debugw(action, "msg", fmt.Sprintf("user data %v for userID %v deleted", hint, userID))
		return nil
	}
}
