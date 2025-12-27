package models

// Visitor represents a visitor to a page
type Visitor struct {
	BaseModel `bson:",inline"`

	Fingerprint string `bson:"fingerprint" json:"fingerprint"`
	Domain      string `bson:"domain" json:"domain"`
	PageID      string `bson:"pageId" json:"pageId"`
}

// CollectionName returns the MongoDB collection name
func (v *Visitor) CollectionName() string {
	return "visitors"
}
