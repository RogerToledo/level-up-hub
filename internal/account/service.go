package account

import (
	"context"
	"fmt"

	"github.com/me/level-up-hub/apperr"
	"github.com/me/level-up-hub/auth"
	"github.com/me/level-up-hub/internal/repository"
	"github.com/me/level-up-hub/internal/pkg/identity"
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
		return apperr.MessageError(apperr.ErrEncryptPassword, err)
	}

	err = s.repo.CreateUser(ctx, repository.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Active:   req.Active,
	})
	if err != nil {
		return apperr.MessageError(fmt.Sprintf(apperr.ErrCreate, apperr.UserPT), err)
	}

	return nil
}

func (s *Service) UpdateUser(ctx context.Context, id string, req CreateUserRequest) error {
	userID, err := identity.ParseID(id)
	if err != nil {
		return apperr.MessageError(apperr.ErrInvalidID, err)
	}

	if _, err = s.repo.FindUserByID(ctx, userID); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return apperr.MessageError(apperr.ErrEncryptPassword, err)
	}

	err = s.repo.UpdateUser(ctx, repository.UpdateUserParams{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Active:   req.Active,
	})
	if err != nil {
		return apperr.MessageError(fmt.Sprintf(apperr.ErrorUpdate, apperr.UserPT), err)
	}

	return nil
}

func (s *Service) DeleteUser(ctx context.Context, id string) error {
	userID, err := identity.ParseID(id)
	if err != nil {
		return apperr.MessageError(apperr.ErrInvalidID, err)
	}

	if _, err = s.repo.FindUserByID(ctx, userID); err != nil {
		return err
	}

	err = s.repo.DeleteUser(ctx, userID)
	if err != nil {
		return apperr.MessageError(fmt.Sprintf(apperr.ErrDelete, apperr.UserPT), err)
	}

	return nil
}

func (s *Service) FindUserByID(ctx context.Context, id string) (*repository.FindUserByIDRow, error) {
	userID, err := identity.ParseID(id)
	if err != nil {
		return nil, apperr.MessageError(apperr.ErrInvalidID, err)
	}

	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, apperr.MessageError(fmt.Sprintf(apperr.ErrFindByID, apperr.UserPT), err)
	}

	return &user, nil
}

func (s *Service) FindUserByEmail(ctx context.Context, email string) (*repository.FindUserByEmailRow, error) {
	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, apperr.MessageError(fmt.Sprintf(apperr.ErrFindByID, apperr.UserPT), err)
	}

	return &user, nil
}

func (s *Service) FindAllUsers(ctx context.Context) ([]repository.FindAllUsersRow, error) {
	users, err := s.repo.FindAllUsers(ctx)
	if err != nil {
		return nil, apperr.MessageError(fmt.Sprintf(apperr.ErrFindAll, apperr.UserPT), err)
	}

	return users, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest, secret string) (LoginResponse, error) {
	user, err := s.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return LoginResponse{}, apperr.MessageError(fmt.Sprintf(apperr.ErrFindByID, apperr.UserPT), err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return LoginResponse{}, apperr.MessageError(apperr.ErrInvalidCredentials, err)
	}

	token, err := auth.GenerateToken(user.ID, user.Role, secret)
	if err != nil {
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
