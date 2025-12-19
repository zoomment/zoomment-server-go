package models

import (
	"github.com/kamva/mgm/v3"
)

// Visitor represents a visitor to a page
type Visitor struct {
	mgm.DefaultModel `bson:",inline"`

	Fingerprint string `bson:"fingerprint" json:"fingerprint"`
	Domain      string `bson:"domain" json:"domain"`
	PageID      string `bson:"pageId" json:"pageId"`
}

// CollectionName returns the MongoDB collection name
func (v *Visitor) CollectionName() string {
	return "visitors"
}

// VisitorCountResponse is the response format for getting visitor count
type VisitorCountResponse struct {
	Count int `json:"count"`
}
