package models

import (
	"github.com/kamva/mgm/v3"
)

// User represents a user in the system
// mgm.DefaultModel provides:
//   - ID (ObjectID)
//   - CreatedAt (time.Time)
//   - UpdatedAt (time.Time)
// Just like Mongoose's { timestamps: true }!
type User struct {
	mgm.DefaultModel `bson:",inline"` // Embed the default fields

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
// This is like setting the model name in Mongoose: mongoose.model('User', schema)
func (u *User) CollectionName() string {
	return "users"
}

// Constants for user roles
const (
	RoleAdmin      = 1
	RoleSuperAdmin = 2
)
