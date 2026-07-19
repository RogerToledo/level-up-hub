package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/me/level-up-hub/backend/internal/repository"
	"github.com/stretchr/testify/mock"
)

// MockQuerier is a mock implementation of repository.Querier
type MockQuerier struct {
	mock.Mock
}

// CreateUser mocks the CreateUser method.
func (m *MockQuerier) CreateUser(ctx context.Context, arg repository.CreateUserParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

// FindUserByID mocks the FindUserByID method.
func (m *MockQuerier) FindUserByID(ctx context.Context, id uuid.UUID) (repository.FindUserByIDRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(repository.FindUserByIDRow), args.Error(1)
}

// FindUserByEmail mocks the FindUserByEmail method.
func (m *MockQuerier) FindUserByEmail(ctx context.Context, email string) (repository.FindUserByEmailRow, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return repository.FindUserByEmailRow{}, args.Error(1)
	}
	return args.Get(0).(repository.FindUserByEmailRow), args.Error(1)
}

// FindAllUsers mocks the FindAllUsers method.
func (m *MockQuerier) FindAllUsers(ctx context.Context) ([]repository.FindAllUsersRow, error) {
	args := m.Called(ctx)
	return args.Get(0).([]repository.FindAllUsersRow), args.Error(1)
}

// FindAllUsersPaginated mocks the FindAllUsersPaginated method.
func (m *MockQuerier) FindAllUsersPaginated(ctx context.Context, arg repository.FindAllUsersPaginatedParams) ([]repository.FindAllUsersPaginatedRow, error) {
	args := m.Called(ctx, arg)
	if args.Get(0) == nil {
		return []repository.FindAllUsersPaginatedRow{}, args.Error(1)
	}
	return args.Get(0).([]repository.FindAllUsersPaginatedRow), args.Error(1)
}

// CountAllUsers mocks the CountAllUsers method.
func (m *MockQuerier) CountAllUsers(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// UpdateUser mocks the UpdateUser method.
func (m *MockQuerier) UpdateUser(ctx context.Context, arg repository.UpdateUserParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

// DeleteUser mocks the DeleteUser method.
func (m *MockQuerier) DeleteUser(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// CreateInitiative mocks the CreateInitiative method.
func (m *MockQuerier) CreateInitiative(ctx context.Context, arg repository.CreateInitiativeParams) (repository.Initiative, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repository.Initiative), args.Error(1)
}

// FindInitiativeByID mocks the FindInitiativeByID method.
func (m *MockQuerier) FindInitiativeByID(ctx context.Context, arg repository.FindInitiativeByIDParams) (repository.FindInitiativeByIDRow, error) {
	args := m.Called(ctx, arg)
	if args.Get(0) == nil {
		return repository.FindInitiativeByIDRow{}, args.Error(1)
	}
	return args.Get(0).(repository.FindInitiativeByIDRow), args.Error(1)
}

// DeleteInitiative mocks the DeleteInitiative method.
func (m *MockQuerier) DeleteInitiative(ctx context.Context, arg repository.DeleteInitiativeParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

// Note: Additional mock methods can be added as needed for specific tests
// For now, only the essential methods for user account tests are implemented
