package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/workspaces", func(ctx *gin.Context) {
		userID := ctx.Query("user_id")

		if userID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    "MISSING_USER_ID",
				"message": "user_id is required",
			})
			return
		}

		if userID == "189512da-360e-4e68-85bc-3d1b4f80b68f" {
			ctx.JSON(http.StatusOK, gin.H{
				"message":    "OK",
				"workspaces": []string{"ws-1", "ws-2"},
			})
			return
		}

		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    "USER_NOT_FOUND",
			"message": "not found user with provided id",
		})
	})

	r.Run(":10012")
}
