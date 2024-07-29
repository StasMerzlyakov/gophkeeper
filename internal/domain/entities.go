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
	MasterPasswordHint string
}

// SessionID is used during registration or authorization process
type SessionID string

// UserID is used to identify user in jwt token
type UserID int64

// AuthorizationData struct is used for pass jwt token to client
type AuthorizationData struct {
}

type ClientStatus string

const (
	ClientStatusOnline  ClientStatus = "ONLINE"
	ClientStatusOffline ClientStatus = "OFFLINE"
)

type EMailStatus string

const (
	EMailBusy      EMailStatus = "BUSY"
	EMailAvailable EMailStatus = "AVAILABLE"
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
	EMail              string
	PasswordHash       string
	PasswordSalt       string
	EncryptedOTPKey    string
	MasterPasswordHint string
	HelloEncrypted     string
}

// UnencryptedMasterKeyData struct is used in registration process on client side
type UnencryptedMasterKeyData struct {
	MasterPassword     string
	MasterPasswordHint string
}

// MasterKeyData struct is used in registration process
type MasterKeyData struct {
	MasterPasswordHint string
	HelloEncrypted     string
}

// HelloData struct is used in authorization process
type HelloData struct {
	HelloEncrypted     string
	MasterPasswordHint string
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
