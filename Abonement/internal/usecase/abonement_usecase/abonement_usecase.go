package user_usecase

import (
	"context"
	"fmt"
	"github.com/DanKo-code/Fitness-Center-Abonement/internal/dtos"
	customErrors "github.com/DanKo-code/Fitness-Center-Abonement/internal/errors"
	"github.com/DanKo-code/Fitness-Center-Abonement/internal/models"
	"github.com/DanKo-code/Fitness-Center-Abonement/internal/repository"
	"github.com/DanKo-code/Fitness-Center-Abonement/internal/usecase"
	"github.com/DanKo-code/Fitness-Center-Abonement/pkg/logger"
	serviceGRPC "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.service"
	"github.com/google/uuid"
	"strings"
	"time"
)

type AbonementUseCase struct {
	abonementRepo repository.AbonementRepository
	serviceClient *serviceGRPC.ServiceClient
	stripeUseCase usecase.StripeUseCase
	cloudUseCase  usecase.CloudUseCase
}

func NewAbonementUseCase(
	abonementRepo repository.AbonementRepository,
	serviceClient *serviceGRPC.ServiceClient,
	stripeUseCase usecase.StripeUseCase,
	cloudUseCase usecase.CloudUseCase,
) *AbonementUseCase {
	return &AbonementUseCase{
		abonementRepo: abonementRepo,
		serviceClient: serviceClient,
		stripeUseCase: stripeUseCase,
		cloudUseCase:  cloudUseCase,
	}
}

func (c *AbonementUseCase) CreateAbonement(
	ctx context.Context,
	cmd *dtos.CreateAbonementCommand,
) (*models.Abonement, error) {

	abonement := &models.Abonement{
		Id:           cmd.Id,
		Title:        cmd.Title,
		Validity:     cmd.Validity,
		VisitingTime: cmd.VisitingTime,
		Photo:        cmd.Photo,
		Price:        cmd.Price,
		UpdatedTime:  time.Now(),
		CreatedTime:  time.Now(),
	}

	stripePriceId, err := c.stripeUseCase.CreateStripeProductAndPrice(abonement.Title, int64(abonement.Price*100), "byn")
	if err != nil {
		return nil, err
	}

	abonement.StripePriceId = stripePriceId

	createdAbonement, err := c.abonementRepo.CreateAbonement(ctx, abonement)
	if err != nil {
		return nil, err
	}

	return createdAbonement, nil
}

func (c *AbonementUseCase) GetAbonementById(
	ctx context.Context,
	id uuid.UUID,
) (*models.Abonement, error) {
	abonement, err := c.abonementRepo.GetAbonementById(ctx, id)
	if err != nil {
		return nil, err
	}

	return abonement, nil
}

func (c *AbonementUseCase) UpdateAbonement(
	ctx context.Context,
	cmd *dtos.UpdateAbonementCommand,
) (*models.Abonement, error) {

	existingAbonement, err := c.abonementRepo.GetAbonementById(ctx, cmd.Id)
	if err != nil {
		return nil, err
	}

	if cmd.Price != 0 {
		newStripePrices, err := c.stripeUseCase.CreateStripePriceAndAssignToProductDeactivateOldPrices(existingAbonement.StripePriceId, int64(cmd.Price*100), "byn")
		if err != nil {
			return nil, err
		}

		cmd.StripePriceId = newStripePrices
	}

	if cmd.Title != "" {
		err := c.stripeUseCase.UpdateStripeProductName(existingAbonement.StripePriceId, cmd.Title)
		if err != nil {
			return nil, err
		}
	}

	err = c.abonementRepo.UpdateAbonement(ctx, cmd)
	if err != nil {
		return nil, err
	}

	updatedAbonement, err := c.abonementRepo.GetAbonementById(ctx, cmd.Id)
	if err != nil {
		return nil, err
	}

	return updatedAbonement, nil
}

func (c *AbonementUseCase) DeleteAbonementById(
	ctx context.Context,
	id uuid.UUID,
) (*models.Abonement, error) {
	abonement, err := c.abonementRepo.GetAbonementById(ctx, id)
	if err != nil {
		return nil, customErrors.AbonementNotFound
	}

	err = c.abonementRepo.DeleteAbonementById(ctx, id)
	if err != nil {
		return nil, err
	}

	//TODO persists
	err = c.stripeUseCase.ArchiveStripeProduct(abonement.StripePriceId)
	if err != nil {
		return nil, err
	}

	if abonement.Photo != "" {
		prefix := "abonement/"
		index := strings.Index(abonement.Photo, prefix)
		var s3PhotoKey string
		if index != -1 {
			s3PhotoKey = abonement.Photo[index+len(prefix):]
		} else {
			logger.ErrorLogger.Printf("Prefix not found")
			return nil, fmt.Errorf("prefix not found")
		}
		err = c.cloudUseCase.DeleteObject(ctx, "abonement/"+s3PhotoKey)
		if err != nil {
			return nil, err
		}
	}

	return abonement, nil
}

func (c *AbonementUseCase) GetAbonementes(
	ctx context.Context,
) ([]*models.Abonement, error) {

	abonementes, err := c.abonementRepo.GetAbonementes(ctx)
	if err != nil {
		return nil, err
	}

	return abonementes, nil
}

func (c *AbonementUseCase) GetAbonementsWithServices(
	ctx context.Context,
) ([]*dtos.AbonementWithServices, error) {
	abonements, err := c.abonementRepo.GetAbonementes(ctx)
	if err != nil {
		logger.ErrorLogger.Printf("Failed GetAbonementes: %s", err)
		return nil, err
	}

	getAbonementsServicesRequest := &serviceGRPC.GetAbonementsServicesRequest{}

	for _, abonement := range abonements {
		getAbonementsServicesRequest.AbonementIds =
			append(
				getAbonementsServicesRequest.AbonementIds,
				abonement.Id.String(),
			)
	}

	getAbonementsServicesResponse, err := (*c.serviceClient).GetAbonementsServices(ctx, getAbonementsServicesRequest)
	if err != nil {
		logger.ErrorLogger.Printf("Failed GetAbonementsServices: %s", err)
		return nil, err
	}

	var abonementWithServices []*dtos.AbonementWithServices

	//add abonements
	for _, abonement := range abonements {

		aws := &dtos.AbonementWithServices{
			Abonement: abonement,
			Services:  nil,
		}

		abonementWithServices = append(abonementWithServices, aws)
	}

	//add services
	for _, extValue := range getAbonementsServicesResponse.AbonementIdsWithServices {

		abonementId := extValue.AbonementId

		for key, value := range abonementWithServices {
			if value.Abonement.Id.String() == abonementId {
				abonementWithServices[key].Services = append(abonementWithServices[key].Services, extValue.ServiceObjects...)
			}
		}
	}

	return abonementWithServices, nil
}

func (c *AbonementUseCase) GetAbonementsByIds(
	ctx context.Context,
	ids []uuid.UUID,
) ([]*models.Abonement, error) {
	abonements, err := c.abonementRepo.GetAbonementsByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	return abonements, err
}
