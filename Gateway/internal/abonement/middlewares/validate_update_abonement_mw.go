package middlewares

import (
	"Gateway/internal/abonement/dtos"
	"Gateway/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func ValidateUpdateAbonementMW() gin.HandlerFunc {
	return func(c *gin.Context) {

		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
			return
		}

		id, idOk := form.Value["id"]
		title, titleOk := form.Value["title"]
		validityPeriod, validityPeriodOk := form.Value["validity_period"]
		visitingTime, visitingTimeOk := form.Value["visiting_time"]
		price, priceOk := form.Value["price"]
		services, servicesOk := form.Value["services"]

		//id validate
		if len(id) != 1 || !idOk {
			logger.ErrorLogger.Printf("id is required for updating")
			c.JSON(http.StatusBadRequest, gin.H{"error": "для обновления требуется идентификатор"})
			c.Set("InvalidUpdate", struct{}{})
			return
		}
		_, err = uuid.Parse(id[0])
		if err != nil {
			logger.ErrorLogger.Printf("id must be uuid")
			c.JSON(http.StatusBadRequest, gin.H{"error": "id must be uuid"})
			c.Set("InvalidUpdate", struct{}{})
			return
		}

		if (len(title) != 1 || !titleOk) &&
			(len(validityPeriod) != 1 || !validityPeriodOk) &&
			(len(visitingTime) != 1 || !visitingTimeOk) &&
			(len(price) != 1 || !priceOk) &&
			(len(services) != 1 || !servicesOk) {
			logger.ErrorLogger.Printf("at least 1 field is required for updating")
			c.JSON(http.StatusBadRequest, gin.H{"error": "для обновления требуется как минимум 1 поле"})
			c.Set("InvalidUpdate", struct{}{})
			return
		}

		//validityPeriod validation
		if validityPeriodOk {
			timeValue, err := strconv.Atoi(validityPeriod[0])
			if err != nil || timeValue < 1 || timeValue > 12 {
				logger.ErrorLogger.Printf("Visiting time must be a number between 1 and 12")
				c.JSON(http.StatusBadRequest, gin.H{"error": "Visiting time must be a number between 1 and 12"})
				c.Set("InvalidUpdate", struct{}{})
				return
			}
		}

		//title validation
		titleValue := title[0]
		if titleOk {
			if len(titleValue) < 3 || len(titleValue) > 100 {
				logger.ErrorLogger.Printf("Title must be between 3 and 100 characters long")
				c.JSON(http.StatusBadRequest, gin.H{"error": "Длина названия должна составлять от 3 до 100 символов"})
				c.Set("InvalidUpdate", struct{}{})
				return
			}
			allowedTitleRegex := `^[a-zA-Zа-яА-Я0-9 ]+$`
			matched, _ := regexp.MatchString(allowedTitleRegex, titleValue)
			if !matched {
				logger.ErrorLogger.Printf("Title can only contain Russian and English letters, digits, and spaces")
				c.JSON(http.StatusBadRequest, gin.H{"error": "Title can only contain Russian and English letters, digits, and spaces"})
				c.Set("InvalidUpdate", struct{}{})
				return
			}
		}

		//price validation
		var parsePrice int64
		if priceOk {
			parsePrice, err = strconv.ParseInt(price[0], 10, 32)
			if err != nil {
				logger.ErrorLogger.Printf("Failed Parse price to Int: %s", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "price must be a number"})
				c.Set("InvalidUpdate", struct{}{})
				return
			}
			if parsePrice < 10 || parsePrice > 1000 {
				logger.ErrorLogger.Printf("price must bigger then 9$ and lower then 1000$")
				c.JSON(http.StatusBadRequest, gin.H{"error": "цена должна быть ниже 10BYN и выше 1000BYN"})
				c.Set("InvalidUpdate", struct{}{})
				return
			}
		}

		var servicesIds []string
		if servicesOk {
			servicesIds = strings.Split(services[0], ",")

			for _, serId := range servicesIds {
				_, err = uuid.Parse(serId)
				if err != nil {
					logger.ErrorLogger.Printf("service id must be uuid")
					c.JSON(http.StatusBadRequest, gin.H{"error": "service id must be uuid"})
					c.Set("InvalidUpdate", struct{}{})
					return
				}
			}
		}

		updateAbonementCommand := &dtos.UpdateAbonementCommand{}

		updateAbonementCommand.Id = id[0]

		if titleOk {
			updateAbonementCommand.Title = titleValue
		}
		if validityPeriodOk {
			updateAbonementCommand.ValidityPeriod = validityPeriod[0]
		}
		if visitingTimeOk {
			updateAbonementCommand.VisitingTime = visitingTime[0]
		}
		if priceOk {
			updateAbonementCommand.Price = parsePrice
		}
		if servicesOk {
			updateAbonementCommand.Services = servicesIds
		}

		c.Set("UpdateAbonementCommand", updateAbonementCommand)

		c.Next()
	}
}
