package utils

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func AlphaNumRegexCheck(value string) bool {
	isValid, _ := regexp.MatchString("^[A-Za-z0-9 ]+$", value)
	return isValid
}

func AlphaRegexCheck(value string) bool {
	isValid, _ := regexp.MatchString("^[A-Za-z ]+$", value)
	return isValid
}

// func NumRegexCheck(value string) bool {
// 	isValid, _ := regexp.MatchString("^[0-9.]+$", value)
// 	return isValid
// }

func CheckValidation(payload interface{}, logger *zap.Logger) ([]string, bool) {
	validate := validator.New()
	err := validate.RegisterValidation("UserNameCheck", CustomNameValidation)
	if err != nil {
		logger.Error(err.Error())
	}
	err = validate.RegisterValidation("AddressCheck", CustomAddressValidation)
	if err != nil {
		logger.Error(err.Error())
	}
	//err = validate.RegisterValidation("NumberCheck", CustomNumValidation)

	err = validate.Struct(payload)
	if err != nil {
		errMsg := make([]string, 0)
		var msg string
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Name":
				msg = fmt.Sprintf("The value '%v' is incorrect for name", err.Value())
				logger.Error("Name validation failed",
					zap.String("Tag", err.Tag()),
					zap.String("Type", err.Type().String()),
					zap.String("Issue", msg))
				errMsg = append(errMsg, msg)
			case "Email":
				msg = fmt.Sprintf("The value '%v' is incorrect for email", err.Value())
				logger.Error("Email validation failed",
					zap.String("Tag", err.Tag()),
					zap.String("Type", err.Type().String()),
					zap.String("Issue", msg))
				errMsg = append(errMsg, msg)
			case "Role":
				msg = fmt.Sprintf("The value '%v' is incorrect. It should be either of these values %v", err.Value(), err.Param())
				logger.Error("Role validation failed",
					zap.String("Tag", err.Tag()),
					zap.String("Type", err.Type().String()),
					zap.String("Issue", msg))
				errMsg = append(errMsg, msg)
			case "Password":
				msg = fmt.Sprintf("The length for password is incorrect. The length is %v", len(err.Value().(string)))
				logger.Error("Password validation failed",
					zap.String("Tag", err.Tag()),
					zap.String("Type", err.Type().String()),
					zap.String("Issue", msg))
				errMsg = append(errMsg, msg)
			case "Address":
				msg = fmt.Sprintf("The value '%v' is incorrect for address", err.Value())
				logger.Error("Address validation failed",
					zap.String("Tag", err.Tag()),
					zap.String("Type", err.Type().String()),
					zap.String("Issue", msg))
				errMsg = append(errMsg, msg)
			case "Latitude":
				msg = fmt.Sprintf("The value '%v' is incorrect for latitude", err.Value())
				logger.Error("Latitude validation failed",
					zap.String("Tag", err.Tag()),
					zap.String("Type", err.Type().String()),
					zap.String("Issue", msg))
				errMsg = append(errMsg, msg)
			case "Longitude":
				msg = fmt.Sprintf("The value '%v' is incorrect for longitude", err.Value())
				logger.Error("Longitude validation failed",
					zap.String("Tag", err.Tag()),
					zap.String("Type", err.Type().String()),
					zap.String("Issue", msg))
				errMsg = append(errMsg, msg)
			case "Price":
				msg = fmt.Sprintf("The value '%v' is incorrect for Price", err.Value())
				logger.Error("Price validation failed",
					zap.String("Tag", err.Tag()),
					zap.String("Type", err.Type().String()),
					zap.String("Issue", msg))
				errMsg = append(errMsg, msg)
			}
		}
		//LogError("Payload's required validation failed", err, "payload", fmt.Sprintf("%#v", err.(validator.FieldError).Field()))
		return errMsg, false
	}
	return nil, true
}

func CustomNameValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return AlphaRegexCheck(value)
}

func CustomAddressValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return AlphaNumRegexCheck(value)
}

// func CustomNumValidation(fl validator.FieldLevel) bool {
// 	value := fl.Field().String()
// 	return NumRegexCheck(value)
// }
