package structs

import "time"

type User struct {
	ID                       int                           `json:"id"`
	Username                 *string                       `json:"username"`
	Name                     string                        `json:"name"`
	Email                    string                        `json:"email"`
	ProfileImageURL          string                        `json:"profile_image_url"`
	Authentication           GeneralNSN                    `json:"authentication"`
	InsertedAt               time.Time                     `json:"inserted_at"`
	LastActive               *time.Time                    `json:"last_active"`
	AuthenticationStrategies *UserAuthenticationStrategies `json:"strategies,omitempty"`
}

type UserTS struct {
	ID                       *int                          `json:"id,omitempty"`
	Username                 *string                       `json:"username,omitempty"`
	Name                     *string                       `json:"name,omitempty"`
	Email                    *string                       `json:"email,omitempty"`
	ProfileImageURL          *string                       `json:"profile_image_url,omitempty"`
	Authentication           *GeneralNSN                   `json:"authentication,omitempty"`
	InsertedAt               *time.Time                    `json:"inserted_at,omitempty"`
	LastActive               *time.Time                    `json:"last_active,omitempty"`
	AuthenticationStrategies *UserAuthenticationStrategies `json:"strategies,omitempty"`
}

type UserAuthenticationStrategies struct {
	UserID   *int    `json:"user_id,omitempty"`
	GoogleID *string `json:"google_id"`
	Password *string `json:"password"`
}
