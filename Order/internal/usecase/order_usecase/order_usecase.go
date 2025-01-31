package order_usecase

import (
	"context"
	"github.com/DanKo-code/FitnessCenter-Order/internal/dtos"
	"github.com/DanKo-code/FitnessCenter-Order/internal/models"
	"github.com/DanKo-code/FitnessCenter-Order/internal/repository"
	"github.com/DanKo-code/FitnessCenter-Order/pkg/logger"
	abonementGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.abonement"
	orderGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.order"
	serviceGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.service"
	userGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type OrderUseCase struct {
	orderRepo       repository.OrderRepository
	abonementClient *abonementGRPC.AbonementClient
	serviceClient   *serviceGRPC.ServiceClient
	userClient      *userGRPC.UserClient
}

func NewOrderUseCase(
	orderRepo repository.OrderRepository,
	abonementClient *abonementGRPC.AbonementClient,
	userClient *userGRPC.UserClient,
	serviceClient *serviceGRPC.ServiceClient,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:       orderRepo,
		abonementClient: abonementClient,
		userClient:      userClient,
		serviceClient:   serviceClient,
	}
}

func (o *OrderUseCase) CreateCoachAbonement(ctx context.Context, cmd *dtos.CreateOrderCommand) (*models.Order, error) {

	getUserByIdRequest := &userGRPC.GetUserByIdRequest{Id: cmd.UserId.String()}

	_, err := (*o.userClient).GetUserById(ctx, getUserByIdRequest)
	if err != nil {
		return nil, err
	}

	getAbonementByIdRequest := &abonementGRPC.GetAbonementByIdRequest{Id: cmd.AbonementId.String()}

	abonement, err := (*o.abonementClient).GetAbonementById(ctx, getAbonementByIdRequest)
	if err != nil {
		return nil, err
	}

	order := &models.Order{
		Id:          uuid.New(),
		AbonementId: cmd.AbonementId,
		UserId:      cmd.UserId,
		Status:      "Valid",
		UpdatedTime: time.Now(),
		CreatedTime: time.Now(),
	}

	number, err := strconv.Atoi(abonement.AbonementObject.Validity)
	if err != nil {
		return nil, err
	}

	expiredTime := order.CreatedTime.AddDate(0, number, 0)

	order.ExpiredTime = expiredTime

	err = o.orderRepo.CreateCoachAbonement(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *OrderUseCase) GetUserOrders(ctx context.Context, userId uuid.UUID) (*orderGRPC.GetUserOrdersResponse, error) {

	getUserByIdRequest := &userGRPC.GetUserByIdRequest{Id: userId.String()}

	_, err := (*o.userClient).GetUserById(ctx, getUserByIdRequest)
	if err != nil {
		return nil, err
	}

	// models orders
	orders, err := o.orderRepo.GetUserOrders(ctx, userId)
	if err != nil {
		return nil, err
	}

	logger.DebugLogger.Printf("user orders: %v\n", orders)

	var abonementIds []string
	for _, abonement := range orders {
		abonementIds = append(abonementIds, abonement.AbonementId.String())
	}

	// proto abonements
	abonements := &abonementGRPC.GetAbonementsByIdsResponse{}
	if len(abonementIds) > 0 {

		getAbonementsByIdsRequest := &abonementGRPC.GetAbonementsByIdsRequest{
			Ids: abonementIds,
		}

		abonements, err = (*o.abonementClient).GetAbonementsByIds(context.TODO(), getAbonementsByIdsRequest)
		if err != nil {
			return nil, err
		}
	}

	abonIdAbonement := map[string]*abonementGRPC.AbonementObject{}
	for _, abonement := range abonements.AbonementObjects {
		abonIdAbonement[abonement.Id] = abonement
	}

	logger.DebugLogger.Printf("user orders abonements: %v\n", abonements)

	// proto abonements with services
	abonIdServicesMap := map[string][]*serviceGRPC.ServiceObject{}
	if len(abonementIds) > 0 {

		var abonIds []string
		for _, abon := range abonements.AbonementObjects {
			abonIds = append(abonIds, abon.Id)
		}

		getAbonementsServicesRequest := &serviceGRPC.GetAbonementsServicesRequest{
			AbonementIds: abonIds,
		}

		abonIdServices, err := (*o.serviceClient).GetAbonementsServices(ctx, getAbonementsServicesRequest)
		if err != nil {
			return nil, err
		}

		for _, object := range abonIdServices.AbonementIdsWithServices {
			abonIdServicesMap[object.AbonementId] = object.ServiceObjects
		}
	}

	logger.DebugLogger.Printf("user abonIdServicesMap: %v\n", abonIdServicesMap)

	orderIdOrder := map[uuid.UUID]*orderGRPC.OrderObject{}
	for _, order := range orders {
		orderObject := &orderGRPC.OrderObject{
			Id:          order.Id.String(),
			UserId:      order.UserId.String(),
			AbonementId: order.AbonementId.String(),
			Status:      order.Status,
			CreatedTime: order.CreatedTime.String(),
			UpdatedTime: order.UpdatedTime.String(),
		}

		orderIdOrder[order.Id] = orderObject
	}

	getUserOrdersResponse := orderGRPC.GetUserOrdersResponse{
		OrderObjectWithAbonementWithServices: nil,
	}

	for _, order := range orders {
		orderObjectWithAbonement := &orderGRPC.OrderObjectWithAbonementWithServices{
			OrderObject:     orderIdOrder[order.Id],
			AbonementObject: abonIdAbonement[order.AbonementId.String()],
			ServiceObjects:  abonIdServicesMap[order.AbonementId.String()],
		}

		getUserOrdersResponse.OrderObjectWithAbonementWithServices = append(getUserOrdersResponse.OrderObjectWithAbonementWithServices, orderObjectWithAbonement)
	}

	return &getUserOrdersResponse, nil
}

func (o *OrderUseCase) SetExpiredOrdersTasks(ctx context.Context) error {
	err := o.orderRepo.SetExpiredOrdersTasks(ctx)
	if err != nil {
		return err
	}

	return nil
}
