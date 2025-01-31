package grpc

import (
	"context"
	"errors"
	"github.com/DanKo-code/Fitness-Center-Abonement/internal/dtos"
	customErrors "github.com/DanKo-code/Fitness-Center-Abonement/internal/errors"
	"github.com/DanKo-code/Fitness-Center-Abonement/internal/usecase"
	"github.com/DanKo-code/Fitness-Center-Abonement/pkg/logger"
	abonementProtobuf "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.abonement"
	serviceGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.service"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"reflect"
	"strings"
	"time"
)

var _ abonementProtobuf.AbonementServer = (*AbonementgRPC)(nil)

type AbonementgRPC struct {
	abonementProtobuf.UnimplementedAbonementServer

	abonementUseCase usecase.AbonementUseCase
	cloudUseCase     usecase.CloudUseCase
	serviceClient    *serviceGRPC.ServiceClient
}

func RegisterAbonementServer(
	gRPC *grpc.Server,
	abonementUseCase usecase.AbonementUseCase,
	cloudUseCase usecase.CloudUseCase,
	serviceClient *serviceGRPC.ServiceClient,
) {
	abonementProtobuf.RegisterAbonementServer(
		gRPC,
		&AbonementgRPC{
			abonementUseCase: abonementUseCase,
			cloudUseCase:     cloudUseCase,
			serviceClient:    serviceClient,
		})
}

func (c *AbonementgRPC) CreateAbonement(g grpc.ClientStreamingServer[abonementProtobuf.CreateAbonementRequest, abonementProtobuf.CreateAbonementResponse]) error {

	abonementData, abonementPhoto, err := GetObjectData(
		&g,
		func(chunk *abonementProtobuf.CreateAbonementRequest) interface{} {
			return chunk.GetAbonementDataForCreate()
		},
		func(chunk *abonementProtobuf.CreateAbonementRequest) []byte {
			return chunk.GetAbonementPhoto()
		},
	)
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid request data")
	}

	if abonementData == nil {
		logger.ErrorLogger.Printf("abonement data is empty")
		return status.Error(codes.InvalidArgument, "abonement data is empty")
	}

	castedAbonementData, ok := abonementData.(*abonementProtobuf.AbonementDataForCreate)
	if !ok {
		logger.ErrorLogger.Printf("abonement data is not of type AbonementProtobuf.AbonementDataForCreate")
		return status.Error(codes.InvalidArgument, "abonement data is not of type AbonementProtobuf.AbonementDataForCreate")
	}

	cmd := &dtos.CreateAbonementCommand{
		Id:           uuid.New(),
		Title:        castedAbonementData.Title,
		Validity:     castedAbonementData.Validity,
		VisitingTime: castedAbonementData.VisitingTime,
		Price:        int(castedAbonementData.Price),
	}

	var photoURL string
	randomID := uuid.New().String()
	if abonementPhoto != nil {
		url, err := c.cloudUseCase.PutObject(context.TODO(), abonementPhoto, "abonement/"+randomID)
		photoURL = url
		if err != nil {
			logger.ErrorLogger.Printf("Failed to create abonement photo in cloud: %v", err)
			return status.Error(codes.Internal, "Failed to create abonement photo in cloud")
		}
	}

	cmd.Photo = photoURL

	abonement, err := c.abonementUseCase.CreateAbonement(context.TODO(), cmd)
	if err != nil {

		if photoURL == "" {
			err := c.cloudUseCase.DeleteObject(context.TODO(), "abonement/"+cmd.Id.String())
			if err != nil {
				logger.ErrorLogger.Printf("Failed to delete abonement photo from cloud: %v", err)
				return status.Error(codes.Internal, "Failed to delete abonement photo in cloud")
			}
		}

		return status.Error(codes.Internal, "Failed to create abonement")
	}

	createAbonementServicesRequest := &serviceGRPC.CreateAbonementServicesRequest{
		AbonementService: &serviceGRPC.AbonementService{
			AbonementId: abonement.Id.String(),
			ServiceId:   castedAbonementData.ServicesIds,
		},
	}
	services, err := (*c.serviceClient).CreateAbonementServices(context.TODO(), createAbonementServicesRequest)
	if err != nil {
		return err
	}

	var abonementsServices *serviceGRPC.GetAbonementsServicesResponse
	if services != nil {
		getAbonementsServicesRequest := &serviceGRPC.GetAbonementsServicesRequest{
			AbonementIds: []string{abonement.Id.String()},
		}
		abonementsServices, err = (*c.serviceClient).GetAbonementsServices(context.TODO(), getAbonementsServicesRequest)
		if err != nil {
			return err
		}
	}

	abonementObject := &abonementProtobuf.AbonementObject{
		Id:           abonement.Id.String(),
		Title:        abonement.Title,
		Validity:     abonement.Validity,
		VisitingTime: abonement.VisitingTime,
		Photo:        abonement.Photo,
		Price:        int32(abonement.Price),
		CreatedTime:  abonement.CreatedTime.String(),
		UpdatedTime:  abonement.UpdatedTime.String(),
	}

	var abonementWithServices *abonementProtobuf.AbonementWithServices
	if abonementsServices != nil {
		abonementWithServices = &abonementProtobuf.AbonementWithServices{
			Abonement: abonementObject,
			Services:  abonementsServices.AbonementIdsWithServices[0].ServiceObjects,
		}
	} else {
		abonementWithServices = &abonementProtobuf.AbonementWithServices{
			Abonement: abonementObject,
			Services:  nil,
		}
	}

	response := &abonementProtobuf.CreateAbonementResponse{
		AbonementWithServices: abonementWithServices,
	}

	err = g.SendAndClose(response)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to send abonement create response: %v", err)
		return status.Error(codes.Internal, "Failed to send abonement create response")
	}

	return nil
}

func (c *AbonementgRPC) GetAbonementById(ctx context.Context, request *abonementProtobuf.GetAbonementByIdRequest) (*abonementProtobuf.GetAbonementByIdResponse, error) {

	abonement, err := c.abonementUseCase.GetAbonementById(ctx, uuid.MustParse(request.Id))
	if err != nil {

		if errors.Is(err, customErrors.AbonementNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, err
	}

	abonementObject := &abonementProtobuf.AbonementObject{
		Id:            abonement.Id.String(),
		Title:         abonement.Title,
		Validity:      abonement.Validity,
		VisitingTime:  abonement.VisitingTime,
		Photo:         abonement.Photo,
		Price:         int32(abonement.Price),
		CreatedTime:   abonement.CreatedTime.String(),
		UpdatedTime:   abonement.UpdatedTime.String(),
		StripePriceId: abonement.StripePriceId,
	}

	response := &abonementProtobuf.GetAbonementByIdResponse{
		AbonementObject: abonementObject,
	}

	return response, nil
}

func (c *AbonementgRPC) UpdateAbonement(g grpc.ClientStreamingServer[abonementProtobuf.UpdateAbonementRequest, abonementProtobuf.UpdateAbonementResponse]) error {
	abonementData, abonementPhoto, err := GetObjectData(
		&g,
		func(chunk *abonementProtobuf.UpdateAbonementRequest) interface{} {
			return chunk.GetAbonementDataForUpdate()
		},
		func(chunk *abonementProtobuf.UpdateAbonementRequest) []byte {
			return chunk.GetAbonementPhoto()
		},
	)
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid request data")
	}

	if abonementData == nil {
		logger.ErrorLogger.Printf("abonement data is empty")
		return status.Error(codes.InvalidArgument, "abonement data is empty")
	}

	castedAbonementData, ok := abonementData.(*abonementProtobuf.AbonementDataForUpdate)
	if !ok {
		logger.ErrorLogger.Printf("abonement data is not of type AbonementProtobuf.AbonementDataForCreate")
		return status.Error(codes.InvalidArgument, "abonement data is not of type AbonementProtobuf.AbonementDataForCreate")
	}

	cmd := &dtos.UpdateAbonementCommand{
		Id:           uuid.MustParse(castedAbonementData.Id),
		Title:        castedAbonementData.Title,
		Validity:     castedAbonementData.Validity,
		VisitingTime: castedAbonementData.VisitingTime,
		Price:        int(castedAbonementData.Price),
		UpdatedTime:  time.Now(),
	}

	existingAbonement, err := c.abonementUseCase.GetAbonementById(context.TODO(), uuid.MustParse(castedAbonementData.Id))
	if err != nil {
		return status.Error(codes.NotFound, "abonement not found")
	}

	var photoURL string
	randomID := uuid.New().String()
	if abonementPhoto != nil {
		if existingAbonement.Photo != "" {
			prefix := "abonement/"
			index := strings.Index(existingAbonement.Photo, prefix)
			var s3PhotoKey string
			if index != -1 {
				s3PhotoKey = existingAbonement.Photo[index+len(prefix):]
			} else {
				logger.ErrorLogger.Printf("Prefix not found")
			}

			exists, err := c.cloudUseCase.ObjectExists(context.TODO(), "abonement/"+s3PhotoKey)
			if err != nil {
				return status.Error(codes.Internal, "can't find previous photo meta")
			}

			if exists {
				err := c.cloudUseCase.DeleteObject(context.TODO(), "abonement/"+s3PhotoKey)
				if err != nil {
					return err
				}
			}
		}

		url, err := c.cloudUseCase.PutObject(context.TODO(), abonementPhoto, "abonement/"+randomID)
		photoURL = url
		if err != nil {
			logger.ErrorLogger.Printf("Failed to create abonement photo in cloud: %v", err)
			return status.Error(codes.Internal, "Failed to create abonement photo in cloud")
		}
	}

	cmd.Photo = photoURL

	abonement, err := c.abonementUseCase.UpdateAbonement(context.TODO(), cmd)
	if err != nil {
		return status.Error(codes.Internal, "Failed to update abonement")
	}

	updateAbonementServicesRequest := &serviceGRPC.UpdateAbonementServicesRequest{
		AbonementService: &serviceGRPC.AbonementService{
			AbonementId: abonement.Id.String(),
			ServiceId:   castedAbonementData.ServicesIds,
		},
	}

	if len(castedAbonementData.ServicesIds) > 0 {
		_, err := (*c.serviceClient).UpdateAbonementServices(context.TODO(), updateAbonementServicesRequest)
		if err != nil {
			return err
		}
	}

	var abonementsServices *serviceGRPC.GetAbonementsServicesResponse
	getAbonementsServicesRequest := &serviceGRPC.GetAbonementsServicesRequest{
		AbonementIds: []string{abonement.Id.String()},
	}
	abonementsServices, err = (*c.serviceClient).GetAbonementsServices(context.TODO(), getAbonementsServicesRequest)
	if err != nil {
		return err
	}

	abonementObject := &abonementProtobuf.AbonementObject{
		Id:           abonement.Id.String(),
		Title:        abonement.Title,
		Validity:     abonement.Validity,
		VisitingTime: abonement.VisitingTime,
		Photo:        abonement.Photo,
		Price:        int32(abonement.Price),
		CreatedTime:  abonement.CreatedTime.String(),
		UpdatedTime:  abonement.UpdatedTime.String(),
	}

	var abonementWithServices *abonementProtobuf.AbonementWithServices
	if abonementsServices != nil {
		abonementWithServices = &abonementProtobuf.AbonementWithServices{
			Abonement: abonementObject,
			Services:  abonementsServices.AbonementIdsWithServices[0].ServiceObjects,
		}
	} else {
		abonementWithServices = &abonementProtobuf.AbonementWithServices{
			Abonement: abonementObject,
			Services:  nil,
		}
	}

	response := &abonementProtobuf.UpdateAbonementResponse{
		AbonementWithServices: abonementWithServices,
	}

	err = g.SendAndClose(response)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to send abonement update response: %v", err)
		return status.Error(codes.Internal, "Failed to send abonement update response")
	}

	return nil
}

func (c *AbonementgRPC) DeleteAbonementById(ctx context.Context, request *abonementProtobuf.DeleteAbonementByIdRequest) (*abonementProtobuf.DeleteAbonementByIdResponse, error) {
	deletedAbonement, err := c.abonementUseCase.DeleteAbonementById(ctx, uuid.MustParse(request.Id))
	if err != nil {
		return nil, err
	}

	abonementObject := &abonementProtobuf.AbonementObject{
		Id:            deletedAbonement.Id.String(),
		Title:         deletedAbonement.Title,
		Validity:      deletedAbonement.Validity,
		VisitingTime:  deletedAbonement.VisitingTime,
		Photo:         deletedAbonement.Photo,
		Price:         int32(deletedAbonement.Price),
		CreatedTime:   deletedAbonement.CreatedTime.String(),
		UpdatedTime:   deletedAbonement.UpdatedTime.String(),
		StripePriceId: deletedAbonement.StripePriceId,
	}

	deleteAbonementByIdResponse := &abonementProtobuf.DeleteAbonementByIdResponse{
		AbonementObject: abonementObject,
	}

	return deleteAbonementByIdResponse, nil
}

func (c *AbonementgRPC) GetAbonements(ctx context.Context, _ *emptypb.Empty) (*abonementProtobuf.GetAbonementsResponse, error) {

	abonementes, err := c.abonementUseCase.GetAbonementes(ctx)
	if err != nil {
		return nil, err
	}

	var abonementObjects []*abonementProtobuf.AbonementObject

	for _, abonement := range abonementes {

		abonementObject := &abonementProtobuf.AbonementObject{
			Id:            abonement.Id.String(),
			Title:         abonement.Title,
			Validity:      abonement.Validity,
			VisitingTime:  abonement.VisitingTime,
			Photo:         abonement.Photo,
			Price:         int32(abonement.Price),
			CreatedTime:   abonement.CreatedTime.String(),
			UpdatedTime:   abonement.UpdatedTime.String(),
			StripePriceId: abonement.StripePriceId,
		}

		abonementObjects = append(abonementObjects, abonementObject)
	}

	response := &abonementProtobuf.GetAbonementsResponse{AbonementObjects: abonementObjects}

	return response, nil
}

func (c *AbonementgRPC) GetAbonementsWithServices(ctx context.Context, _ *emptypb.Empty) (*abonementProtobuf.GetAbonementsWithServicesResponse, error) {
	abonementsWithServices, err := c.abonementUseCase.GetAbonementsWithServices(ctx)
	if err != nil {
		return nil, err
	}

	var abonementsWithServicesForResponse []*abonementProtobuf.AbonementWithServices
	for _, abonementWithServices := range abonementsWithServices {

		abonementObject := &abonementProtobuf.AbonementObject{
			Id:            abonementWithServices.Abonement.Id.String(),
			Title:         abonementWithServices.Abonement.Title,
			Validity:      abonementWithServices.Abonement.Validity,
			VisitingTime:  abonementWithServices.Abonement.VisitingTime,
			Photo:         abonementWithServices.Abonement.Photo,
			Price:         int32(abonementWithServices.Abonement.Price),
			CreatedTime:   abonementWithServices.Abonement.CreatedTime.String(),
			UpdatedTime:   abonementWithServices.Abonement.UpdatedTime.String(),
			StripePriceId: abonementWithServices.Abonement.StripePriceId,
		}

		abonementWithServices := &abonementProtobuf.AbonementWithServices{
			Abonement: abonementObject,
			Services:  abonementWithServices.Services,
		}

		abonementsWithServicesForResponse = append(abonementsWithServicesForResponse, abonementWithServices)
	}

	response := &abonementProtobuf.GetAbonementsWithServicesResponse{
		AbonementsWithServices: abonementsWithServicesForResponse,
	}

	return response, nil
}

func (c *AbonementgRPC) GetAbonementsByIds(ctx context.Context, request *abonementProtobuf.GetAbonementsByIdsRequest) (*abonementProtobuf.GetAbonementsByIdsResponse, error) {

	var ids []uuid.UUID
	for _, id := range request.Ids {
		ids = append(ids, uuid.MustParse(id))
	}

	abonements, err := c.abonementUseCase.GetAbonementsByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	getAbonementsByIdsResponse := &abonementProtobuf.GetAbonementsByIdsResponse{
		AbonementObjects: nil,
	}
	for _, abonement := range abonements {
		abonementObject := &abonementProtobuf.AbonementObject{
			Id:            abonement.Id.String(),
			Title:         abonement.Title,
			Validity:      abonement.Validity,
			VisitingTime:  abonement.VisitingTime,
			Photo:         abonement.Photo,
			Price:         int32(abonement.Price),
			CreatedTime:   abonement.CreatedTime.String(),
			UpdatedTime:   abonement.UpdatedTime.String(),
			StripePriceId: abonement.StripePriceId,
		}

		getAbonementsByIdsResponse.AbonementObjects = append(getAbonementsByIdsResponse.AbonementObjects, abonementObject)
	}

	return getAbonementsByIdsResponse, nil
}

func GetObjectData[T any, R any](
	g *grpc.ClientStreamingServer[T, R],
	extractObjectData func(chunk *T) interface{},
	extractObjectPhoto func(chunk *T) []byte,
) (interface{},
	[]byte,
	error,
) {
	var objectData interface{}
	var objectPhoto []byte

	for {
		chunk, err := (*g).Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.ErrorLogger.Printf("Error getting chunk: %v", err)
			return nil, nil, err
		}

		if ud := extractObjectData(chunk); ud != nil && !reflect.ValueOf(ud).IsNil() {
			objectData = ud
		}

		if uf := extractObjectPhoto(chunk); uf != nil {
			objectPhoto = append(objectPhoto, uf...)
		}
	}

	return objectData, objectPhoto, nil
}
