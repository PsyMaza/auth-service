package response

// swagger:model TokenPair
type TokenPair struct {
	// AccessToken at
	AccessToken string `json:"accessToken"`
	// RefreshToken rt
	RefreshToken string `json:"refreshToken"`
}
