package domain

type User struct {
	ID               string `json:"id" bson:"_id"`
	Username         string `json:"username" bson:"username"`
	Region           string `json:"region" bson:"region"`
	OrganizationType string `json:"organization_type" bson:"organization_type"`
	Role             string `json:"role" bson:"role"`
	Contacts         struct {
		Email       string `json:"email" bson:"email"`
		PhoneNumber string `json:"phone_number" bson:"phone_number"`
		TelegramID  string `json:"telegram_id" bson:"telegram_id"`
	} `json:"contacts" bson:"contacts"`
	Permissions struct {
		Email    bool `json:"email" bson:"email"`
		Phone    bool `json:"phone" bson:"phone"`
		Telegram bool `json:"telegram" bson:"telegram"`
		Push     bool `json:"push" bson:"push"`
	} `json:"permissions" bson:"permissions"`
}
