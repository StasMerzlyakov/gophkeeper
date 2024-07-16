package domain

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

// AuthorizationData struct is used for pass jwt token to client
type AuthorizationData struct {
	AccessToken string
}

type EMailStatus string

const (
	EMailBusy      EMailStatus = "BUSY"
	EMailAvailable EMailStatus = "FREE"
)

type AuthorizationState string

const (
	// User created, OTP QR sent
	AuthorizationStateInit AuthorizationState = "INIT"
	// OTP passed
	AuthorizationStateCompleted AuthorizationState = "COMPLETED"
)
