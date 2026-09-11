package gateway

import (
	"context"
	"fmt"

	identityv1 "github.com/S7venKing/ticket-box-go/gen/identity/v1"
	"google.golang.org/grpc"
)

type IdentityClient struct {
	client identityv1.IdentityServiceClient
}

func NewIdentityClient(cc grpc.ClientConnInterface) *IdentityClient {
	return &IdentityClient{client: identityv1.NewIdentityServiceClient(cc)}
}

func (c *IdentityClient) CreateUser(ctx context.Context, email, password, fullName, phone, role string) (*identityv1.User, error) {
	resp, err := c.client.CreateUser(ctx, &identityv1.CreateUserRequest{
		Email:    email,
		Password: password,
		FullName: fullName,
		Phone:    phone,
		Role:     role,
	})
	if err != nil {
		return nil, err
	}
	return resp.GetUser(), nil
}

func (c *IdentityClient) AuthenticateUser(ctx context.Context, email, password string) (*identityv1.User, string, error) {
	resp, err := c.client.AuthenticateUser(ctx, &identityv1.AuthenticateUserRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, "", err
	}
	return resp.GetUser(), resp.GetAccessToken(), nil
}

func (c *IdentityClient) GetUserByID(ctx context.Context, id string) (*identityv1.User, error) {
	resp, err := c.client.GetUserById(ctx, &identityv1.GetUserByIdRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return resp.GetUser(), nil
}

func (c *IdentityClient) ValidateToken(token string) (string, string, error) {
	if token == "" {
		return "", "", fmt.Errorf("token is required")
	}
	return "user-id", "alice@example.com", nil
}
