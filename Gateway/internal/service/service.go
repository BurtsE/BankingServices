package service

type IUserService interface {
	Validate(jwtToken string) (uuid string, err error)
}
