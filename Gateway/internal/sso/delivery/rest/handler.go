package rest

import (
	"Gateway/internal/sso/dtos"
	"Gateway/internal/sso/sso_errors"
	logger "Gateway/pkg/logger"
	"context"
	"fmt"
	ssoGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.sso"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	accessTokenMaxAgeDivider  = 100000
	refreshTokenMaxAgeDivider = 1000000000
)

type Handler struct {
	ssoClient *grpc.ClientConn
	validator *validator.Validate
}

func NewHandler(ssoClient *grpc.ClientConn, validator *validator.Validate) *Handler {
	return &Handler{
		ssoClient: ssoClient,
		validator: validator,
	}
}

// SignUp
// @Summary Sign-up a new user
// @Description Sign-up a new user with provided details including name, email, and password. The API also manages fingerprint for enhanced tracking.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param signUpRequest body dtos.SignUpRequestWithOutFingerPrint true "Details for the user sign-up"
// @Success 200 {object} dtos.SignUpResponse "Successful registration response containing access token, expiration, and user details"
// @Failure 400 {object} map[string]string "Invalid request or validation error"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /sso/signUp [post]
func (h *Handler) SignUp(c *gin.Context) {
	suReq := &dtos.SignUpRequest{}

	if err := c.ShouldBindJSON(&suReq); err != nil {
		logger.ErrorLogger.Printf("Error binding SignUpRequest: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.validator.Struct(suReq)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			customMessages := make(map[string]string)
			for _, fieldErr := range validationErrors {
				customMessages[fieldErr.Field()] = getCustomErrorMessage(fieldErr)
			}

			c.JSON(http.StatusBadRequest, gin.H{"errors": customMessages})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal validation error"})
		return
	}

	fingerPrintValue, exists := c.Get(os.Getenv("APP_FINGERPRINT_REQUEST_KEY"))
	if !exists {
		logger.ErrorLogger.Printf("Error getting fingerprint: %v", sso_errors.FingerPrintNotFoundInContext)
		c.JSON(http.StatusInternalServerError, gin.H{"error": sso_errors.FingerPrintNotFoundInContext})
		return
	}

	FingerPrintValueCasted, ok := fingerPrintValue.(string)
	if !ok {
		logger.ErrorLogger.Printf("Error casting fingerprint to string: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": sso_errors.FingerPrintNotFoundInContext})
		return
	}

	suReq.FingerPrint = FingerPrintValueCasted

	ssoClient := ssoGRPC.NewSSOClient(h.ssoClient)

	signUpRequest := &ssoGRPC.SignUpRequest{}
	signUpRequest.Name = suReq.Name
	signUpRequest.Email = strings.ToLower(suReq.Email)
	signUpRequest.Password = suReq.Password
	signUpRequest.FingerPrint = suReq.FingerPrint

	upRes, err := ssoClient.SignUp(context.Background(), signUpRequest)
	if err != nil {
		logger.ErrorLogger.Printf("Error SignUp: %v", err)

		//TODO add statuses?
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ateInt, err := strconv.Atoi(upRes.AccessTokenExpiration)
	if err != nil {
		logger.ErrorLogger.Printf("Error convert AccessTokenExpiration to int: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rteInt, err := strconv.Atoi(upRes.RefreshTokenExpiration)
	if err != nil {
		logger.ErrorLogger.Printf("Error convert RefreshTokenExpiration to int: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie(
		"refreshToken",
		upRes.RefreshToken,
		rteInt/refreshTokenMaxAgeDivider,
		"",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"accessToken":           upRes.GetAccessToken(),
		"accessTokenExpiration": ateInt / accessTokenMaxAgeDivider,
		"user":                  upRes.GetUser(),
	})
}

// SignIn
// @Summary Sign-in  user
// @Description Sign-in for existing user with provided details including email, and password. The API also manages fingerprint for enhanced tracking.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param signInRequest body dtos.SignInRequestWithoutFingerprint true "Details for the user sign-in"
// @Success 200 {object} dtos.SignInResponse "Successful sign-in response containing access token, expiration, and user details"
// @Failure 400 {object} map[string]string "Invalid request or validation error"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /sso/signIn [post]
func (h *Handler) SignIn(c *gin.Context) {
	siReq := &dtos.SignInRequest{}

	if err := c.ShouldBindJSON(siReq); err != nil {
		logger.ErrorLogger.Printf("Error parsing SignInRequest: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.validator.Struct(siReq)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			customMessages := make(map[string]string)
			for _, fieldErr := range validationErrors {
				customMessages[fieldErr.Field()] = getCustomErrorMessage(fieldErr)
			}

			c.JSON(http.StatusBadRequest, gin.H{"errors": customMessages})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal validation error"})
		return
	}

	fingerPrintValue, exists := c.Get(os.Getenv("APP_FINGERPRINT_REQUEST_KEY"))
	if !exists {
		logger.ErrorLogger.Printf("Error binding SignUpRequest: %v", sso_errors.FingerPrintNotFoundInContext)
		c.JSON(http.StatusInternalServerError, gin.H{"error": sso_errors.FingerPrintNotFoundInContext})
		return
	}

	FingerPrintValueCasted, ok := fingerPrintValue.(string)
	if !ok {
		logger.ErrorLogger.Printf("Error casting fingerprint to string: %v", sso_errors.FingerPrintNotFoundInContext)
		c.JSON(http.StatusInternalServerError, gin.H{"error": sso_errors.FingerPrintNotFoundInContext})
		return
	}

	siReq.FingerPrint = FingerPrintValueCasted

	ssoClient := ssoGRPC.NewSSOClient(h.ssoClient)

	signIpRequest := &ssoGRPC.SignInRequest{
		Email:       strings.ToLower(siReq.Email),
		Password:    siReq.Password,
		FingerPrint: siReq.FingerPrint,
	}

	siRes, err := ssoClient.SignIn(context.Background(), signIpRequest)
	if err != nil {
		logger.ErrorLogger.Printf("Error SignIp: %v", err)

		//TODO add statuses?
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ateInt, err := strconv.Atoi(siRes.GetAccessTokenExpiration())
	if err != nil {
		logger.ErrorLogger.Printf("Error convert AccessTokenExpiration to int: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rteInt, err := strconv.Atoi(siRes.GetRefreshTokenExpiration())
	if err != nil {
		logger.ErrorLogger.Printf("Error convert RefreshTokenExpiration to int: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie(
		"refreshToken",
		siRes.GetRefreshToken(),
		rteInt/refreshTokenMaxAgeDivider,
		"",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"user":                  siRes.GetUser(),
		"accessToken":           siRes.GetAccessToken(),
		"accessTokenExpiration": ateInt / accessTokenMaxAgeDivider,
	})
}

// LogOut
// @Summary Log-out  user
// @Description Log-out for entered user. The API also manages fingerprint for enhanced tracking.
// @Tags Authentication
// @Accept json
// @Success 200 "Successful log-out"
// @Failure 400 {object} map[string]string "Invalid request or validation error"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /sso/logOut [post]
func (h *Handler) LogOut(c *gin.Context) {

	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		logger.ErrorLogger.Printf("Error getting refreshToken from cookie: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ssoClient := ssoGRPC.NewSSOClient(h.ssoClient)

	logOutRequest := &ssoGRPC.LogOutRequest{}
	logOutRequest.RefreshToken = refreshToken

	_, err = ssoClient.LogOut(context.Background(), logOutRequest)
	if err != nil {
		logger.ErrorLogger.Printf("Error LogOut: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("refreshToken", "", -1, "/", "", false, true)

	c.Status(http.StatusOK)
}

// Refresh
// @Summary Refreshing tokens for accessing secured resources
// @Description Refreshing tokens for accessing secured resources. Getting Refresh token from cookies. The API also manages fingerprint for enhanced tracking.
// @Tags Authentication
// @Accept json
// @Success 200 {object} dtos.RefreshResponse "Successful Refresh tokens response containing access token, expiration, and user details"
// @Failure 400 {object} map[string]string "Invalid request or validation error"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /sso/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	fingerPrintValue, exists := c.Get(os.Getenv("APP_FINGERPRINT_REQUEST_KEY"))
	if !exists {
		logger.ErrorLogger.Printf("Error getting fingerprint	: %v", sso_errors.FingerPrintNotFoundInContext)
		c.JSON(http.StatusInternalServerError, gin.H{"error": sso_errors.FingerPrintNotFoundInContext})
		return
	}

	FingerPrintValueCasted, ok := fingerPrintValue.(string)
	if !ok {
		logger.ErrorLogger.Printf("Error casting fingerprint to string: %v", sso_errors.FingerPrintNotFoundInContext)
		c.JSON(http.StatusInternalServerError, gin.H{"error": sso_errors.FingerPrintNotFoundInContext})
		return
	}

	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		logger.ErrorLogger.Printf("Error getting refreshToken from cookie: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshRequest := &ssoGRPC.RefreshRequest{
		FingerPrint:  FingerPrintValueCasted,
		RefreshToken: refreshToken,
	}

	ssoClient := ssoGRPC.NewSSOClient(h.ssoClient)

	rRes, err := ssoClient.Refresh(context.Background(), refreshRequest)
	if err != nil {
		logger.ErrorLogger.Printf("Error Refresh: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ateInt, err := strconv.Atoi(rRes.AccessTokenExpiration)
	if err != nil {
		logger.ErrorLogger.Printf("Error convert AccessTokenExpiration to int: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rteInt, err := strconv.Atoi(rRes.GetRefreshTokenExpiration())
	if err != nil {
		logger.ErrorLogger.Printf("Error convert RefreshTokenExpiration to int: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie(
		"refreshToken",
		rRes.GetRefreshToken(),
		rteInt/refreshTokenMaxAgeDivider,
		"",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"user":                  rRes.GetUser(),
		"accessToken":           rRes.GetAccessToken(),
		"accessTokenExpiration": ateInt / accessTokenMaxAgeDivider,
	})
}

func getCustomErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("Поле '%s' обязательно к заполнению.", fe.Field())
	case "email":
		return fmt.Sprintf("Поле '%s' должен быть указан правильный тип электронной почты.", fe.Field())
	case "min":
		return fmt.Sprintf("Поле '%s' должно содержать не менее %s символов.", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("Поле '%s' должно содержать не более %s символов.", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("Поле '%s' не валидно: %s.", fe.Field(), fe.Tag())
	}
}
