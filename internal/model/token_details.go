package model

type TokenDetails struct {
	AccessToken  string `bson:"access_token"`
	RefreshToken string `bson:"refresh_token"`
	AtExpires    int64  `bson:"at_expires"`
	RtExpires    int64  `bson:"rt_expires"`
}
