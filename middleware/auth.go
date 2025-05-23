package middleware

import (
	"coresamples/common"
	"coresamples/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Strips 'TOKEN ' prefix from token string
func stripBearerPrefixFromTokenString(tok string) (string, error) {
	// Should be a bearer token
	if len(tok) > 7 && strings.ToUpper(tok[0:7]) == "BEARER " {
		return tok[7:], nil
	}
	return tok, nil
}

// AuthorizationHeaderExtractor Extract  token from Authorization header
// Uses PostExtractionFilter to strip "TOKEN " prefix from header
var AuthorizationHeaderExtractor = &request.PostExtractionFilter{
	Extractor: request.HeaderExtractor{"Authorization"},
	Filter:    stripBearerPrefixFromTokenString,
}

// MyAuth2Extractor Extractor for OAuth2 access tokens.  Looks in 'Authorization'
// header then 'access_token' argument for a token.
var MyAuth2Extractor = &request.MultiExtractor{
	AuthorizationHeaderExtractor,
	request.ArgumentExtractor{"access_token"},
}

// UpdateContextUserModel Helper to write user info to the context
func UpdateContextUserModel(c *gin.Context, claims jwt.MapClaims) {
	c.Set("account_type", claims["role"])
	switch role := util.InterStringToString(claims["role"]); role {
	case "clinic":
		c.Set("account_id", claims["clinic_id"])
	case "customer":
		c.Set("account_id", claims["customer_id"])
	case "patient":
		c.Set("account_id", claims["patient_id"])
	default:
	}

}

func AuthMiddleware(auto401 bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.FullPath(), "healthcheck") {
			c.JSON(200, gin.H{
				"message": "pong",
			})
			c.Abort()
			return
		}
		if strings.Contains(c.FullPath(), "swagger") ||
			strings.Contains(c.FullPath(), "webhook") ||
			strings.Contains(c.FullPath(), "statusCheck") ||
			strings.Contains(c.FullPath(), "patient/guest_login") {
			return
		}

		token, err := request.ParseFromRequest(c.Request, MyAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
			return []byte(common.Secrets.JWTSecret), nil
		})
		if err != nil {
			if auto401 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			}
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			accountType := claims["role"].(string)
			if accountType != "clinic" && strings.Contains(c.FullPath(), "membership") {
				if auto401 {
					c.AbortWithStatusJSON(http.StatusUnauthorized, "account should be clinic, getting "+accountType)
				}
				return
			}
			UpdateContextUserModel(c, claims)
		}
		c.Next()
	}
}
