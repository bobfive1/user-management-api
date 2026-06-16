package userprofile

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type UserProfileMock struct {
	mock.Mock
}

func NewUserProfileMockMock() *UserProfileMock {
	return &UserProfileMock{}
}

func (m *UserProfileMock) Create(ctx context.Context, request UserProfile) (UserProfile, error) {
	args := m.Called()
	return args.Get(0).(UserProfile), args.Error(1)
}

func (m *UserProfileMock) List(ctx context.Context) ([]UserProfile, error) {
	args := m.Called()
	return args.Get(0).([]UserProfile), args.Error(1)
}

func (m *UserProfileMock) GetByID(ctx context.Context, userID string) (UserProfile, error) {
	args := m.Called()
	return args.Get(0).(UserProfile), args.Error(1)
}

func (m *UserProfileMock) Update(ctx context.Context, userID string, request UserProfile) (UserProfile, error) {
	args := m.Called()
	return args.Get(0).(UserProfile), args.Error(1)
}

func (m *UserProfileMock) Delete(ctx context.Context, userID string) error {
	args := m.Called()
	return args.Error(0)
}
