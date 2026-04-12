package account

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/me/level-up-hub/apperr"
	"github.com/me/level-up-hub/auth"
	"github.com/me/level-up-hub/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *repository.Queries
}

func NewService(repo *repository.Queries) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("password encryption failed",
            slog.String("error", err.Error()),
        )
		return apperr.MessageError(apperr.ErrEncryptPassword, err)
	}

	err = s.repo.CreateUser(ctx, repository.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Active:   req.Active,
	})
	if err != nil {
		 slog.Error("user creation failed",
            slog.String("error", err.Error()),
            slog.String("email", req.Email),
        )
		return apperr.MessageError(fmt.Sprintf(apperr.ErrCreate, apperr.UserPT), err)
	}

	return nil
}

func (s *Service) UpdateUser(ctx context.Context, id uuid.UUID, req CreateUserRequest) error {
	if _, err := s.repo.FindUserByID(ctx, id); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("password encryption failed",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		return apperr.MessageError(apperr.ErrEncryptPassword, err)
	}

	err = s.repo.UpdateUser(ctx, repository.UpdateUserParams{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Active:   req.Active,
	})
	if err != nil {
		slog.Error("user update failed",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		return apperr.MessageError(fmt.Sprintf(apperr.ErrorUpdate, apperr.UserPT), err)
	}

	return nil
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.FindUserByID(ctx, id); err != nil {
		return err
	}

	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		slog.Error("user deletion failed",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		return apperr.MessageError(fmt.Sprintf(apperr.ErrDelete, apperr.UserPT), err)
	}

	return nil
}

func (s *Service) FindUserByID(ctx context.Context, id uuid.UUID) (*repository.FindUserByIDRow, error) {
	user, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		slog.Error("failed to find user by ID",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		return nil, apperr.MessageError(fmt.Sprintf(apperr.ErrFindByID, apperr.UserPT), err)
	}

	return &user, nil
}

func (s *Service) FindUserByEmail(ctx context.Context, email string) (*repository.FindUserByEmailRow, error) {
	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		slog.Error("failed to find user by email",
			slog.String("error", err.Error()),
			slog.String("email", email),
		)
		return nil, apperr.MessageError(fmt.Sprintf(apperr.ErrFindByID, apperr.UserPT), err)
	}

	return &user, nil
}

func (s *Service) FindAllUsers(ctx context.Context) ([]repository.FindAllUsersRow, error) {
	users, err := s.repo.FindAllUsers(ctx)
	if err != nil {
		slog.Error("failed to find all users",
			slog.String("error", err.Error()),
		)
		return nil, apperr.MessageError(fmt.Sprintf(apperr.ErrFindAll, apperr.UserPT), err)
	}

	return users, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest, secret string) (LoginResponse, error) {
	user, err := s.FindUserByEmail(ctx, req.Email)
	if err != nil {
		slog.Warn("login failed: user not found", slog.String("email", req.Email))
		return LoginResponse{}, apperr.MessageError(fmt.Sprintf(apperr.ErrFindByID, apperr.UserPT), err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		slog.Warn("login failed - invalid password",
            slog.String("email", req.Email),
            slog.String("user_id", user.ID.String()),
        )
		return LoginResponse{}, apperr.MessageError(apperr.ErrInvalidCredentials, err)
	}

	token, err := auth.GenerateToken(user.ID, string(user.Role), secret)
	if err != nil {
		 slog.Error("failed to generate token",
            slog.String("error", err.Error()),
            slog.String("user_id", user.ID.String()),
        )
		return LoginResponse{}, apperr.MessageError(apperr.ErrGenerateToken, err)
	}

	return LoginResponse{
		Token: token,
		User: UserResponse{
			ID:       user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
			Active:   user.Active,
		},
	}, nil
}
