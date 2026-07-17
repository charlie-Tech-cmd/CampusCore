package config

import "time"

// AuthConfig contains authentication and authorization settings.
type AuthConfig struct {
	// JWTSecret is the secret key used to sign JWT tokens.
	JWTSecret string

	// AccessTokenExpiry defines how long an access token remains valid.
	AccessTokenExpiry time.Duration

	// RefreshTokenExpiry defines how long a refresh token remains valid.
	RefreshTokenExpiry time.Duration

	// Issuer identifies the application that issued the token.
	Issuer string

	// Audience identifies the intended recipient(s) of the token.
	Audience string

	// CookieSecure determines whether authentication cookies
	// should only be transmitted over HTTPS.
	CookieSecure bool

	// CookieHTTPOnly prevents JavaScript access to authentication cookies.
	CookieHTTPOnly bool

	// CookieSameSite specifies the SameSite policy for cookies.
	CookieSameSite string
}

// loadAuthConfig loads authentication configuration from the environment.
func loadAuthConfig() AuthConfig {
	return AuthConfig{
		JWTSecret:          getEnv("JWT_SECRET", "change-me-before-production"),
		AccessTokenExpiry:  getEnvAsDuration("ACCESS_TOKEN_EXPIRY", "15m"),
		RefreshTokenExpiry: getEnvAsDuration("REFRESH_TOKEN_EXPIRY", "168h"),
		Issuer:             getEnv("JWT_ISSUER", "CampusCore"),
		Audience:           getEnv("JWT_AUDIENCE", "CampusCore"),
		CookieSecure:       getEnvAsBool("COOKIE_SECURE", false),
		CookieHTTPOnly:     getEnvAsBool("COOKIE_HTTP_ONLY", true),
		CookieSameSite:     getEnv("COOKIE_SAME_SITE", "Lax"),
	}
}
