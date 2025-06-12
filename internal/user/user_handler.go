package user

import (
	"net/http"

	"boilerplate/internal/user/model"
	"boilerplate/pkg/response"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *UserService
}

func NewUserHandler(userService *UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(c echo.Context) error {
	var input model.RegisterInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "invalid request payload", err)
	}

	if err := c.Validate(&input); err != nil {
		return response.ValidationError(c, err)
	}

	user, err := h.userService.Register(input)
	if err != nil {
		return response.BadRequest(c, "registration failed", err)
	}

	return response.Success(c, http.StatusCreated, "User registered successfully", user)
}

func (h *UserHandler) Login(c echo.Context) error {
	var input model.LoginInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "invalid request payload", err)
	}

	if err := c.Validate(&input); err != nil {
		return response.ValidationError(c, err)
	}

	token, err := h.userService.Login(input)
	if err != nil {
		return response.Unauthorized(c, "invalid credentials", err)
	}

	return response.Success(c, http.StatusOK, "Login successful", map[string]string{
		"token": token,
	})
}

func (h *UserHandler) Logout(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	if err := h.userService.Logout(userID); err != nil {
		return response.InternalServerError(c, "logout failed", err)
	}

	return response.Success(c, http.StatusOK, "Logout successful", nil)
}

func (h *UserHandler) GetMe(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		return response.NotFound(c, "user not found", err)
	}

	return response.Success(c, http.StatusOK, "User profile retrieved successfully", user)
}

func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	var input model.UpdateProfileInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "invalid request payload", err)
	}

	if err := c.Validate(&input); err != nil {
		return response.ValidationError(c, err)
	}

	user, err := h.userService.UpdateProfile(userID, input)
	if err != nil {
		return response.BadRequest(c, "failed to update profile", err)
	}

	return response.Success(c, http.StatusOK, "Profile updated successfully", user)
}

func (h *UserHandler) DeleteAccount(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	if err := h.userService.DeleteAccount(userID); err != nil {
		return response.InternalServerError(c, "failed to delete account", err)
	}

	return response.Success(c, http.StatusOK, "Account deleted successfully", nil)
}

func (h *UserHandler) GetAllUsers(c echo.Context) error {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		return response.InternalServerError(c, "failed to get all users", err)
	}

	return response.Success(c, http.StatusOK, "Users retrieved successfully", users)
}
