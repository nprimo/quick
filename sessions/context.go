package sessions

import "context"

const (
	userIDKey    = "userID"
	csrfTokenKey = "csrfToken"
)

func WithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) int {
	if userID, ok := ctx.Value(userIDKey).(int); ok {
		return userID
	}
	return 0
}

func WithCSRFToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, csrfTokenKey, token)
}

func GetCSRFToken(ctx context.Context) string {
	if token, ok := ctx.Value(csrfTokenKey).(string); ok {
		return token
	}
	return ""
}
