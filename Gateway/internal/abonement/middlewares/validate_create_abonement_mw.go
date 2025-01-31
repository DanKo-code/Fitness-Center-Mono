package middlewares

import (
	"Gateway/internal/abonement/abonement_errors"
	"Gateway/internal/abonement/dtos"
	"Gateway/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func ValidateCreateAbonementMW() gin.HandlerFunc {
	return func(c *gin.Context) {

		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
			return
		}

		title := form.Value["title"]
		validityPeriod := form.Value["validity_period"]
		visitingTime := form.Value["visiting_time"]
		price := form.Value["price"]
		services := form.Value["services"]

		if (len(title) != 1 || title[0] == "") ||
			(len(validityPeriod) != 1 || validityPeriod[0] == "") ||
			(len(visitingTime) != 1 || visitingTime[0] == "") ||
			(len(price) != 1 || price[0] == "") ||
			(len(services) != 1 || services[0] == "") {
			logger.ErrorLogger.Printf(abonement_errors.OnlyPhotoOptional.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": abonement_errors.OnlyPhotoOptional.Error()})
			return
		}

		//validityPeriod validation
		timeValue, err := strconv.Atoi(validityPeriod[0])
		if err != nil || timeValue < 1 || timeValue > 12 {
			logger.ErrorLogger.Printf("Visiting time must be a number between 1 and 12")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Visiting time must be a number between 1 and 12"})
			return
		}

		//title validation
		titleValue := title[0]
		if len(titleValue) < 3 || len(titleValue) > 100 {
			logger.ErrorLogger.Printf("Title must be between 3 and 100 characters long")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Длина названия должна составлять от 3 до 100 символов"})
			return
		}
		allowedTitleRegex := `^[a-zA-Zа-яА-Я0-9 ]+$`
		matched, _ := regexp.MatchString(allowedTitleRegex, titleValue)
		if !matched {
			logger.ErrorLogger.Printf("Title can only contain Russian and English letters, digits, and spaces")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Title can only contain Russian and English letters, digits, and spaces"})
			return
		}

		//price validation
		parsePrice, err := strconv.ParseInt(price[0], 10, 32)
		if err != nil {
			logger.ErrorLogger.Printf("Failed Parse price to Int: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "price must be a number"})
			return
		}
		if parsePrice < 10 || parsePrice > 1000 {
			logger.ErrorLogger.Printf("price must bigger then 9$ and lower then 1000$")
			c.JSON(http.StatusBadRequest, gin.H{"error": "цена должна быть ниже 10BYN и выше 1000BYN"})
			return
		}

		//services validation
		servicesIds := strings.Split(services[0], ",")
		if len(servicesIds) == 0 {
			logger.ErrorLogger.Printf("at least one service is required")
			c.JSON(http.StatusBadRequest, gin.H{"error": "требуется по крайней мере одна услуга"})
			return
		}

		createAbonementCommand := &dtos.CreateAbonementCommand{
			Title:          titleValue,
			ValidityPeriod: validityPeriod[0],
			VisitingTime:   visitingTime[0],
			Price:          parsePrice,
			Services:       servicesIds,
		}

		c.Set("CreateAbonementCommand", createAbonementCommand)

		c.Next()
	}
}
