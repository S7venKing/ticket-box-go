package grpc

import (
	"context"
	"time"

	identityv1 "github.com/s7venking/ticket-box/identity/gen/identity/v1"
	"github.com/s7venking/ticket-box/identity/internal/application/command"
	"github.com/s7venking/ticket-box/identity/internal/application/dto"
	"github.com/s7venking/ticket-box/identity/internal/application/query"
)

type IdentityServer struct {
	identityv1.UnimplementedIdentityServiceServer

	createUserHandler     *command.CreateUserHandler
	getUserByIDHandler    *query.GetUserByIDHandler
	getUserByEmailHandler *query.GetUserByEmailHandler
	updateProfileHandler  *command.UpdateProfileHandler
	activateUserHandler   *command.ActivateUserHandler
	deactivateUserHandler *command.DeactivateUserHandler
}

func NewIdentityServer(
	createUserHandler *command.CreateUserHandler,
	getUserByIDHandler *query.GetUserByIDHandler,
	getUserByEmailHandler *query.GetUserByEmailHandler,
	updateProfileHandler *command.UpdateProfileHandler,
	activateUserHandler *command.ActivateUserHandler,
	deactivateUserHandler *command.DeactivateUserHandler,
) *IdentityServer {
	return &IdentityServer{
		createUserHandler:     createUserHandler,
		getUserByIDHandler:    getUserByIDHandler,
		getUserByEmailHandler: getUserByEmailHandler,
		updateProfileHandler:  updateProfileHandler,
		activateUserHandler:   activateUserHandler,
		deactivateUserHandler: deactivateUserHandler,
	}
}

func (s *IdentityServer) CreateUser(
	ctx context.Context,
	req *identityv1.CreateUserRequest,
) (*identityv1.CreateUserResponse, error) {

	result, err := s.createUserHandler.Handle(
		ctx,
		command.CreateUserCommand{
			Email:    req.GetEmail(),
			Password: req.GetPassword(),
			FullName: req.GetFullName(),
			Phone:    req.GetPhone(),
		},
	)

	if err != nil {
		return nil, mapError(err)
	}

	return &identityv1.CreateUserResponse{
		User: toProtoUser(result),
	}, nil
}

func toProtoUser(
	u *dto.UserDTO,
) *identityv1.User {
	return &identityv1.User{
		Id:        u.ID.String(),
		Email:     u.Email,
		FullName:  u.FullName,
		Phone:     u.Phone,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339Nano),
	}
}
