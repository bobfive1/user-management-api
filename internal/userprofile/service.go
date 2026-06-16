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

func (s *UserProfileService) Create(ctx context.Context, request InsertUserProfileRequest) (UserProfile, error) {
	userProfile := UserProfile{
		UserID:    request.UserID,
		Password:  HashPassword(request.Password),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Address:   request.Address,
		Birthdate: request.Birthdate,
		Email:     request.Email,
	}

	userProfile, err := s.repo.Create(ctx, userProfile)
	if err != nil {
		return UserProfile{}, err
	}
	return userProfile, nil
}

func (s *UserProfileService) List(ctx context.Context) ([]UserProfile, error) {
	listUserProfile, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return listUserProfile, nil
}

func (s *UserProfileService) GetByID(ctx context.Context, userID string) (UserProfile, error) {
	userProfile, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return UserProfile{}, err
	}
	return userProfile, nil
}

func (s *UserProfileService) Update(ctx context.Context, userID string, request UpdateUserProfileRequest) (UserProfile, error) {
	userProfile := UserProfile{
		Password:  HashPassword(request.Password),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Address:   request.Address,
		Birthdate: request.Birthdate,
		Email:     request.Email,
	}

	userProfile, err := s.repo.Update(ctx, userID, userProfile)
	if err != nil {
		return UserProfile{}, err
	}
	return userProfile, nil
}

func (s *UserProfileService) Delete(ctx context.Context, userID string) error {
	err := s.repo.Delete(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserProfileService) Login(ctx context.Context, request UserProfileLoginRequest) (UserProfile, error) {
	user, err := s.repo.GetByID(ctx, request.UserID)
	if err != nil {
		return UserProfile{}, err
	}
	//check password
	if !CheckPasswordHash(request.Password, user.Password) {
		return UserProfile{}, ErrorPassNotMatch
	}
	return user, nil
}
