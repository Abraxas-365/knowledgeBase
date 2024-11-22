package chatuser

import "context"

type Repository interface {
	GetChatUserByID(ctx context.Context, chatUserID string) (*ChatUser, error)
	CreateChatUser(ctx context.Context, u ChatUser) (*ChatUser, error)
}
