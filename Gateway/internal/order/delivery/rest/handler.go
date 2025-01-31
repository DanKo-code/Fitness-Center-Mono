package rest

import (
	"Gateway/internal/order/dtos"
	"Gateway/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	orderGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.order"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/webhook"
	"io/ioutil"
	"log"
	"net/http"
)

const endpointSecret = "whsec_07715d3864ab33c2cb0d10293ce78d99e311a165f98a15a40f4c1cd2711b286d"

type Handler struct {
	orderClient *orderGRPC.OrderClient
}

func NewHandler(orderClient *orderGRPC.OrderClient) *Handler {
	return &Handler{
		orderClient: orderClient,
	}
}

func (h *Handler) HandleCheckoutSessionCompleted(c *gin.Context) {

	logger.InfoLogger.Printf("Enter HandleCheckoutSessionCompleted Handler")

	payload, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read body"})
		return
	}

	logger.DebugLogger.Printf("Stripe-Signature Header: " + c.Request.Header.Get("Stripe-Signature"))

	event, err := webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"), endpointSecret)
	if err != nil {
		log.Printf("Webhook signature verification failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}

	var clientId string
	var abonementId string
	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			log.Printf("Failed to parse session: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse event"})
			return
		}

		logger.DebugLogger.Printf("session.Metadata: %v", session.Metadata)

		clientId = session.Metadata["client_id"]
		abonementId = session.Metadata["abonement_id"]

	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	createOrderRequest := &orderGRPC.CreateOrderRequest{
		OrderDataForCreate: &orderGRPC.OrderDataForCreate{
			UserId:      clientId,
			AbonementId: abonementId,
		},
	}

	logger.DebugLogger.Printf("createOrderRequest: %v", createOrderRequest)

	order, err := (*h.orderClient).CreateOrder(context.TODO(), createOrderRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order": order.GetOrderObject(),
	})
}

func (h *Handler) CreateCheckoutSession(c *gin.Context) {

	ccsDto := &dtos.CreateCheckoutSessionDTO{}

	if err := c.ShouldBindJSON(&ccsDto); err != nil {
		logger.ErrorLogger.Printf("Error binding CreateCheckoutSessionRequest: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createCheckoutSessionRequest := &orderGRPC.CreateCheckoutSessionRequest{
		ClientId:      ccsDto.ClientId,
		AbonementId:   ccsDto.AbonementId,
		StripePriceId: ccsDto.StripePriceId,
	}

	logger.DebugLogger.Printf("createCheckoutSessionRequest: %v", createCheckoutSessionRequest)

	session, err := (*h.orderClient).CreateCheckoutSession(context.TODO(), createCheckoutSessionRequest)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessionUrl": session.GetSessionUrl()})
}

func (h *Handler) GetUserOrders(c *gin.Context) {
	id := c.Param("userId")

	convertedId, err := uuid.Parse(id)
	if err != nil {

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid id format")})
		return
	}

	_ = convertedId

	getUserOrdersRequest := &orderGRPC.GetUserOrdersRequest{
		UserId: id,
	}

	orders, err := (*h.orderClient).GetUserOrders(context.TODO(), getUserOrdersRequest)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders.OrderObjectWithAbonementWithServices,
	})
}
