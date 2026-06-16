package userprofile

import (
	"errors"
	"net/http"

	errInt "github.com/bobfive1/user-management-api/internal/error"
	"github.com/bobfive1/user-management-api/internal/logger"
	"github.com/bobfive1/user-management-api/internal/validation"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

const (
	LoginFailMessage    = "Invalid login or password"
	UserProfileNotExist = "The requested user profile does not exist."
)

type UserProfileHandler struct {
	service UserProfileService
	log     *zap.SugaredLogger
}

func NewUserProfileHandler(service UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{service: service, log: logger.GetLogger("UserProfileHandler")}
}

func (h *UserProfileHandler) RegisterRoutes(router gin.IRouter) {
	profiles := router.Group("/userprofiles")
	profiles.POST("", h.create)
	profiles.GET("", h.list)
	profiles.GET("/:user_id", h.getByID)
	profiles.PUT("/:user_id", h.update)
	profiles.DELETE("/:user_id", h.delete)
	profiles.POST("/login", h.login)
}

func (h *UserProfileHandler) create(c *gin.Context) {
	var request InsertUserProfileRequest

	request, err := validation.ShouldBindJSONWithValidate(c, request)
	if err != nil {
		c.Error(err)
		return
	}

	profile, err := h.service.Create(c, request)
	if err != nil {
		c.Error(errInt.NewInternalServerError(err.Error()))
		return
	}

	response := NewUserProfileResponse("200", "Success", profile)
	c.JSON(http.StatusOK, response)
}

func (h *UserProfileHandler) list(c *gin.Context) {
	profiles, err := h.service.List(c)
	if err != nil {
		c.Error(errInt.NewInternalServerError(err.Error()))
		return
	}

	response := NewUserProfileResponse("200", "Success", profiles)
	c.JSON(http.StatusOK, response)
}

func (h *UserProfileHandler) getByID(c *gin.Context) {
	userID := c.Param("user_id")

	profile, err := h.service.GetByID(c, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.Error(errInt.NewNotFoundError(UserProfileNotExist))
			return
		}
		c.Error(errInt.NewInternalServerError(err.Error()))
		return
	}

	response := NewUserProfileResponse("200", "Success", profile)
	c.JSON(http.StatusOK, response)
}

func (h *UserProfileHandler) update(c *gin.Context) {
	userID := c.Param("user_id")

	var request UpdateUserProfileRequest
	request, err := validation.ShouldBindJSONWithValidate(c, request)
	if err != nil {
		c.Error(err)
		return
	}

	profile, err := h.service.Update(c, userID, request)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.Error(errInt.NewNotFoundError(UserProfileNotExist))
			return
		}
		c.Error(errInt.NewInternalServerError(err.Error()))
		return
	}

	response := NewUserProfileResponse("200", "Success", profile)
	c.JSON(http.StatusOK, response)
}

func (h *UserProfileHandler) delete(c *gin.Context) {
	userID := c.Param("user_id")

	if err := h.service.Delete(c, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.Error(errInt.NewNotFoundError(UserProfileNotExist))
			return
		}
		c.Error(errInt.NewInternalServerError(err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserProfileHandler) login(c *gin.Context) {
	var request UserProfileLoginRequest

	request, err := validation.ShouldBindJSONWithValidate(c, request)
	if err != nil {
		c.Error(err)
		return
	}

	profile, err := h.service.Login(c, request)
	if err != nil {

		if errors.Is(err, ErrorPassNotMatch) || errors.Is(err, pgx.ErrNoRows) {
			c.Error(errInt.NewBadRequestError("400", LoginFailMessage, nil))
			return
		}
		c.Error(errInt.NewInternalServerError(err.Error()))
		return
	}

	response := NewUserProfileResponse("200", "Success", profile)
	c.JSON(http.StatusOK, response)
}
