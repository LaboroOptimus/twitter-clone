package authjwt

import "context"

type ctxKey string

const (
	KeyUserID   ctxKey = "userID"
	KeyNickname ctxKey = "nickname"
	KeyJTI      ctxKey = "jti"
)

func WithUserID(ctx context.Context, id uint) context.Context {
	return context.WithValue(ctx, KeyUserID, id)
}
func UserIDFrom(ctx context.Context) (uint, bool) {
	v := ctx.Value(KeyUserID)
	id, ok := v.(uint)
	return id, ok
}

func WithNickname(ctx context.Context, nick string) context.Context {
	return context.WithValue(ctx, KeyNickname, nick)
}
func NicknameFrom(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(KeyNickname).(string)
	return v, ok
}

func WithJTI(ctx context.Context, jti string) context.Context {
	return context.WithValue(ctx, KeyJTI, jti)
}
func JTIFrom(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(KeyJTI).(string)
	return v, ok
}
