package middlewares

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(jwksURL, issuer, audience string) gin.HandlerFunc {
	keySet := jwkCache(jwksURL)

	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		if tokenStr == auth {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid auth header"})
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}

			kid, ok := t.Header["kid"].(string)
			if !ok {
				return nil, fmt.Errorf("missing kid")
			}

			return keySet.GetKey(kid)
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid claims"})
			return
		}

		// validate issuer
		if claims["iss"] != issuer {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid issuer"})
			return
		}

		// validate audience
		if !containsAudience(claims["aud"], audience) {
			c.AbortWithStatusJSON(403, gin.H{"error": "invalid audience"})
			return
		}

		// attach identity
		c.Set("sub", claims["sub"])
		c.Set("claims", claims)

		c.Next()
	}
}

func containsAudience(tokenAud any, requiredAud string) bool {
	aud, ok := tokenAud.(string)
	if !ok {
		return false
	}

	return aud == requiredAud
}
