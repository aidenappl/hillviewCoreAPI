package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/jwt"
	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/structs"
)

type contextKey struct {
	name string
}

type userContextKey string

var (
	UserContextType userContextKey = "user"
)

var JWTClaimsCtxKey = &contextKey{"jwt_claims"}

var UserModelCtxKey = &contextKey{"user_model"}

var StudentAllowedRoutes = []string{
	"/core/v1.1/admin/video",
	"/core/v1.1/admin/playlist",
	"/core/v1.1/admin/upload",
}

func WithClaimsValue(ctx context.Context) *jwt.HVJwtClaims {
	val, ok := ctx.Value(JWTClaimsCtxKey).(*jwt.HVJwtClaims)
	if !ok {
		return nil
	}

	return val
}

func WithUserModelValue(ctx context.Context) *structs.User {
	val, ok := ctx.Value(UserContextType).(*structs.User)
	if !ok {
		return nil
	}

	return val
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func MuxHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, "+
			"Content-Type, "+
			"Cookie, "+
			"Accept-Encoding, "+
			"Connection, "+
			"Content-Length")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Server", "Go")
		next.ServeHTTP(w, r)
	})
}

func TokenHandlers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get token from header
		rawToken := r.Header.Get("Authorization")
		splitToken := strings.Split(rawToken, "Bearer ")

		next.ServeHTTP(w, r)

		if len(splitToken) != 2 {
			return
		}
		rawToken = splitToken[1]

		if len(rawToken) < 1 {
			return
		}

		// parse token validity
		token, err := jwt.ParseJWT(rawToken)
		if err != nil {
			return
		}

		claims := token.Claims.(*jwt.HVJwtClaims)

		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			return
		}

		if userID != 0 && len(r.RequestURI) != 0 && len(r.Method) != 0 {
			err = query.InsertRequestLog(db.DB, userID, r.RequestURI, r.Method)
			if err != nil {
				log.Println(fmt.Errorf("failed to insert request log: %w", err))
			}
		}

	})
}

func AccessTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get token from header
		rawToken := r.Header.Get("Authorization")
		splitToken := strings.Split(rawToken, "Bearer ")

		if len(splitToken) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		rawToken = splitToken[1]

		if len(rawToken) < 1 {
			http.Error(w, "Missing Authorization token", http.StatusUnauthorized)
			return
		}

		// parse token validity
		token, err := jwt.ParseJWT(rawToken)
		if err != nil {
			if strings.Contains(err.Error(), "token is expired") {
				http.Error(w, "Token is expired", http.StatusUnauthorized)
			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
			return
		}

		claims := token.Claims.(*jwt.HVJwtClaims)

		claimsValid, resp, err := jwt.ValidJWT(r.Context(), rawToken, claims, &jwt.HVJwtClaims{Type: jwt.AccessToken})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if !claimsValid {
			if resp.Expired {
				http.Error(w, "Token is expired", http.StatusUnauthorized)
			}
			if resp.Revoked {
				http.Error(w, "Token is revoked", http.StatusUnauthorized)
			}
			if resp.InvalidIssuer || resp.Err || resp.Invalid {
				http.Error(w, "Invalid token, bad issuer, response, or invalid", http.StatusUnauthorized)
			}
			return
		}

		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			log.Println("failed to convert user id to int", err.Error())
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		user, err := query.FindUser(db.DB, query.FindUserRequest{ID: &userID})
		if err != nil {
			log.Println("failed to find user by id", err.Error())
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserContextType, user)
		r = r.WithContext(ctx)

		// check user permissions
		if user.Authentication.ID == 1 || user.Authentication.ID == 9 {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if user.Authentication.ShortName == "student" {
			// check route
			for _, route := range StudentAllowedRoutes {
				if strings.Contains(r.RequestURI, route) {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "you do not have permission to access this resource", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
