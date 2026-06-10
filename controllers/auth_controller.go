package controllers

import (
	"max/auth/models"
	"max/auth/services"
	"max/auth/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var input models.AuthInput

	// Validar JSON entrante contra el struct AuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, "Datos inválidos", err.Error())
		return
	}

	if err := ctrl.service.RegisterUser(input); err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, "Error registrando usuario", nil)
		return
	}

	utils.RespondJSON(c, http.StatusCreated, "Usuario registrado exitosamente", nil)
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var input models.AuthInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, "Datos inválidos", err.Error())
		return
	}

	token, err := ctrl.service.LoginUser(input)
	if err != nil {
		utils.RespondJSON(c, http.StatusUnauthorized, "Error de autenticación", err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusOK, "Login exitoso", gin.H{"token": token})
}
