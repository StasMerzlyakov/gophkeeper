package domain

import "github.com/golang-jwt/jwt/v4"

// EMailData struct contains registration and login data
type EMailData struct {
	EMail    string
	Password string
}

// OPTPass struct is used for OTP password pass
type OTPPass struct {
	Pass string
}

// MasterKey struct is used in master key initialization process
type MasterKey struct {
	EncryptedMasterKey string
	MasterKeyPassHint  string
}

// SessionID is used during registration or authorization process
type SessionID string

// UserID is used to identify user in jwt token
type UserID string

// AuthorizationData struct is used for pass jwt token to client
type AuthorizationData struct {
}

type EMailStatus string

const (
	EMailBusy      EMailStatus = "BUSY"
	EMailAvailable EMailStatus = "FREE"
)

type RegistrationState string

const (
	// User created, OTP QR sent
	RegistrationStateInit RegistrationState = "INIT"
	// OTP passed, masterKey is required
	RegistrationStateAuth RegistrationState = "Auth"
)

// RegistrationData struct contains Registration process data
type RegistrationData struct {
	EMail           string
	PasswordHash    string
	PasswordSalt    string
	EncryptedOTPKey string
	State           RegistrationState
}

// FullRegistrationData stuct contains all registration data for user creation in StateFullStorage
type FullRegistrationData struct {
	UserID             UserID
	EMail              string
	PasswordHash       string
	PasswordSalt       string
	EncryptedOTPKey    string
	EncryptedMasterKey string
	MasterKeyHint      string
	HelloEncrypted     string
}

// MasterKeyData struct used in registration process
type MasterKeyData struct {
	EncryptedMasterKey string
	MasterKeyHint      string
	HelloEncrypted     string
}

type AuthorizationState string

const (
	// User created, OTP QR sent
	AuthorizationStateInit AuthorizationState = "INIT"
	// OTP passed
	AuthorizationStateCompleted AuthorizationState = "COMPLETED"
)

// LoginData stuct used in authorization process
type LoginData struct {
	UserID          UserID
	EMail           string
	PasswordHash    string
	PasswordSalt    string
	EncryptedOTPKey string
}

// HashData struct contains user password information for saving
type HashData struct {
	Hash string
	Salt string
}

type JWTToken string

type Claims struct {
	jwt.RegisteredClaims
	UserID UserID
}
