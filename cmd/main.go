package main

import (
	"net/http"
	"xplr-opal/api/generated"
	"xplr-opal/infra/opa"
	"xplr-opal/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type App struct {
	opaClient opa.OPAClient
}

func (a *App) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func (a *App) ListTodos(c *gin.Context) {
	resourceOwnerID := "189512da-360e-4e68-85bc-3d1b4f80b68f"
	input := map[string]any{
		"user": map[string]any{
			"id":    c.GetString("sub"),
			"roles": []string{"user"},
		},
		"action": "list",
		"resource": map[string]string{
			"type":     "todo",
			"owner_id": resourceOwnerID,
		},
	}

	allow, err := a.opaClient.Allow(c.Request.Context(), "todos", input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "AUTHZ_ERROR",
			"message": err.Error(),
		})
		return
	}

	if !allow {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    "FORBIDDEN",
			"message": "You do not have permission to access this resource",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          uuid.New(),
		"description": "homework",
	})
}

func main() {
	r := gin.Default()

	authorized := r.Group("/")
	authorized.Use(
		middlewares.JWTAuthMiddleware(
			"http://localhost:8080/realms/myrealm/protocol/openid-connect/certs",
			"http://localhost:8080/realms/myrealm",
			"account",
		),
	)

	opaClient := opa.NewOPAClient("http://localhost:8181")
	app := &App{
		opaClient: opaClient,
	}
	generated.RegisterHandlers(authorized, app)

	r.Run(":10011")
}
