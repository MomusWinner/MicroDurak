package http

import (
	"errors"
	"net/http"

	"github.com/MommusWinner/MicroDurak/internal/services/auth/domain/cases"
	"github.com/MommusWinner/MicroDurak/internal/services/auth/domain/props"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	useCase *cases.AuthUseCase
}

func NewAuthHandler(useCase *cases.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		useCase: useCase,
	}
}

type AuthResponse struct {
	PlayerID string `json:"player_id"`
	Token    string `json:"token"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Age      int16  `json:"age" validate:"required,gte=0,lte=130"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

var internalServerError = echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")

// Register creates a new user account
// @Summary Register a new user
// @Description Creates a new user account with the provided information
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} AuthResponse "User successfully registered"
// @Failure 400
// @Failure 409
// @Failure 500
// @Router /register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	r := new(RegisterRequest)
	if err := c.Bind(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := c.Validate(r); err != nil {
		return err
	}

	resp, err := h.useCase.Register(props.RegisterReq{Name: r.Name, Age: int(r.Age), Email: r.Email, Password: r.Password})

	if err != nil {
		if errors.Is(err, cases.ErrEmailAlreadyTaken) {
			return c.String(http.StatusConflict, "Email Taken")
		}
		return internalServerError
	}

	return c.JSON(http.StatusCreated, AuthResponse{
		PlayerID: resp.PlayerId.String(),
		Token:    resp.Token,
	})
}

// Login authenticates a user and returns a token
// @Summary Login user
// @Description Authenticates a user with email and password, returns JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200
// @Failure 400
// @Failure 403
// @Failure 500
// @Router /login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	r := new(LoginRequest)
	if err := c.Bind(r); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	resp, err := h.useCase.Login(props.LoginReq{Email: r.Email, Password: r.Password})

	if err != nil {
		if errors.Is(err, cases.ErrLoginFailed) {
			return c.String(http.StatusForbidden, "Login failed")
		} else {
			return internalServerError
		}
	}

	return c.JSON(http.StatusOK, AuthResponse{
		PlayerID: resp.PlayerId,
		Token:    resp.Token,
	})
}
