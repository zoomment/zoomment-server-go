package constants

const (
	// Reaction limits
	MaxReactionLength = 20

	// JWT expiration (1 year in hours)
	JWTExpirationHours = 365 * 24

	// Validation limits
	MaxPageIDLength = 500
	MaxDomainLength = 253
	MaxEmailLength  = 254

	// Date format for email notifications
	DateFormat = "02 Jan 2006 - 15:04"
)
