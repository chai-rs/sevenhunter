package dto

type LoginReq struct {
	Email    string `json:"email" validate:"required,email,min=5,max=200"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type RegisterReq struct {
	Name     string `json:"name" validate:"required,min=2,max=32"`
	Email    string `json:"email" validate:"required,email,min=5,max=200"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type AuthResp struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	User         *UserResp `json:"user"`
}
