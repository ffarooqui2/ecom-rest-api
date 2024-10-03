package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ffarooqui2/ecom/config"
	"github.com/ffarooqui2/ecom/types"
	"github.com/ffarooqui2/ecom/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "userID"

// WithJWTAuth is a middleware that checks if the request has a valid JWT token and adds the user to the context if it is valid
func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc { 
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the token from the request
		tokenString := utils.GetTokenFromRequest(r)

		// If the token is empty, return an error
		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		// If the token is invalid, return an error
		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		// Get the userID from the token
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		// Convert the userID to an int
		userID, err := strconv.Atoi(str)
		if err != nil {
			log.Printf("failed to convert userID to int: %v", err)
			permissionDenied(w)
			return
		}

		// Get the user from the store
		u, err := store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			permissionDenied(w)
			return
		}

		// Add the user to the context - this will allow us to access the user in the handler
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)

		// Call the function if the token is valid
		handlerFunc(w, r)
	}
}

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(int(userID)),
		"expiresAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

// validateJWT validates a JWT token
func validateJWT(tokenString string) (*jwt.Token, error) {
	// Parse the token
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key
		return []byte(config.Envs.JWTSecret), nil
	})
}

// permissionDenied is a helper function that writes a permission denied error to the response 
func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}

	return userID
}