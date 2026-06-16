package userprofile

import (
	"context"
	"errors"

	"github.com/bobfive1/user-management-api/internal/logger"

	"go.uber.org/zap"
)

var (
	ErrorPassNotMatch = errors.New("Password not match")
)

type UserProfileService struct {
	logger *zap.SugaredLogger
	repo   UserProfileRepository
}

func NewUserProfileService(repo UserProfileRepository) UserProfileService {
	return UserProfileService{
		logger: logger.GetLogger("UserProfile service"),
		repo:   repo,
	}
}

func (s *UserProfileService) Create(ctx context.Context, request InsertUserProfileRequest) (UserProfileDisplay, error) {
	userProfile := UserProfile{
		UserID:    request.UserID,
		Password:  HashPassword(request.Password),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Address:   request.Address,
		Birthdate: request.Birthdate,
		Email:     request.Email,
	}

	createdUserProfile, err := s.repo.Create(ctx, userProfile)
	if err != nil {
		return UserProfileDisplay{}, err
	}
	return createdUserProfile, nil
}

func (s *UserProfileService) List(ctx context.Context) ([]UserProfileDisplay, error) {
	listUserProfile, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return listUserProfile, nil
}

func (s *UserProfileService) GetByID(ctx context.Context, userID string) (UserProfileDisplay, error) {
	userProfile, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return UserProfileDisplay{}, err
	}
	return userProfile, nil
}

func (s *UserProfileService) Update(ctx context.Context, userID string, request UpdateUserProfileRequest) (UserProfileDisplay, error) {
	userProfile := UserProfile{
		Password:  HashPassword(request.Password),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Address:   request.Address,
		Birthdate: request.Birthdate,
		Email:     request.Email,
	}

	updatedUserProfile, err := s.repo.Update(ctx, userID, userProfile)
	if err != nil {
		return UserProfileDisplay{}, err
	}
	return updatedUserProfile, nil
}

func (s *UserProfileService) Delete(ctx context.Context, userID string) error {
	err := s.repo.Delete(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserProfileService) Login(ctx context.Context, request UserProfileLoginRequest) (UserProfileDisplay, error) {
	user, err := s.repo.GetByIDWithPassword(ctx, request.UserID)
	if err != nil {
		return UserProfileDisplay{}, err
	}
	//check password
	if !CheckPasswordHash(request.Password, user.Password) {
		return UserProfileDisplay{}, ErrorPassNotMatch
	}
	return NewUserProfileDisplay(user), nil
}
