package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pascaldekloe/jwt"
	"net/http"
	"os"
	"scoop-order/cmd/src"
	"strings"
	"time"
)

func(server *Server) authenticate() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		if len(ctx.Request.Header["Authorization"]) == 0{
			ctx.Set("permission", "denied")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("invalid or missing authentication token <Authorization>")))
			return
		}
		authorizationHeader := ctx.Request.Header["Authorization"][0]

		if authorizationHeader == "" {
			ctx.Set("permission", "denied")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("invalid or missing authentication token <Authorization:empty>")))
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")

		if len(headerParts) < 2 {
			ctx.Set("permission", "denied")
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		if headerParts[0] != "JWT" {
			ctx.Set("permission", "denied")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("unsupported authorization type %s", headerParts[0])))
			return
		}
		token := headerParts[1]
		claims, err := jwt.HMACCheck([]byte(token), []byte(os.Getenv("JWT_SECRET")))

		if err != nil {
			ctx.Set("permission", "denied")
			fmt.Println(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("invalid or missing authentication token, %s", err.Error())))
			return
		}

		// Check token expiration
		if !claims.Valid(time.Now()) {
			ctx.Set("permission", "denied")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("expired token")))
			return
		}

		// Decode token to get claims
		data, err := src.DecodeToken(token)
		if err != nil {
			ctx.Set("permission", "denied")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}


		// Identify user by ID using information provided in claims
		user, err := server.transaction.SelectUser(ctx, int32(data.UserID))
		if err != nil {
			ctx.Set("permission", "denied")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// We need to set the user model to context so that we can use for permission based authentication

		server.InfoLogger.Println("Authorization: "+ authorizationHeader, nil)
		server.Logger.SetPrefix("DEBUG - ")
		server.Logger.Println("Authorization: "+ authorizationHeader)
		ctx.Set("user", user)
		ctx.Set("permission", "allowed")
		ctx.Next()
	}
}

func (server *Server) requireAuthentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		permission, _ := ctx.Get("permission")
		if permission != "allowed"{
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("you must be authenticated to access this resource")))
			return
		}
		ctx.Next()
	}
}
