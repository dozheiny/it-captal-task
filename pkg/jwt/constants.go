package jwt

import "time"

const (
	accessTokenLifeTime = time.Hour * 5
	accessTokenExpired  = "access token expired"
	accessTokenIsWrong  = "access token is wrong"
	internalServerError = "internal server error"
)
