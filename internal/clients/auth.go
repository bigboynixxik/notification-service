package clients

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	auth "notification-service/pkg/api/auth/v1"
)

type AuthClient struct {
	client auth.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthClient(addr string) (*AuthClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("clients.NewAuthClient, failed to connect: %w", err)
	}
	client := auth.NewAuthServiceClient(conn)
	return &AuthClient{client: client, conn: conn}, nil
}

func (c *AuthClient) Close() error {
	return c.conn.Close()
}

func (c *AuthClient) BindTelegram(ctx context.Context, token string, chatID int64) (*auth.BindTelegramResponse, error) {
	resp, err := c.client.BindTelegram(ctx, &auth.BindTelegramRequest{
		Token:  token,
		ChatId: chatID,
	})
	if err != nil {
		return nil, fmt.Errorf("auth client: failed to bind: %w", err)
	}
	return resp, nil
}

func (c *AuthClient) GetUserChatID(ctx context.Context, userID string) (*auth.GetUserChatIDResponse, error) {
	resp, err := c.client.GetUserChatID(ctx, &auth.GetUserChatIDRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("auth client: failed to get chat id: %w", err)
	}
	return resp, nil
}
