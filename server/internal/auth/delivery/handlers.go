package delivery

import (
	"net/http"
	"strings"

	"github.com/antonpodkur/ShareMy/config"
	"github.com/antonpodkur/ShareMy/internal/auth"
	"github.com/antonpodkur/ShareMy/internal/models"
	"github.com/gin-gonic/gin"
)

type authHandlers struct {
    cfg *config.Config
    authUsecase auth.Usecase 
}

func NewAuthHandlers(cfg *config.Config, authUsecase auth.Usecase) auth.Handlers {
    return &authHandlers {
        cfg: cfg,
        authUsecase: authUsecase,
    }
}

func (ah *authHandlers) SignUp() gin.HandlerFunc {
    return func(c *gin.Context) {
       var user *models.SignUpInput 

       if err := c.ShouldBindJSON(&user); err != nil {
           c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
           return
       }

       if user.Password != user.PasswordConfirm {
           c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
           return
       }

       newUser, err := ah.authUsecase.SignUp(user)
       if err != nil {
           if strings.Contains(err.Error(), "email already exists") {
               c.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
               return
           }
           c.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
           return
       }

       c.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": models.UserFilteredResponse(newUser)}})
    }
}

func (ah *authHandlers) SignIn() gin.HandlerFunc {
	panic("not implemented") // TODO: Implement
}

func (ah *authHandlers) RefreshAccessToken() gin.HandlerFunc {
	panic("not implemented") // TODO: Implement
}

func (ah *authHandlers) LogOut() gin.HandlerFunc {
	panic("not implemented") // TODO: Implement
}

func (ah *authHandlers) GetMe() gin.HandlerFunc {
	panic("not implemented") // TODO: Implement
}
