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
	listUsersHandler        *query.ListUsersHandler
	updateProfileHandler    *command.UpdateProfileHandler
	activateUserHandler     *command.ActivateUserHandler
	deactivateUserHandler   *command.DeactivateUserHandler
	assignOrganizerHandler  *command.AssignUserToOrganizerHandler
	createOrganizerHandler  *command.CreateOrganizerHandler
	getOrganizerHandler     *query.GetOrganizerHandler
	listOrganizersHandler   *query.ListOrganizersHandler
	updateOrganizerHandler  *command.UpdateOrganizerHandler
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
	listUsersHandlers ...*query.ListUsersHandler,
) *IdentityServer {
	var listUsersHandler *query.ListUsersHandler
	if len(listUsersHandlers) > 0 {
		listUsersHandler = listUsersHandlers[0]
	}
	return &IdentityServer{
		logger:                  logger,
		jwtSecret:               jwtSecret,
		createUserHandler:       createUserHandler,
		authenticateUserHandler: authenticateUserHandler,
		getUserByIDHandler:      getUserByIDHandler,
		getUserByEmailHandler:   getUserByEmailHandler,
		listUsersHandler:        listUsersHandler,
		updateProfileHandler:    updateProfileHandler,
		activateUserHandler:     activateUserHandler,
		deactivateUserHandler:   deactivateUserHandler,
	}
}

func (s *IdentityServer) SetAssignOrganizerHandler(handler *command.AssignUserToOrganizerHandler) {
	s.assignOrganizerHandler = handler
}

func (s *IdentityServer) SetOrganizerHandlers(
	create *command.CreateOrganizerHandler,
	get *query.GetOrganizerHandler,
	list *query.ListOrganizersHandler,
	update *command.UpdateOrganizerHandler,
) {
	s.createOrganizerHandler = create
	s.getOrganizerHandler = get
	s.listOrganizersHandler = list
	s.updateOrganizerHandler = update
}

func (s *IdentityServer) AssignUserToOrganizer(ctx context.Context, req *identityv1.AssignUserToOrganizerRequest) (*identityv1.AssignUserToOrganizerResponse, error) {
	if s.assignOrganizerHandler == nil {
		return nil, status.Error(codes.Unimplemented, "membership management unavailable")
	}
	userID, err := parseUserID(req.GetUserId())
	if err != nil {
		return nil, err
	}
	organizerID, err := uuid.Parse(req.GetOrganizerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid organizer id")
	}
	result, err := s.assignOrganizerHandler.Handle(ctx, command.AssignUserToOrganizerCommand{UserID: userID, OrganizerID: organizerID})
	if err != nil {
		return nil, s.mapError(ctx, err)
	}
	return &identityv1.AssignUserToOrganizerResponse{User: toProtoUser(result)}, nil
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
		Role:     req.GetRole(),
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

func (s *IdentityServer) ListUsers(
	ctx context.Context,
	req *identityv1.ListUsersRequest,
) (*identityv1.ListUsersResponse, error) {
	offset, limit := int(req.GetOffset()), int(req.GetLimit())
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	result, err := s.listUsersHandler.Handle(ctx, offset, limit)
	if err != nil {
		return nil, s.mapError(ctx, err)
	}
	users := make([]*identityv1.User, 0, len(result))
	for _, item := range result {
		users = append(users, toProtoUser(dto.FromUser(item)))
	}
	return &identityv1.ListUsersResponse{Users: users}, nil
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

func (s *IdentityServer) CreateOrganizer(ctx context.Context, req *identityv1.CreateOrganizerRequest) (*identityv1.CreateOrganizerResponse, error) {
	result, err := s.createOrganizerHandler.Handle(ctx, command.CreateOrganizerCommand{
		Name: req.GetName(), Email: req.GetEmail(), Phone: req.GetPhone(), Slug: req.GetSlug(),
	})
	if err != nil {
		return nil, s.mapError(ctx, err)
	}
	return &identityv1.CreateOrganizerResponse{Organizer: toProtoOrganizer(result)}, nil
}

func (s *IdentityServer) GetOrganizerById(ctx context.Context, req *identityv1.GetOrganizerByIdRequest) (*identityv1.GetOrganizerByIdResponse, error) {
	id, err := parseOrganizerID(req.GetId())
	if err != nil {
		return nil, err
	}
	result, err := s.getOrganizerHandler.Handle(ctx, id)
	if err != nil {
		return nil, s.mapError(ctx, err)
	}
	return &identityv1.GetOrganizerByIdResponse{Organizer: toProtoOrganizer(result)}, nil
}

func (s *IdentityServer) ListOrganizers(ctx context.Context, req *identityv1.ListOrganizersRequest) (*identityv1.ListOrganizersResponse, error) {
	offset, limit := int(req.GetOffset()), int(req.GetLimit())
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	result, err := s.listOrganizersHandler.Handle(ctx, offset, limit)
	if err != nil {
		return nil, s.mapError(ctx, err)
	}
	items := make([]*identityv1.Organizer, 0, len(result))
	for _, item := range result {
		items = append(items, toProtoOrganizer(item))
	}
	return &identityv1.ListOrganizersResponse{Organizers: items}, nil
}

func (s *IdentityServer) UpdateOrganizer(ctx context.Context, req *identityv1.UpdateOrganizerRequest) (*identityv1.UpdateOrganizerResponse, error) {
	id, err := parseOrganizerID(req.GetId())
	if err != nil {
		return nil, err
	}
	result, err := s.updateOrganizerHandler.Handle(ctx, command.UpdateOrganizerCommand{
		ID: id, Name: req.GetName(), Email: req.GetEmail(), Phone: req.GetPhone(), Slug: req.GetSlug(),
	})
	if err != nil {
		return nil, s.mapError(ctx, err)
	}
	return &identityv1.UpdateOrganizerResponse{Organizer: toProtoOrganizer(result)}, nil
}

func parseOrganizerID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, status.Error(codes.InvalidArgument, "invalid organizer id")
	}
	return id, nil
}

func parseUserID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	return id, nil
}

func toProtoOrganizer(o *dto.OrganizerDTO) *identityv1.Organizer {
	return &identityv1.Organizer{
		Id: o.ID.String(), Name: o.Name, Email: o.Email, Phone: o.Phone, Slug: o.Slug,
		IsActive: o.IsActive, CreatedAt: timestamppb.New(o.CreatedAt), UpdatedAt: timestamppb.New(o.UpdatedAt),
	}
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
		Role:      u.Role,
		OrganizerId: func() string {
			if u.OrganizerID == nil {
				return ""
			}
			return u.OrganizerID.String()
		}(),
	}
}
