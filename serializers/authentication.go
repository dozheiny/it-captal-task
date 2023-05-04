package serializers

import "errors"

var (
	userNameCannotEmpty     = errors.New("username cannot empty")
	passwordCannotEmpty     = errors.New("password cannot empty")
	refreshTokenCannotEmpty = errors.New("refresh token cannot empty")
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

func (i *Login) Validation() error {
	if len(i.Username) == 0 {
		return userNameCannotEmpty
	}

	if len(i.Password) == 0 {
		return passwordCannotEmpty
	}

	return nil
}

func (i *RefreshToken) Validation() error {
	if len(i.RefreshToken) == 0 {
		return refreshTokenCannotEmpty
	}

	return nil
}
