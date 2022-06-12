package models

import "time"

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    time.Time
	RtExpires    time.Time
}
