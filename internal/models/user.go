package models

type User struct {
	Base
	Name     string `json:"name"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"-"`
	Image    string `json:"path_to_image"`
}

type RegisterUser struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required,min=3,max=50"`
	Password string `json:"password"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"`
}

type Token struct {
	User         User
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
