package serializers

import "errors"

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var (
	userNameCannotEmpty = errors.New("username cannot empty")
	passwordCannotEmpty = errors.New("password cannot empty")
)

func (i *Login) Validation() error {
	if len(i.Username) == 0 {
		return userNameCannotEmpty
	}

	if len(i.Password) == 0 {
		return passwordCannotEmpty
	}

	return nil
}
