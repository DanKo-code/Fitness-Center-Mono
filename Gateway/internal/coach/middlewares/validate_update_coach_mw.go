package middlewares

import (
	"Gateway/internal/coach/dtos"
	"Gateway/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"regexp"
	"strings"
)

func ValidateUpdateCoachMW() gin.HandlerFunc {
	return func(c *gin.Context) {

		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
			return
		}

		id, idOk := form.Value["id"]
		name, nameOk := form.Value["name"]
		description, descriptionOk := form.Value["description"]
		services, servicesOk := form.Value["services"]

		logger.DebugLogger.Printf(
			"id: %v\n"+
				"name: %v\n"+
				"description: %v\n"+
				"services: %v\n",
			id, name, description, services,
		)

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

		if (len(name) != 1 || !nameOk) &&
			(len(description) != 1 || !descriptionOk) &&
			(len(services) != 1 || !servicesOk) {
			logger.ErrorLogger.Printf("at least 1 field is required for updating")
			c.JSON(http.StatusBadRequest, gin.H{"error": "для обновления требуется как минимум 1 поле"})
			c.Set("InvalidUpdate", struct{}{})
			return
		}

		//name validation
		nameValue := name[0]
		if len(nameValue) < 2 || len(nameValue) > 100 {
			logger.ErrorLogger.Printf("Name must be between 2 and 100 characters long")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Длина имени должна составлять от 2 до 100 символов"})
			return
		}
		allowedNameRegex := `^[a-zA-Zа-яА-Я0-9 ]+$`
		matched, _ := regexp.MatchString(allowedNameRegex, nameValue)
		if !matched {
			logger.ErrorLogger.Printf("Name can only contain Russian and English letters, digits, and spaces")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Name can only contain Russian and English letters, digits, and spaces"})
			return
		}

		//description validation
		descriptionValue := description[0]
		if len(descriptionValue) < 10 || len(nameValue) > 500 {
			logger.ErrorLogger.Printf("Description must be between 10 and 500 characters long")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Длина описания должна составлять от 10 до 500 символов"})
			return
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

		updateCoachCommand := &dtos.UpdateCoachCommand{}

		updateCoachCommand.Id = id[0]

		if nameOk {
			updateCoachCommand.Name = nameValue
		}
		if descriptionOk {
			updateCoachCommand.Description = descriptionValue
		}
		if servicesOk {
			updateCoachCommand.Services = servicesIds
		}

		c.Set("UpdateCoachCommand", updateCoachCommand)

		c.Next()
	}
}
