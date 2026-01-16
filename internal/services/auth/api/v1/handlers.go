package v1

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

var internalServerError = echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")

func (h *AuthHandler) Register(c echo.Context) error {
	type RegisterRequest struct {
		Name     string `json:"name" validate:"required"`
		Age      int16  `json:"age" validate:"required,gte=0,lte=130"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

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

func (h *AuthHandler) Login(c echo.Context) error {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

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
