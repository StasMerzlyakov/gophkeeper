package domain

type RegistrationState string

const (
	// User created, OTP QR sent
	RegistrationStateInit RegistrationState = "INIT"
	// OTP passed, masterKey is required
	RegistrationStateAuthPassed RegistrationState = "AUTH"
)

// RegistrationData struct contains Registration process
type RegistrationData struct {
	EMail            string
	PasswordHash     string
	PasswordSalt     string
	EncryptedOTPPass string
	State            RegistrationState
}
