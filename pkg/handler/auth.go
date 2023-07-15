package handler

import (
	todo "github.com/LittleMikle/ToDo_List"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var input todo.User
	err := c.BindJSON(&input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		//log.Error().Msgf("failed with signUp handler: %w", err)
		return
	}
	id, err := h.services.Autorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	err := c.BindJSON(&input)
	if err != nil {
		log.Error().Msgf("failed with signIn JSON input: %w", err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		//log.Error().Msgf("failed with signIn handler: %w", err)
	}
	token, err := h.services.Autorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		log.Error().Msgf("failed with signIn generating token: %w", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": token,
	})
}
