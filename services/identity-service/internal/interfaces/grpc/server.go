package grpc

import (
	"context"
	"time"

	"github.com/google/uuid"

	identityv1 "github.com/s7venking/ticket-box/identity/gen/identity/v1"
	"github.com/s7venking/ticket-box/identity/internal/application/command"
	"github.com/s7venking/ticket-box/identity/internal/application/dto"
	"github.com/s7venking/ticket-box/identity/internal/application/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *IdentityServer) DeactivateUser(
	ctx context.Context,
	req *identityv1.DeactivateUserRequest,
) (*identityv1.DeactivateUserResponse, error) {

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid user id",
		)
	}

	result, err := s.deactivateUserHandler.Handle(
		ctx,
		command.DeactivateUserCommand{
			UserID: id,
		},
	)

	if err != nil {
		return nil, mapError(err)
	}

	return &identityv1.DeactivateUserResponse{
		User: toProtoUser(result),
	}, nil
}

func (s *IdentityServer) ActivateUser(
	ctx context.Context,
	req *identityv1.ActivateUserRequest,
) (*identityv1.ActivateUserResponse, error) {

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid user id",
		)
	}

	result, err := s.activateUserHandler.Handle(
		ctx,
		command.ActivateUserCommand{
			UserID: id,
		},
	)

	if err != nil {
		return nil, mapError(err)
	}

	return &identityv1.ActivateUserResponse{
		User: toProtoUser(result),
	}, nil
}

func (s *IdentityServer) UpdateUserProfile(
	ctx context.Context,
	req *identityv1.UpdateUserProfileRequest,
) (*identityv1.UpdateUserProfileResponse, error) {

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid user id",
		)
	}

	result, err := s.updateProfileHandler.Handle(
		ctx,
		command.UpdateProfileCommand{
			UserID:   id,
			FullName: req.GetFullName(),
			Phone:    req.GetPhone(),
		},
	)

	if err != nil {
		return nil, mapError(err)
	}

	return &identityv1.UpdateUserProfileResponse{
		User: toProtoUser(result),
	}, nil
}

func (s *IdentityServer) GetUserById(
	ctx context.Context,
	req *identityv1.GetUserByIdRequest,
) (*identityv1.GetUserByIdResponse, error) {

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid user id",
		)
	}

	result, err := s.getUserByIDHandler.Handle(
		ctx,
		query.GetUserByIDQuery{
			ID: id,
		},
	)

	if err != nil {
		return nil, mapError(err)
	}

	return &identityv1.GetUserByIdResponse{
		User: toProtoUser(result),
	}, nil
}

func (s *IdentityServer) GetUserByEmail(
	ctx context.Context,
	req *identityv1.GetUserByEmailRequest,
) (*identityv1.GetUserByEmailResponse, error) {

	result, err := s.getUserByEmailHandler.Handle(
		ctx,
		query.GetUserByEmailQuery{
			Email: req.GetEmail(),
		},
	)

	if err != nil {
		return nil, mapError(err)
	}

	return &identityv1.GetUserByEmailResponse{
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
