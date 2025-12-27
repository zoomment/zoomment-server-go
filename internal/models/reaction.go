package models

// Reaction represents a user's reaction (emoji) on a page
type Reaction struct {
	BaseModel `bson:",inline"`

	Fingerprint string `bson:"fingerprint" json:"fingerprint"`
	Domain      string `bson:"domain" json:"domain"`
	PageID      string `bson:"pageId" json:"pageId"`
	Reaction    string `bson:"reaction" json:"reaction"`
}

// CollectionName returns the MongoDB collection name
func (r *Reaction) CollectionName() string {
	return "reactions"
}

// ReactionAggregation represents the count of each reaction type
type ReactionAggregation struct {
	Reaction string `bson:"_id" json:"_id"`
	Count    int    `bson:"count" json:"count"`
}
