package postgres

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func (st *storage) IsEMailAvailable(ctx context.Context, email string) (bool, error) {
	var count int
	err := st.pPool.QueryRow(ctx, "select count(userID) from user_info where email = $1", email).Scan(&count)

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
		 values ($1, $2, $3, $4, $5, $6, $7) returning userId;
	  	`, data.EMail, data.PasswordHash, data.PasswordSalt, data.EncryptedOTPKey,
		data.EncryptedMasterKey, data.MasterKeyHint, data.HelloEncrypted).Scan(&userID); err != nil {
		return fmt.Errorf("%w - %s", domain.ErrServerInternal, err.Error())
	} else {
		return nil
	}
}

func (st *storage) GetLoginData(ctx context.Context, email string) (*domain.LoginData, error) {

	var loginData domain.LoginData
	err := st.pPool.QueryRow(ctx, "select userID, email, pass_hash, pass_salt, otp_key from user_info where email = $1", email).
		Scan(&loginData.UserID, &loginData.EMail, &loginData.PasswordHash, &loginData.PasswordSalt, &loginData.EncryptedOTPKey)

	if err != nil {
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
	err = st.pPool.QueryRow(ctx, "select hello_encrypted, master_key, master_hint from user_info where userId = $1", userID).
		Scan(&helloData.HelloEncrypted, &helloData.EncryptedMasterKey, &helloData.MasterKeyHint)
	if err != nil {
		return nil, fmt.Errorf("%w - %s", domain.ErrServerInternal, err.Error())
	} else {
		return &helloData, nil
	}
}
