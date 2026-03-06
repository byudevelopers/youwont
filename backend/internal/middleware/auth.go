package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"youwont.api/internal/model"
)

// UserFinder is the interface the middleware needs to look up users.
// Defined here (consumer-defined) to avoid coupling to the full service.
type UserFinder interface {
	FindBySupabaseID(ctx context.Context, supabaseID string) (*model.User, error)
}

type Auth struct {
	jwtSecret  []byte
	userFinder UserFinder
}

func NewAuth(jwtSecret string, userFinder UserFinder) *Auth {
	return &Auth{
		jwtSecret:  []byte(jwtSecret),
		userFinder: userFinder,
	}
}

// Required is the Echo middleware that enforces authentication on protected routes.
func (a *Auth) Required(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		tokenStr := extractBearer(c)
		if tokenStr == "" {
			return c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "missing authorization token"))
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return a.jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "invalid token"))
		}

		sub, err := token.Claims.GetSubject()
		if err != nil || sub == "" {
			return c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "invalid token claims"))
		}

		user, err := a.userFinder.FindBySupabaseID(c.Request().Context(), sub)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errResp("INTERNAL", "server error"))
		}
		if user == nil {
			return c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "user not found"))
		}

		c.Set("user", user)
		return next(c)
	}
}

// ExtractSubFromToken parses the JWT and returns the sub claim.
// Used by POST /users which needs the supabase_id but no existing MongoDB user.
func (a *Auth) ExtractSubFromToken(c *echo.Context) (string, error) {
	tokenStr := extractBearer(c)
	if tokenStr == "" {
		return "", jwt.ErrTokenMalformed
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return a.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	return token.Claims.GetSubject()
}

// UserFromContext retrieves the authenticated user from the Echo context.
func UserFromContext(c *echo.Context) *model.User {
	u, _ := c.Get("user").(*model.User)
	return u
}

func extractBearer(c *echo.Context) string {
	auth := c.Request().Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return auth[7:]
	}
	return ""
}

func errResp(code, message string) map[string]interface{} {
	return map[string]interface{}{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	}
}
