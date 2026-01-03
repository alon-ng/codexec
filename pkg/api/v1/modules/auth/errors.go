package auth

var (
	ErrEmailAlreadyExists    = "User with this email already exists"
	ErrSignupFailed          = "Failed to sign up user"
	ErrLoginFailed           = "Failed to login"
	ErrInvalidCredentials    = "Invalid email or password"
	ErrTokenGenerationFailed = "Failed to generate token"
)
