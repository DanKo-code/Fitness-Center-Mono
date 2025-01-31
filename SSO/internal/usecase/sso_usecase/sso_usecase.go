package sso_usecase

import (
	"SSO/internal/constants"
	"SSO/internal/dtos"
	"SSO/internal/models"
	"SSO/internal/repository"
	"SSO/internal/ssoErrors"
	"SSO/pkg/logger"
	"context"
	"errors"
	userGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"strconv"
	"time"
)

var (
	userRole = "client"
)

type SSOUseCase struct {
	ssoRepo    repository.SSORepository
	userClient *grpc.ClientConn
}

func NewSSOUseCase(ssoRepo repository.SSORepository, userClient *grpc.ClientConn) *SSOUseCase {
	return &SSOUseCase{
		ssoRepo:    ssoRepo,
		userClient: userClient,
	}
}

type payload struct {
	userId uuid.UUID
	email  string
	role   string
}

func (ssoUC *SSOUseCase) SignUp(ctx context.Context, cmd *dtos.SignUpRequestCommand) (*dtos.SignUpResponseCommand, error) {

	nuc := userGRPC.NewUserClient(ssoUC.userClient)

	_, err := nuc.GetUserByEmail(ctx, &userGRPC.GetUserByEmailRequest{Email: cmd.Email})
	if err == nil {
		return nil, ssoErrors.UserAlreadyExists
	} else {

		st, ok := status.FromError(err)

		if !ok {
			return nil, err
		}

		switch st.Code() {
		case codes.NotFound:
		default:
			return nil, err
		}
	}

	stream, err := nuc.CreateUser(context.TODO())
	if err != nil {
		logger.ErrorLogger.Printf("Error getted stream for creating user")
		return nil, err
	}

	userDataForCreate := &userGRPC.UserDataForCreate{
		Email:    cmd.Email,
		Role:     userRole,
		Password: cmd.Password,
		Name:     cmd.Name,
	}

	createUserRequest := &userGRPC.CreateUserRequest{
		Payload: &userGRPC.CreateUserRequest_UserDataForCreate{
			UserDataForCreate: userDataForCreate,
		},
	}

	err = stream.Send(createUserRequest)
	if err != nil {
		return nil, err
	}

	createdUser, err := stream.CloseAndRecv()
	if err != nil {
		logger.ErrorLogger.Printf("Failed get createdUser")
		return nil, err
	}

	if createdUser.UserObject == nil {
		return nil, ssoErrors.VoidUserData
	}

	payload := payload{uuid.MustParse(createdUser.UserObject.Id), createdUser.UserObject.Email, createdUser.UserObject.Role}

	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	accessToken, err := GenerateAccessToken(payload, jwtSecret)
	if err != nil {
		logger.ErrorLogger.Printf("Error GenerateAccessToken: %v", err)
		return nil, err
	}

	refreshToken, err := GenerateRefreshToken(uuid.MustParse(createdUser.UserObject.Id), jwtSecret)
	if err != nil {
		logger.ErrorLogger.Printf("Error GenerateRefreshToken: %v", err)
		return nil, err
	}

	refreshSession := &models.RefreshSessions{
		Id:           uuid.New(),
		UserId:       uuid.MustParse(createdUser.UserObject.Id),
		RefreshToken: refreshToken,
		FingerPrint:  cmd.FingerPrint,
		CreatedTime:  time.Now(),
		UpdatedTime:  time.Now(),
	}

	_, err = ssoUC.ssoRepo.CreateRefreshSession(ctx, refreshSession)
	if err != nil {
		logger.ErrorLogger.Printf("Error CreateRefreshSession: %v", err)
		return nil, err
	}

	userObject := &userGRPC.UserObject{
		Id:          createdUser.UserObject.Id,
		Name:        createdUser.UserObject.Name,
		Email:       createdUser.UserObject.Email,
		Role:        createdUser.UserObject.Role,
		Photo:       createdUser.UserObject.Photo,
		CreatedTime: createdUser.UserObject.CreatedTime,
		UpdatedTime: createdUser.UserObject.UpdatedTime,
	}

	signUpResponseCommand := &dtos.SignUpResponseCommand{
		AccessToken:            accessToken,
		RefreshToken:           refreshToken,
		AccessTokenExpiration:  strconv.FormatInt(int64(constants.AccessTokenExpiration), 10),
		RefreshTokenExpiration: strconv.FormatInt(int64(constants.RefreshTokenExpiration), 10),
		User:                   userObject,
	}

	return signUpResponseCommand, nil
}

func (ssoUC *SSOUseCase) SignIn(ctx context.Context, siReq *dtos.SignInRequestCommand) (*dtos.SignInResponseCommand, error) {

	nuc := userGRPC.NewUserClient(ssoUC.userClient)

	user, err := nuc.GetUserByEmail(ctx, &userGRPC.GetUserByEmailRequest{Email: siReq.Email})
	if err != nil {

		st, ok := status.FromError(err)

		if ok == false {
			return nil, err
		}

		switch st.Code() {
		case codes.NotFound:
		default:
			return nil, err
		}

		return nil, err
	}

	checkPasswordRequest := &userGRPC.CheckPasswordRequest{
		UserId:   user.UserObject.Id,
		Password: siReq.Password,
	}

	_, err = nuc.CheckPassword(ctx, checkPasswordRequest)
	if err != nil {
		logger.ErrorLogger.Printf("Error CheckPassword: %v", err)
		return nil, ssoErrors.InvalidPassword
	}

	err = ssoUC.ssoRepo.DeleteRefreshSessionByUserId(ctx, uuid.MustParse(user.UserObject.Id))
	if err != nil && !errors.Is(err, ssoErrors.RefreshSessionNotFound) {
		logger.ErrorLogger.Printf("Error DeleteRefreshSessionByUserId: %v", err)
		return nil, err
	}

	payload := payload{uuid.MustParse(user.UserObject.Id), user.UserObject.Email, user.UserObject.Role}
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	accessToken, err := GenerateAccessToken(payload, jwtSecret)
	if err != nil {
		logger.ErrorLogger.Printf("Error GenerateAccessToken: %v", err)
		return nil, err
	}

	refreshToken, err := GenerateRefreshToken(uuid.MustParse(user.UserObject.Id), jwtSecret)
	if err != nil {
		logger.ErrorLogger.Printf("Error GenerateRefreshToken: %v", err)
		return nil, err
	}

	refreshSession := &models.RefreshSessions{}
	refreshSession.Id = uuid.New()
	refreshSession.UserId = uuid.MustParse(user.UserObject.Id)
	refreshSession.RefreshToken = refreshToken
	refreshSession.FingerPrint = siReq.FingerPrint
	refreshSession.CreatedTime = time.Now()
	refreshSession.UpdatedTime = time.Now()

	_, err = ssoUC.ssoRepo.CreateRefreshSession(ctx, refreshSession)
	if err != nil {
		logger.ErrorLogger.Printf("Error CreateRefreshSession: %v", err)
		return nil, err
	}

	siRes := &dtos.SignInResponseCommand{}
	siRes.RefreshToken = refreshToken
	siRes.AccessToken = accessToken
	siRes.AccessTokenExpiration = strconv.FormatInt(int64(constants.AccessTokenExpiration), 10)
	siRes.RefreshTokenExpiration = strconv.FormatInt(int64(constants.RefreshTokenExpiration), 10)
	siRes.User = &userGRPC.UserObject{
		Id:          user.UserObject.Id,
		Name:        user.UserObject.Name,
		Email:       user.UserObject.Email,
		Role:        user.UserObject.Role,
		Photo:       user.UserObject.Photo,
		CreatedTime: user.UserObject.CreatedTime,
		UpdatedTime: user.UserObject.UpdatedTime,
	}

	return siRes, nil
}

func (ssoUC *SSOUseCase) LogOut(ctx context.Context, refreshToken string) error {

	err := ssoUC.ssoRepo.DeleteRefreshSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		logger.ErrorLogger.Printf("Error LogOut: %v", err)
		return err
	}

	return nil
}

func (ssoUC *SSOUseCase) Refresh(ctx context.Context, cmd *dtos.RefreshRequestCommand) (*dtos.RefreshResponseCommand, error) {

	nuc := userGRPC.NewUserClient(ssoUC.userClient)

	token, err := VerifyRefreshToken(cmd.RefreshToken)
	if err != nil {
		logger.ErrorLogger.Printf("Error VerifyRefreshToken: %v", err)
		return nil, err
	}

	refreshSession, err := ssoUC.ssoRepo.GetRefreshSession(ctx, cmd.RefreshToken)
	if err != nil {
		return nil, err
	}

	if cmd.FingerPrint != refreshSession.FingerPrint {
		logger.ErrorLogger.Printf("Error not valid FingerPrint: %v", cmd.FingerPrint)
		return nil, ssoErrors.InvalidFingerPrint
	}

	if err := ssoUC.ssoRepo.DeleteRefreshSessionByRefreshToken(ctx, cmd.RefreshToken); err != nil {
		logger.ErrorLogger.Printf("Error DeleteRefreshSessionByRefreshToken: %v", cmd.RefreshToken)
		return nil, err
	}

	var userId uuid.UUID
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if id, ok := claims["id"].(string); ok {
			userId, err = uuid.Parse(id)
			if err != nil {
				logger.ErrorLogger.Printf("Error uuid.Parse: %v", err)
				return nil, err
			}
		} else {
			logger.ErrorLogger.Printf("Invalid RefreshToken")
			return nil, ssoErrors.InvalidRefreshToken
		}
	}

	user, err := nuc.GetUserById(ctx, &userGRPC.GetUserByIdRequest{Id: userId.String()})
	if err != nil {

		st, ok := status.FromError(err)

		if !ok {
			return nil, err
		}

		switch st.Code() {
		case codes.NotFound:
			return nil, ssoErrors.UserNotFound
		default:
			return nil, err
		}
	}

	payload := payload{uuid.MustParse(user.UserObject.Id), user.UserObject.Email, user.UserObject.Role}

	accessTokenNew, err := GenerateAccessToken(payload, []byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		logger.ErrorLogger.Printf("Error GenerateAccessToken: %v", err)
		return nil, err
	}

	refreshTokenNew, err := GenerateRefreshToken(uuid.MustParse(user.UserObject.Id), []byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		logger.ErrorLogger.Printf("Error GenerateRefreshToken: %v", err)
		return nil, err
	}

	refreshSessionNew := &models.RefreshSessions{
		Id:           uuid.New(),
		UserId:       uuid.MustParse(user.UserObject.Id),
		RefreshToken: refreshTokenNew,
		FingerPrint:  cmd.FingerPrint,
		CreatedTime:  time.Now(),
		UpdatedTime:  time.Now(),
	}

	_, err = ssoUC.ssoRepo.CreateRefreshSession(ctx, refreshSessionNew)
	if err != nil {
		return nil, err
	}

	refreshResponseCommand := &dtos.RefreshResponseCommand{
		AccessToken:            accessTokenNew,
		RefreshToken:           refreshTokenNew,
		AccessTokenExpiration:  strconv.FormatInt(int64(constants.AccessTokenExpiration), 10),
		RefreshTokenExpiration: strconv.FormatInt(int64(constants.RefreshTokenExpiration), 10),
		User: &userGRPC.UserObject{
			Id:          user.UserObject.Id,
			Name:        user.UserObject.Name,
			Email:       user.UserObject.Email,
			Role:        user.UserObject.Role,
			Photo:       user.UserObject.Photo,
			CreatedTime: user.UserObject.CreatedTime,
			UpdatedTime: user.UserObject.UpdatedTime,
		},
	}

	return refreshResponseCommand, nil
}

func GenerateAccessToken(payload payload, jwtSecret []byte) (string, error) {
	claims := jwt.MapClaims{
		"user_id": payload.userId,
		"role":    payload.role,
		"exp":     time.Now().Add(constants.AccessTokenExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken(clientId uuid.UUID, jwtSecret []byte) (string, error) {
	claims := jwt.MapClaims{
		"id":  clientId,
		"exp": time.Now().Add(constants.RefreshTokenExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func VerifyRefreshToken(refreshToken string) (*jwt.Token, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
