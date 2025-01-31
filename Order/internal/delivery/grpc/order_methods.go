package grpc

import (
	"context"
	"github.com/DanKo-code/FitnessCenter-Order/internal/dtos"
	"github.com/DanKo-code/FitnessCenter-Order/internal/usecase"
	"github.com/DanKo-code/FitnessCenter-Order/pkg/logger"
	orderProtobuf "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.order"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"google.golang.org/grpc"
	"log"
)

type OrdergRPC struct {
	orderProtobuf.UnimplementedOrderServer

	orderUseCase usecase.OrderUseCase
}

func Register(gRPC *grpc.Server, orderUseCase usecase.OrderUseCase) {
	orderProtobuf.RegisterOrderServer(gRPC, &OrdergRPC{orderUseCase: orderUseCase})
}

func (o OrdergRPC) CreateOrder(ctx context.Context, request *orderProtobuf.CreateOrderRequest) (*orderProtobuf.CreateOrderResponse, error) {

	logger.DebugLogger.Printf("request.OrderDataForCreate: %v", request.OrderDataForCreate)

	cmd := &dtos.CreateOrderCommand{
		UserId:      uuid.MustParse(request.OrderDataForCreate.UserId),
		AbonementId: uuid.MustParse(request.OrderDataForCreate.AbonementId),
	}

	order, err := o.orderUseCase.CreateCoachAbonement(ctx, cmd)
	if err != nil {
		return nil, err
	}

	orderObject := &orderProtobuf.OrderObject{
		Id:          order.Id.String(),
		UserId:      order.UserId.String(),
		AbonementId: order.AbonementId.String(),
		Status:      order.Status,
		CreatedTime: order.CreatedTime.String(),
		UpdatedTime: order.UpdatedTime.String(),
	}

	response := &orderProtobuf.CreateOrderResponse{
		OrderObject: orderObject,
	}

	return response, nil
}

func (o OrdergRPC) GetUserOrders(ctx context.Context, request *orderProtobuf.GetUserOrdersRequest) (*orderProtobuf.GetUserOrdersResponse, error) {
	response, err := o.orderUseCase.GetUserOrders(ctx, uuid.MustParse(request.UserId))
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (o OrdergRPC) CreateCheckoutSession(ctx context.Context, request *orderProtobuf.CreateCheckoutSessionRequest) (*orderProtobuf.CreateCheckoutSessionResponse, error) {

	stripe.Key = "sk_test_51PxOL0A75DCPwyUvr31hX8Ju84gJa8CuRgT2o7RA5eRVfhPSwtRRmfpxVYbPkCpSNSF4xPytvonrhMq7qZtkeewb00SO5G61FT"

	domain := "http://localhost:3333/main/abonnements"
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				Price:    stripe.String(request.StripePriceId),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(domain),
		CancelURL:  stripe.String(domain),

		Metadata: map[string]string{
			"client_id":    request.ClientId,
			"abonement_id": request.AbonementId,
		},
	}

	s, err := session.New(params)
	if err != nil {
		log.Printf("session.New: %v", err)
	}

	createCheckoutSessionResponse := &orderProtobuf.CreateCheckoutSessionResponse{
		SessionUrl: s.ID,
	}

	return createCheckoutSessionResponse, nil
}
