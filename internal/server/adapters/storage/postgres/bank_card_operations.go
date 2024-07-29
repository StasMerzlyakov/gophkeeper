package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (st *storage) GetBankCardList(ctx context.Context) ([]domain.EncryptedBankCard, error) {
	userID, err := domain.GetUserID(ctx)
	action := domain.GetAction(0)
	log := domain.GetCtxLogger(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w userID is not set", domain.ErrServerInternal)
	}

	rows, err := st.pPool.Query(ctx, "select number, content from bank_card where user_id = $1", userID)

	if err != nil {
		log.Infow(action, "err", err.Error())
		return nil, fmt.Errorf("%w - %s", domain.ErrServerInternal, err.Error())
	}

	defer rows.Close()

	var result []domain.EncryptedBankCard

	for rows.Next() {
		var card domain.EncryptedBankCard
		err = rows.Scan(&card.Number, &card.Content)
		if err != nil {
			log.Infow(action, "err", err.Error())
			return nil, fmt.Errorf("%w - %s", domain.ErrServerInternal, err.Error())
		}
		result = append(result, card)
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

func (st *storage) CreateBankCard(ctx context.Context, bnkCard *domain.EncryptedBankCard) error {
	userID, err := domain.GetUserID(ctx)
	action := domain.GetAction(0)
	log := domain.GetCtxLogger(ctx)
	if err != nil {
		log.Infow(action, "err", err.Error())
		return fmt.Errorf("%w userID is not set", domain.ErrServerInternal)
	}
	var id int64

	if err := st.pPool.QueryRow(ctx,
		`insert into bank_card(number, content, user_id) values ($1, $2, $3) on conflict("number","user_id") do nothing returning id;
	  	`, bnkCard.Number, bnkCard.Content, userID).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// email already registered
			log.Infow(action, "err", fmt.Sprintf("bank card with number %s already registered", bnkCard.Number))
			return fmt.Errorf("%w - bank card with number %s already registered", domain.ErrClientDataIncorrect, bnkCard.Number)
		}
		log.Infow(action, "err", err.Error())
		return fmt.Errorf("%w - %s", domain.ErrServerInternal, err.Error())
	} else {
		log.Debugw(action, "msg", fmt.Sprintf("bank card %v for userID %v registered, id = %v", bnkCard.Number, userID, id))
		return nil
	}
}

func (st *storage) UpdateBankCard(ctx context.Context, bnkCard *domain.EncryptedBankCard) error {
	userID, err := domain.GetUserID(ctx)
	action := domain.GetAction(0)
	log := domain.GetCtxLogger(ctx)
	if err != nil {
		log.Infow(action, "err", err.Error())
		return fmt.Errorf("%w userID is not set", domain.ErrServerInternal)
	}
	var id int64

	if err := st.pPool.QueryRow(ctx,
		`update bank_card set content = $1 where number = $2 and user_id = $3 returning id;
	  	`, bnkCard.Content, bnkCard.Number, userID).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// email already registered
			log.Infow(action, "err", fmt.Sprintf("bank card with number %v for user %v nit exists", bnkCard.Number, userID))
			return fmt.Errorf("%w - bank card with number %v for user %v nit exists", domain.ErrClientDataIncorrect, bnkCard.Number, userID)
		}
		log.Infow(action, "err", err.Error())
		return fmt.Errorf("%w - %s", domain.ErrServerInternal, err.Error())
	} else {
		log.Debugw(action, "msg", fmt.Sprintf("bank card %v for userID %v updated", bnkCard.Number, userID))
		return nil
	}
}

func (st *storage) DeleteBankCard(ctx context.Context, number string) error {
	userID, err := domain.GetUserID(ctx)
	action := domain.GetAction(0)
	log := domain.GetCtxLogger(ctx)
	if err != nil {
		log.Infow(action, "err", err.Error())
		return fmt.Errorf("%w userID is not set", domain.ErrServerInternal)
	}
	var id int64

	if err := st.pPool.QueryRow(ctx,
		`delete from bank_card where number = $1 and user_id = $2 returning id;
	  	`, number, userID).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// email already registered
			log.Infow(action, "err", fmt.Sprintf("bank card with number %v for user %v nit exists", number, userID))
			return fmt.Errorf("%w - bank card with number %v for user %v nit exists", domain.ErrClientDataIncorrect, number, userID)
		}
		log.Infow(action, "err", err.Error())
		return fmt.Errorf("%w - %s", domain.ErrServerInternal, err.Error())
	} else {
		log.Debugw(action, "msg", fmt.Sprintf("bank card %v for userID %v deleted", number, userID))
		return nil
	}
}
