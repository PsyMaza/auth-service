package models

import "time"

type TokenDetails struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	AtExpires    time.Time `json:"-"`
	RtExpires    time.Time `json:"-"`
}
