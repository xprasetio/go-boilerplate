package middleware

import (
	"boilerplate/internal/user/model"
	"boilerplate/pkg/response"

	"github.com/labstack/echo/v4"
)

func AdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userObj := c.Get("user")
			if userObj == nil {
				return response.Unauthorized(c, "unauthorized", nil)
			}

			user, ok := userObj.(*model.User)
			if !ok {
				return response.Unauthorized(c, "invalid user data", nil)
			}

			if user == nil {
				return response.Unauthorized(c, "user not found", nil)
			}

			if !user.IsAdmin() {
				return response.Forbidden(c, "access denied: admin role required", nil)
			}

			return next(c)
		}
	}
}
