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

func CheckValidation(payload interface{}) bool {
	validate := validator.New()
	err := validate.RegisterValidation("UserNameCheck", CustomNameValidation)
	err = validate.RegisterValidation("AddressCheck", CustomAddressValidation)
	//err = validate.RegisterValidation("NumberCheck", CustomNumValidation)

	err = validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Name":
				zap.L().Error("Name validation failed",
					zap.String("Tag", err.Tag()),
					zap.String("Type", err.Type().String()),
					zap.String("Issue", fmt.Sprintf("The value %v is incorrect for name", err.Value())))
			case "Email":
				zap.L().Error("Email validation failed",
					zap.String("Tag", err.Tag()),
					zap.String("Type", err.Type().String()),
					zap.String("Issue", fmt.Sprintf("The value %v is incorrect for email", err.Value())))
			case "Role":
				zap.L().Error("Role validation failed",
					zap.String("Tag", err.Tag()),
					zap.String("Type", err.Type().String()),
					zap.String("Issue", fmt.Sprintf("The value %v is incorrect. It should be either of these values %v", err.Value(), err.Param())))
			case "Password":
				zap.L().Error("Password validation failed",
					zap.String("Tag", err.Tag()),
					zap.String("Type", err.Type().String()),
					zap.String("Issue", fmt.Sprintf("The length for password is incorrect. The length is %v", len(err.Value().(string)))))
			case "Address":
				zap.L().Error("Address validation failed",
					zap.String("Tag", err.Tag()),
					zap.String("Type", err.Type().String()),
					zap.String("Issue", fmt.Sprintf("The value %v is incorrect for address", err.Value())))
			case "Latitude":
				zap.L().Error("Latitude validation failed",
					zap.String("Tag", err.Tag()),
					zap.String("Type", err.Type().String()),
					zap.String("Issue", fmt.Sprintf("The value %v is incorrect for latitude", err.Value())))
			case "Longitude":
				zap.L().Error("Longitude validation failed",
					zap.String("Tag", err.Tag()),
					zap.String("Type", err.Type().String()),
					zap.String("Issue", fmt.Sprintf("The value %v is incorrect for longitude", err.Value())))
			case "Price":
				zap.L().Error("Price validation failed",
					zap.String("Tag", err.Tag()),
					zap.String("Type", err.Type().String()),
					zap.String("Issue", fmt.Sprintf("The value %v is incorrect for Price", err.Value())))
			}
		}
		//LogError("Payload's required validation failed", err, "payload", fmt.Sprintf("%#v", err.(validator.FieldError).Field()))
		return false
	}
	return true
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
