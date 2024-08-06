package domain

import (
	"fmt"
	"strconv"
)

// EncryptedBankCard used on server side
type EncryptedBankCard struct {
	Number  string
	Content string
}

// BankCardView used in view -> controller communication
type BankCardView struct {
	Number      string
	ExpiryMonth string
	ExpiryYear  string
	CVV         string
}

func (bcv BankCardView) ToBankCard() (*BankCard, error) {
	bankCard := &BankCard{
		Number: bcv.Number,
		CVV:    bcv.CVV,
	}

	expMonth, err := strconv.Atoi(bcv.ExpiryMonth)
	if err != nil {
		return nil, fmt.Errorf("bankCard data err - wrong expiryMoth value")
	}
	bankCard.ExpiryMonth = expMonth

	expEear, err := strconv.Atoi(bcv.ExpiryYear)
	if err != nil {
		return nil, fmt.Errorf("bankCard data err - wrong expiryYear value")
	}

	if expEear < 100 {
		expEear += 2000
	}
	bankCard.ExpiryYear = expEear

	return bankCard, nil
}

// BankCard bank card data
type BankCard struct {
	// Type is an optional string with one of the supported card types
	Type string `json:"type,omitempty"`
	// Number is the credit card number
	Number string `json:"number,omitempty"`
	// ExpiryMonth is the credit card expiration month
	ExpiryMonth int `json:"exporityMonth,omitempty"`
	// ExpiryYear is the credit card expiration year
	ExpiryYear int `json:"exporityYear,omitempty"`
	// CVV is the credit card CVV code
	CVV string `json:"cvv,omitempty"`
}

// EncryptedUserPasswordData used on server side
type EncryptedUserPasswordData struct {
	Hint    string
	Content string
}

// UserPasswordData user login/password data
type UserPasswordData struct {
	// SiteURL or other hint
	Hint string `json:"hint,omitempty"`
	// Login is user login
	Login string `json:"login,omitempty"`
	// Password is user password
	Passwrod string `json:"password,omitempty"`
}
