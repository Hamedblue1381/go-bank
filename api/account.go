package api

import "github.com/gin-gonic/gin"

type createAccountRequest struct {
	Owner   string `json:"owner" binding:"required"`
	Curreny string `json:"currency" binding:"required,oneof=USD EUR IRL"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	ctx.ShouldBindJSON(&req)
}
