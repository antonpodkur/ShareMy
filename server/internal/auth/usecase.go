package auth

import "github.com/antonpodkur/ShareMy/internal/models"

type Usecase interface {
    SignUp(*models.SignUpInput) (*models.UserDBResponse, error)
    SignIn(*models.SignInInput) (*models.UserDBResponse, error)
    GetUserById(string) (*models.UserDBResponse, error)
    GetUserByEmail(string) (*models.UserDBResponse, error)
}
