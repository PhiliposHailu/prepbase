package domain

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Bio       string    `json:"bio"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// If DeletedAt is NOT nil, the user is soft-deleted.
	// Anonymization Logic: When fetching Questions or Comments authored by this user,
	// the presentation layer should check this field (or the Usecase should mask it)
	// and display "[Anonymous]" instead of the actual Username to protect privacy, 
	// while preserving the data in the DB for audit/ban purposes.
	DeletedAt *time.Time `json:"-"`
}

type UserRepository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
	GetByID(id string) (*User, error)
	Update(user *User) error
}
type UserUsecase interface {
	Register(user *User) error
	Login(email string, password string) (accessToken string, refreshToken string, err error)
	GetProfile(id string) (*User, error)
	UpdateProfile(user *User) error
	PromoteUser(adminID string, targetUserID string) error
	DeleteUser(actorID string, actorRole string, targetID string) error
}