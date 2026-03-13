package usecase

import (
	"errors"
	"strings"

	"github.com/philipos/prepbase/domain"
)

type userUsecase struct {
	userRepo domain.UserRepository
	pwdSvc   domain.PasswordService
	jwtSvc   domain.JWTService
}

func NewUserUsecase(r domain.UserRepository, p domain.PasswordService, j domain.JWTService) domain.UserUsecase {
	return &userUsecase{
		userRepo: r,
		pwdSvc:   p,
		jwtSvc:   j,
	}
}

// REGISTRATION & LOGIN

func (u *userUsecase) Register(user *domain.User) error {
	// Validaiton 
	if strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.Password) == "" {
		return errors.New("email and password cannot be empty")
	}

	existingUser, _ := u.userRepo.GetByEmail(user.Email)
	if existingUser != nil {
		return errors.New("email is already registered")
	}

	// Hash Password
	hashedPwd, err := u.pwdSvc.HashPassword(user.Password)
	if err != nil {
		return errors.New("failed to secure password")
	}
	user.Password = hashedPwd

	if user.Role == "" {
		user.Role = "user"
	}

	// Save to DB
	return u.userRepo.Create(user)
}

func (u *userUsecase) Login(email string, password string) (string, string, error) {
	// Find User
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid email or password")
	}

	// Soft Delete Check: Prevent deactivated users from logging in!
	if user.DeletedAt != nil {
		return "", "", errors.New("this account has been deactivated")
	}

	// Verify Password
	err = u.pwdSvc.ComparePassword(user.Password, password)
	if err != nil {
		return "", "", errors.New("invalid email or password")
	}

	// Generate Tokens
	accessToken, err := u.jwtSvc.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}

	refreshToken, err := u.jwtSvc.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	return accessToken, refreshToken, nil
}

// PROFILE MANAGEMENT

func (u *userUsecase) GetProfile(id string) (*domain.User, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) UpdateProfile(user *domain.User) error {
	// Security: We only allow updating Bio and Username.
	// Password changes require a specific flow (Forgot Password).
	// if strings.TrimSpace(user.Username) == "" {
	// 	return errors.New("username cannot be empty")
	// }
	return u.userRepo.Update(user)
}

// ADMIN ACTIONS & SOFT DELETE

func (u *userUsecase) PromoteUser(adminID string, targetUserID string) error {
	adminUser, err := u.userRepo.GetByID(adminID)
	if err != nil || adminUser.Role != "admin" {
		return errors.New("unauthorized: only admins can promote users")
	}

	// Find target user
	targetUser, err := u.userRepo.GetByID(targetUserID)
	if err != nil {
		return errors.New("target user not found")
	}

	// Update Role
	targetUser.Role = "admin"
	
	return u.userRepo.Update(targetUser)
}

func (u *userUsecase) DeleteUser(actorID string, actorRole string, targetID string) error {
	// Authorization Rule:
	// A user can delete themselves OR an Admin can delete anyone.
	if actorID != targetID && actorRole != "admin" {
		return errors.New("unauthorized: you can only delete your own account")
	}

	// Call the Soft Delete in the repository
	return u.userRepo.Delete(targetID)
}