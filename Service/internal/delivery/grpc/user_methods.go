package grpc

import (
	"Service/internal/dtos"
	customErrors "Service/internal/errors"
	"Service/internal/usecase"
	"Service/pkg/logger"
	"context"
	"errors"
	serviceProtobuf "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.service"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"time"
)

type ServicegRPC struct {
	serviceProtobuf.UnimplementedServiceServer

	ServiceUseCase usecase.ServiceUseCase
	cloudUseCase   usecase.CloudUseCase
}

func Register(gRPC *grpc.Server, ServiceUseCase usecase.ServiceUseCase, cloudUseCase usecase.CloudUseCase) {
	serviceProtobuf.RegisterServiceServer(gRPC, &ServicegRPC{ServiceUseCase: ServiceUseCase, cloudUseCase: cloudUseCase})
}

func (u *ServicegRPC) CreateService(
	g grpc.ClientStreamingServer[
		serviceProtobuf.CreateServiceRequest,
		serviceProtobuf.CreateServiceResponse,
	]) error {

	serviceData, servicePhoto, err := GetObjectData(
		&g,
		func(chunk *serviceProtobuf.CreateServiceRequest) interface{} {
			return chunk.GetServiceDataForCreate()
		},
		func(chunk *serviceProtobuf.CreateServiceRequest) []byte {
			return chunk.GetServicePhoto()
		},
	)
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid request data")
	}

	if serviceData == nil {
		logger.ErrorLogger.Printf("service data is empty")
		return status.Error(codes.InvalidArgument, "service data is empty")
	}

	castedServiceData, ok := serviceData.(*serviceProtobuf.ServiceDataForCreate)
	if !ok {
		logger.ErrorLogger.Printf("service data is not of type ServiceProtobuf.ServiceDataForCreate")
		return status.Error(codes.InvalidArgument, "service data is not of type ServiceProtobuf.ServiceDataForCreate")
	}

	cmd := &dtos.CreateServiceCommand{
		Id:    uuid.New(),
		Title: castedServiceData.Title,
		Photo: "",
	}

	var photoURL string
	randomID := uuid.New().String()
	if servicePhoto != nil {
		url, err := u.cloudUseCase.PutObject(context.TODO(), servicePhoto, "service/"+randomID)
		photoURL = url
		if err != nil {
			logger.ErrorLogger.Printf("Failed to create service photo in cloud: %v", err)
			return status.Error(codes.Internal, "Failed to create service photo in cloud")
		}
	}

	cmd.Photo = photoURL

	service, err := u.ServiceUseCase.CreateService(context.TODO(), cmd)
	if err != nil {
		return err
	}

	serviceObject := &serviceProtobuf.ServiceObject{
		Id:    service.Id.String(),
		Title: service.Title,
		Photo: service.Photo,
	}

	response := &serviceProtobuf.CreateServiceResponse{
		ServiceObject: serviceObject,
	}

	err = g.SendAndClose(response)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to send service create response: %v", err)
		return status.Error(codes.Internal, "Failed to send service create response")
	}

	return nil
}

func (u *ServicegRPC) GetServiceById(ctx context.Context,
	request *serviceProtobuf.GetServiceByIdRequest,
) (*serviceProtobuf.GetServiceByIdResponse, error) {

	service, err := u.ServiceUseCase.GetServiceById(ctx, uuid.MustParse(request.Id))
	if err != nil {

		if errors.Is(err, customErrors.ServiceNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, err
	}

	serviceObject := &serviceProtobuf.ServiceObject{
		Id:          service.Id.String(),
		Title:       service.Title,
		Photo:       service.Photo,
		CreatedTime: service.CreatedTime.String(),
		UpdatedTime: service.UpdatedTime.String(),
	}

	response := &serviceProtobuf.GetServiceByIdResponse{
		ServiceObject: serviceObject,
	}

	return response, nil
}

func (u *ServicegRPC) UpdateService(
	g grpc.ClientStreamingServer[serviceProtobuf.UpdateServiceRequest, serviceProtobuf.UpdateServiceResponse],
) error {

	serviceData, servicePhoto, err := GetObjectData(
		&g,
		func(chunk *serviceProtobuf.UpdateServiceRequest) interface{} {
			return chunk.GetServiceDataForUpdate()
		},
		func(chunk *serviceProtobuf.UpdateServiceRequest) []byte {
			return chunk.GetServicePhoto()
		},
	)
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid request data")
	}

	if serviceData == nil {
		logger.ErrorLogger.Printf("service data is empty")
		return status.Error(codes.InvalidArgument, "service data is empty")
	}

	castedServiceData, ok := serviceData.(*serviceProtobuf.ServiceDataForUpdate)
	if !ok {
		logger.ErrorLogger.Printf("service data is not of type ServiceProtobuf.ServiceDataForUpdate")
		return status.Error(codes.InvalidArgument, "service data is not of type ServiceProtobuf.ServiceDataForUpdate")
	}

	cmd := &dtos.UpdateServiceCommand{
		Id:          uuid.MustParse(castedServiceData.Id),
		Title:       castedServiceData.Title,
		UpdatedTime: time.Now(),
	}

	var photoURL string
	randomID := uuid.New().String()
	if servicePhoto != nil {
		url, err := u.cloudUseCase.PutObject(context.TODO(), servicePhoto, "service/"+randomID)
		photoURL = url
		if err != nil {
			logger.ErrorLogger.Printf("Failed to create service photo in cloud: %v", err)
			return status.Error(codes.Internal, "Failed to create service photo in cloud")
		}
	}

	cmd.Photo = photoURL

	service, err := u.ServiceUseCase.UpdateService(context.TODO(), cmd)
	if err != nil {
		return status.Error(codes.Internal, "Failed to create service")
	}

	serviceObject := &serviceProtobuf.ServiceObject{
		Id:          service.Id.String(),
		Title:       service.Title,
		Photo:       service.Photo,
		CreatedTime: service.CreatedTime.String(),
		UpdatedTime: service.UpdatedTime.String(),
	}

	response := &serviceProtobuf.UpdateServiceResponse{
		ServiceObject: serviceObject,
	}

	err = g.SendAndClose(response)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to send Service update response: %v", err)
		return err
	}

	return nil
}

func (u *ServicegRPC) DeleteServiceById(
	ctx context.Context,
	request *serviceProtobuf.DeleteServiceByIdRequest,
) (*serviceProtobuf.DeleteServiceByIdResponse, error) {

	service, err := u.ServiceUseCase.DeleteServiceById(ctx, uuid.MustParse(request.Id))
	if err != nil {
		return nil, err
	}

	response := &serviceProtobuf.DeleteServiceByIdResponse{ServiceObject: &serviceProtobuf.ServiceObject{
		Id:          service.Id.String(),
		Title:       service.Title,
		Photo:       service.Photo,
		CreatedTime: service.CreatedTime.String(),
		UpdatedTime: service.UpdatedTime.String(),
	}}

	return response, nil
}

func (u *ServicegRPC) GetServices(
	ctx context.Context,
	_ *emptypb.Empty,
) (*serviceProtobuf.GetServicesResponse, error) {

	services, err := u.ServiceUseCase.GetServices(ctx)
	if err != nil {
		return nil, err
	}

	var serviceObjects []*serviceProtobuf.ServiceObject

	for _, service := range services {
		serviceObjects = append(serviceObjects, &serviceProtobuf.ServiceObject{
			Id:          service.Id.String(),
			Title:       service.Title,
			Photo:       service.Photo,
			CreatedTime: service.CreatedTime.String(),
			UpdatedTime: service.UpdatedTime.String(),
		})
	}

	response := &serviceProtobuf.GetServicesResponse{ServiceObject: serviceObjects}

	return response, nil
}

func (u *ServicegRPC) CreateCoachServices(
	ctx context.Context,
	request *serviceProtobuf.CreateCoachServicesRequest,
) (*serviceProtobuf.CreateCoachServicesResponse, error) {

	var servicesIds []uuid.UUID
	for _, serviceId := range request.CoachService.ServiceId {
		servicesIds = append(servicesIds, uuid.MustParse(serviceId))
	}

	cmd := &dtos.CreateCoachServicesCommand{
		CoachId:     uuid.MustParse(request.CoachService.CoachId),
		ServicesIds: servicesIds,
	}

	coachServices, err := u.ServiceUseCase.CreateCoachServices(ctx, cmd)
	if err != nil {
		return nil, err
	}

	var coachServicesIds []string
	for _, coachService := range coachServices {
		coachServicesIds = append(coachServicesIds, coachService.Id.String())
	}

	coachService := &serviceProtobuf.CoachService{
		CoachId:   request.CoachService.CoachId,
		ServiceId: coachServicesIds,
	}

	createCoachServicesResponse := &serviceProtobuf.CreateCoachServicesResponse{
		CoachService: coachService,
	}

	return createCoachServicesResponse, nil
}

func (u *ServicegRPC) CreateAbonementServices(
	ctx context.Context,
	request *serviceProtobuf.CreateAbonementServicesRequest,
) (*serviceProtobuf.CreateAbonementServicesResponse, error) {

	var servicesIds []uuid.UUID
	for _, serviceId := range request.AbonementService.ServiceId {
		servicesIds = append(servicesIds, uuid.MustParse(serviceId))
	}

	cmd := &dtos.CreateAbonementServicesCommand{
		AbonementId: uuid.MustParse(request.AbonementService.AbonementId),
		ServicesIds: servicesIds,
	}

	abonementServices, err := u.ServiceUseCase.CreateAbonemntServices(ctx, cmd)
	if err != nil {
		return nil, err
	}

	var abonementServicesIds []string
	for _, abonementService := range abonementServices {
		abonementServicesIds = append(abonementServicesIds, abonementService.Id.String())
	}

	abonementService := &serviceProtobuf.AbonementService{
		AbonementId: request.AbonementService.AbonementId,
		ServiceId:   abonementServicesIds,
	}

	createAbonementServicesResponse := &serviceProtobuf.CreateAbonementServicesResponse{
		AbonementService: abonementService,
	}

	return createAbonementServicesResponse, nil
}

func (u *ServicegRPC) GetAbonementsServices(
	ctx context.Context,
	request *serviceProtobuf.GetAbonementsServicesRequest,
) (*serviceProtobuf.GetAbonementsServicesResponse, error) {

	var abonementIdsUUID []uuid.UUID
	for _, id := range request.AbonementIds {
		abonementIdsUUID = append(abonementIdsUUID, uuid.MustParse(id))
	}

	abonementIdWithServicesResponse, err := u.ServiceUseCase.GetAbonementsServices(ctx, abonementIdsUUID)
	if err != nil {
		return nil, err
	}

	getAbonementsServicesResponse := &serviceProtobuf.GetAbonementsServicesResponse{}

	for ai, aiws := range abonementIdWithServicesResponse {

		abonementIdWithServices := &serviceProtobuf.AbonementIdWithServices{
			AbonementId:    ai.String(),
			ServiceObjects: nil,
		}

		for _, service := range aiws {

			serviceObject := &serviceProtobuf.ServiceObject{
				Id:          service.Id.String(),
				Title:       service.Title,
				Photo:       service.Photo,
				CreatedTime: service.CreatedTime.String(),
				UpdatedTime: service.UpdatedTime.String(),
			}

			abonementIdWithServices.ServiceObjects = append(abonementIdWithServices.ServiceObjects, serviceObject)
		}

		getAbonementsServicesResponse.AbonementIdsWithServices =
			append(
				getAbonementsServicesResponse.AbonementIdsWithServices,
				abonementIdWithServices,
			)
	}

	return getAbonementsServicesResponse, nil
}

func (u *ServicegRPC) GetCoachesServices(
	ctx context.Context,
	request *serviceProtobuf.GetCoachesServicesRequest,
) (*serviceProtobuf.GetCoachesServicesResponse, error) {

	var coachIdsUUID []uuid.UUID
	for _, id := range request.CoachIds {
		coachIdsUUID = append(coachIdsUUID, uuid.MustParse(id))
	}

	coachIdWithServicesResponse, err := u.ServiceUseCase.GetCoachesServices(ctx, coachIdsUUID)
	if err != nil {
		return nil, err
	}

	getCoachesServicesResponse := &serviceProtobuf.GetCoachesServicesResponse{}
	for ai, aiws := range coachIdWithServicesResponse {

		coachIdWithServices := &serviceProtobuf.CoachIdWithServices{
			CoachId:        ai.String(),
			ServiceObjects: nil,
		}

		for _, service := range aiws {

			serviceObject := &serviceProtobuf.ServiceObject{
				Id:          service.Id.String(),
				Title:       service.Title,
				Photo:       service.Photo,
				CreatedTime: service.CreatedTime.String(),
				UpdatedTime: service.UpdatedTime.String(),
			}

			coachIdWithServices.ServiceObjects = append(coachIdWithServices.ServiceObjects, serviceObject)
		}

		getCoachesServicesResponse.CoachIdsWithServices =
			append(
				getCoachesServicesResponse.CoachIdsWithServices,
				coachIdWithServices,
			)
	}

	return getCoachesServicesResponse, nil
}

func (u *ServicegRPC) UpdateAbonementServices(
	ctx context.Context,
	request *serviceProtobuf.UpdateAbonementServicesRequest,
) (*serviceProtobuf.UpdateAbonementServicesResponse, error) {

	var servicesIds []uuid.UUID
	for _, serviceId := range request.AbonementService.ServiceId {
		servicesIds = append(servicesIds, uuid.MustParse(serviceId))
	}

	abonementServices, err := u.ServiceUseCase.UpdateAbonementServices(ctx, uuid.MustParse(request.AbonementService.AbonementId), servicesIds)
	if err != nil {
		return nil, err
	}

	var abonementServicesIds []string
	for _, abonementService := range abonementServices {
		abonementServicesIds = append(abonementServicesIds, abonementService.Id.String())
	}

	abonementService := &serviceProtobuf.AbonementService{
		AbonementId: request.AbonementService.AbonementId,
		ServiceId:   abonementServicesIds,
	}

	updateAbonementServicesResponse := &serviceProtobuf.UpdateAbonementServicesResponse{
		AbonementService: abonementService,
	}

	return updateAbonementServicesResponse, nil
}

func (u *ServicegRPC) UpdateCoachServices(
	ctx context.Context,
	request *serviceProtobuf.UpdateCoachServicesRequest,
) (*serviceProtobuf.UpdateCoachServicesResponse, error) {
	var servicesIds []uuid.UUID
	for _, serviceId := range request.CoachService.ServiceId {
		servicesIds = append(servicesIds, uuid.MustParse(serviceId))
	}

	coachServices, err := u.ServiceUseCase.UpdateCoachServices(ctx, uuid.MustParse(request.CoachService.CoachId), servicesIds)
	if err != nil {
		return nil, err
	}

	var coachServicesIds []string
	for _, coachService := range coachServices {
		coachServicesIds = append(coachServicesIds, coachService.Id.String())
	}

	coachService := &serviceProtobuf.CoachService{
		CoachId:   request.CoachService.CoachId,
		ServiceId: coachServicesIds,
	}

	updateCoachServicesResponse := &serviceProtobuf.UpdateCoachServicesResponse{
		CoachService: coachService,
	}

	return updateCoachServicesResponse, nil
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

		if ud := extractObjectData(chunk); ud != nil {
			objectData = ud
		}

		if uf := extractObjectPhoto(chunk); uf != nil {
			objectPhoto = append(objectPhoto, uf...)
		}
	}

	return objectData, objectPhoto, nil
}
