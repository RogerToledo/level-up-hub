package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/me/level-up-hub/internal/repository"
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

// CreateActivity mocks the CreateActivity method.
func (m *MockQuerier) CreateActivity(ctx context.Context, arg repository.CreateActivityParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

// FindActivityByID mocks the FindActivityByID method.
func (m *MockQuerier) FindActivityByID(ctx context.Context, arg repository.FindActivityByIDParams) (repository.Activity, error) {
	args := m.Called(ctx, arg)
	if args.Get(0) == nil {
		return repository.Activity{}, args.Error(1)
	}
	return args.Get(0).(repository.Activity), args.Error(1)
}

// UpdateActivityProgress mocks the UpdateActivityProgress method.
func (m *MockQuerier) UpdateActivityProgress(ctx context.Context, arg repository.UpdateActivityProgressParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

// DeleteActivity mocks the DeleteActivity method.
func (m *MockQuerier) DeleteActivity(ctx context.Context, arg repository.DeleteActivityParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

// Note: Additional mock methods can be added as needed for specific tests
// For now, only the essential methods for user account tests are implemented
