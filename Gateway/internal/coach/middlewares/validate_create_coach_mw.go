package middlewares

import (
	"Gateway/internal/coach/coach_errors"
	"Gateway/internal/coach/dtos"
	"Gateway/pkg/logger"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidateCreateCoachMW() gin.HandlerFunc {
	return func(c *gin.Context) {

		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
			return
		}

		name := form.Value["name"]
		description := form.Value["description"]
		services := form.Value["services"]
		shift := form.Value["shift"]
		email := form.Value["email"]

		if (len(name) != 1 || name[0] == "") ||
			(len(description) != 1 || description[0] == "") ||
			(len(services) != 1 || services[0] == "") ||
			(len(shift) != 1 || shift[0] == "") ||
			(len(email) != 1 || email[0] == "") {
			logger.ErrorLogger.Printf(coach_errors.OnlyPhotoOptional.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": coach_errors.OnlyPhotoOptional.Error()})
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

		//services validation
		servicesIds := strings.Split(services[0], ",")
		if len(servicesIds) == 0 {
			logger.ErrorLogger.Printf("at least one service is required")
			c.JSON(http.StatusBadRequest, gin.H{"error": "требуется по крайней мере одна услуга"})
			return
		}

		//shift
		shiftValue := shift[0]
		if !(shiftValue == "1" || shiftValue == "2") {
			logger.ErrorLogger.Printf("not valid schema")
			c.JSON(http.StatusBadRequest, gin.H{"error": "невалидная смена"})
			return
		}

		// email validation

		emailValue := email[0]
		slog.Info("emailValue", emailValue)
		emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		matched, _ = regexp.MatchString(emailRegex, emailValue)
		if !matched {
			logger.ErrorLogger.Printf("Invalid email format")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат email"})
			return
		}

		createCoachCommand := &dtos.CreateCoachCommand{
			Name:        nameValue,
			Description: descriptionValue,
			Services:    servicesIds,
			Shift:       shiftValue,
			Email:       emailValue,
		}

		c.Set("CreateCoachCommand", createCoachCommand)

		c.Next()
	}
}
