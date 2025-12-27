package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseModel provides ID and timestamps matching Node.js mongoose format
// Embed this in all models instead of mgm.DefaultModel
type BaseModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// PrepareID prepares the ID before creating (implements mgm.Model interface)
func (m *BaseModel) PrepareID(id interface{}) (interface{}, error) {
	if idStr, ok := id.(string); ok && idStr == "" {
		return nil, nil
	}
	return id, nil
}

// GetID returns the model ID (implements mgm.Model interface)
func (m *BaseModel) GetID() interface{} {
	return m.ID
}

// SetID sets the model ID (implements mgm.Model interface)
func (m *BaseModel) SetID(id interface{}) {
	if oid, ok := id.(primitive.ObjectID); ok {
		m.ID = oid
	}
}

// Creating sets timestamps when creating a new document (implements mgm.CreatingHook)
func (m *BaseModel) Creating() error {
	if m.ID.IsZero() {
		m.ID = primitive.NewObjectID()
	}
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return nil
}

// Saving updates the timestamp when saving (implements mgm.SavingHook)
func (m *BaseModel) Saving() error {
	m.UpdatedAt = time.Now()
	return nil
}

