package models

// User represents a user in the system
type User struct {
	BaseModel `bson:",inline"`

	Name       string `bson:"name" json:"name"`
	Email      string `bson:"email" json:"email"`
	Role       int    `bson:"role" json:"role"`
	IsVerified bool   `bson:"isVerified" json:"isVerified"`
}

// NewUser creates a new user with default values
func NewUser(email string) *User {
	return &User{
		Email:      email,
		Role:       RoleAdmin,
		IsVerified: false,
	}
}

// CollectionName returns the MongoDB collection name
func (u *User) CollectionName() string {
	return "users"
}

// Constants for user roles
const (
	RoleAdmin      = 1
	RoleSuperAdmin = 2
)
