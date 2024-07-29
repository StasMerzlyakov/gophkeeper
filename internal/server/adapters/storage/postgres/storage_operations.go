package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (st *storage) IsEMailAvailable(ctx context.Context, email string) (bool, error) {
	var count int
	err := st.pPool.QueryRow(ctx, "select count(user_id) from user_info where email = $1", email).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("%w - %s", domain.ErrServerInternal, err.Error())
	} else {
		return count == 0, nil
	}
}

func (st *storage) Registrate(ctx context.Context, data *domain.FullRegistrationData) error {

	var userID int64

	if err := st.pPool.QueryRow(ctx,
		`insert into user_info(email, pass_hash, pass_salt, otp_key, master_key, master_hint, hello_encrypted) 
		 values ($1, $2, $3, $4, $5, $6, $7) on conflict("email") do nothing returning user_id;
	  	`, data.EMail, data.PasswordHash, data.PasswordSalt, data.EncryptedOTPKey,
		data.EncryptedMasterKey, data.MasterKeyHint, data.HelloEncrypted).Scan(&userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// email already registered
			return fmt.Errorf("%w - email %s already registered", domain.ErrClientDataIncorrect, data.EMail)
		}
		return fmt.Errorf("%w - %s", domain.ErrServerInternal, err.Error())
	} else {
		return nil
	}
}

func (st *storage) GetLoginData(ctx context.Context, email string) (*domain.LoginData, error) {

	var loginData domain.LoginData
	err := st.pPool.QueryRow(ctx, "select user_id, email, pass_hash, pass_salt, otp_key from user_info where email = $1", email).
		Scan(&loginData.UserID, &loginData.EMail, &loginData.PasswordHash, &loginData.PasswordSalt, &loginData.EncryptedOTPKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// email already registered
			return nil, fmt.Errorf("%w - email %s not registered", domain.ErrClientDataIncorrect, email)
		}
		return nil, fmt.Errorf("%w - %s", domain.ErrServerInternal, err.Error())
	} else {
		return &loginData, nil
	}
}

func (st *storage) GetHelloData(ctx context.Context) (*domain.HelloData, error) {
	userID, err := domain.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w userID is not set", domain.ErrServerInternal)
	}

	var helloData domain.HelloData
	err = st.pPool.QueryRow(ctx, "select hello_encrypted, master_key, master_hint from user_info where user_id = $1", userID).
		Scan(&helloData.HelloEncrypted, &helloData.EncryptedMasterKey, &helloData.MasterKeyPassHint)
	if err != nil {
		return nil, fmt.Errorf("%w - %s", domain.ErrServerInternal, err.Error())
	} else {
		return &helloData, nil
	}
}
