package mcpadapter

// contextKey is an unexported type for context keys defined in this package.
type contextKey string

// TokenKey is the context key for the user's auth token.
const TokenKey = contextKey("UIT_AUTH_TOKEN")
