package grpc

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	identityv1 "github.com/S7venKing/ticket-box-go/gen/identity/v1"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/command"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/query"
)

type IdentityServer struct {
	identityv1.UnimplementedIdentityServiceServer

	logger    *slog.Logger
	jwtSecret string

	createUserHandler       *command.CreateUserHandler
	authenticateUserHandler *command.AuthenticateUserHandler
	getUserByIDHandler      *query.GetUserByIDHandler
	getUserByEmailHandler   *query.GetUserByEmailHandler
	updateProfileHandler    *command.UpdateProfileHandler
	activateUserHandler     *command.ActivateUserHandler
	deactivateUserHandler   *command.DeactivateUserHandler
}

func NewIdentityServer(
	logger *slog.Logger,
	jwtSecret string,
	createUserHandler *command.CreateUserHandler,
	authenticateUserHandler *command.AuthenticateUserHandler,
	getUserByIDHandler *query.GetUserByIDHandler,
	getUserByEmailHandler *query.GetUserByEmailHandler,
	updateProfileHandler *command.UpdateProfileHandler,
	activateUserHandler *command.ActivateUserHandler,
	deactivateUserHandler *command.DeactivateUserHandler,
) *IdentityServer {
	return &IdentityServer{
		logger:                  logger,
		jwtSecret:               jwtSecret,
		createUserHandler:       createUserHandler,
		authenticateUserHandler: authenticateUserHandler,
		getUserByIDHandler:      getUserByIDHandler,
		getUserByEmailHandler:   getUserByEmailHandler,
		updateProfileHandler:    updateProfileHandler,
		activateUserHandler:     activateUserHandler,
		deactivateUserHandler:   deactivateUserHandler,
	}
}

func (s *IdentityServer) CreateUser(
	ctx context.Context,
	req *identityv1.CreateUserRequest,
) (*identityv1.CreateUserResponse, error) {
	result, err := s.createUserHandler.Handle(ctx, command.CreateUserCommand{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		FullName: req.GetFullName(),
		Phone:    req.GetPhone(),
	})
	if err != nil {
		return nil, s.mapError(ctx, err)
	}

	return &identityv1.CreateUserResponse{User: toProtoUser(result)}, nil
}

func (s *IdentityServer) AuthenticateUser(
	ctx context.Context,
	req *identityv1.AuthenticateUserRequest,
) (*identityv1.AuthenticateUserResponse, error) {
	result, err := s.authenticateUserHandler.Handle(ctx, command.AuthenticateUserCommand{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, s.mapError(ctx, err)
	}

	token, err := generateJWT(s.jwtSecret, result.ID.String(), result.Email)
	if err != nil {
		return nil, s.mapError(ctx, err)
	}

	return &identityv1.AuthenticateUserResponse{
		User:        toProtoUser(result),
		AccessToken: token,
	}, nil
}

func (s *IdentityServer) GetUserById(
	ctx context.Context,
	req *identityv1.GetUserByIdRequest,
) (*identityv1.GetUserByIdResponse, error) {
	id, err := parseUserID(req.GetId())
	if err != nil {
		return nil, err
	}

	result, err := s.getUserByIDHandler.Handle(ctx, query.GetUserByIDQuery{ID: id})
	if err != nil {
		return nil, s.mapError(ctx, err)
	}

	return &identityv1.GetUserByIdResponse{User: toProtoUser(result)}, nil
}

func (s *IdentityServer) GetUserByEmail(
	ctx context.Context,
	req *identityv1.GetUserByEmailRequest,
) (*identityv1.GetUserByEmailResponse, error) {
	result, err := s.getUserByEmailHandler.Handle(ctx, query.GetUserByEmailQuery{Email: req.GetEmail()})
	if err != nil {
		return nil, s.mapError(ctx, err)
	}

	return &identityv1.GetUserByEmailResponse{User: toProtoUser(result)}, nil
}

func (s *IdentityServer) UpdateUserProfile(
	ctx context.Context,
	req *identityv1.UpdateUserProfileRequest,
) (*identityv1.UpdateUserProfileResponse, error) {
	id, err := parseUserID(req.GetId())
	if err != nil {
		return nil, err
	}

	result, err := s.updateProfileHandler.Handle(ctx, command.UpdateProfileCommand{
		UserID:   id,
		FullName: req.GetFullName(),
		Phone:    req.GetPhone(),
	})
	if err != nil {
		return nil, s.mapError(ctx, err)
	}

	return &identityv1.UpdateUserProfileResponse{User: toProtoUser(result)}, nil
}

func (s *IdentityServer) ActivateUser(
	ctx context.Context,
	req *identityv1.ActivateUserRequest,
) (*identityv1.ActivateUserResponse, error) {
	id, err := parseUserID(req.GetId())
	if err != nil {
		return nil, err
	}

	result, err := s.activateUserHandler.Handle(ctx, command.ActivateUserCommand{UserID: id})
	if err != nil {
		return nil, s.mapError(ctx, err)
	}

	return &identityv1.ActivateUserResponse{User: toProtoUser(result)}, nil
}

func (s *IdentityServer) DeactivateUser(
	ctx context.Context,
	req *identityv1.DeactivateUserRequest,
) (*identityv1.DeactivateUserResponse, error) {
	id, err := parseUserID(req.GetId())
	if err != nil {
		return nil, err
	}

	result, err := s.deactivateUserHandler.Handle(ctx, command.DeactivateUserCommand{UserID: id})
	if err != nil {
		return nil, s.mapError(ctx, err)
	}

	return &identityv1.DeactivateUserResponse{User: toProtoUser(result)}, nil
}

func parseUserID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	return id, nil
}

func toProtoUser(u *dto.UserDTO) *identityv1.User {
	return &identityv1.User{
		Id:        u.ID.String(),
		Email:     u.Email,
		FullName:  u.FullName,
		Phone:     u.Phone,
		IsActive:  u.IsActive,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}
