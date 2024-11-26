package utils

import (
	"context"
	"fmt"
	"regexp"

	"github.com/gauraveg/rmsapp/logger"
	"github.com/go-playground/validator/v10"
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

func CheckValidation(ctx context.Context, payload interface{}, loggers *logger.ZapLogger) ([]string, bool) {
	validate := validator.New()
	err := validate.RegisterValidation("UserNameCheck", CustomNameValidation)
	if err != nil {
		loggers.ErrorWithContext(ctx, err.Error())
	}
	err = validate.RegisterValidation("AddressCheck", CustomAddressValidation)
	if err != nil {
		loggers.ErrorWithContext(ctx, err.Error())
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
				loggers.ErrorWithContext(ctx, map[string]string{"message": "Name validation failed", "tag": err.Tag(), "type": err.Type().String(), "issue": msg})
				errMsg = append(errMsg, msg)
			case "Email":
				msg = fmt.Sprintf("The value '%v' is incorrect for email", err.Value())
				loggers.ErrorWithContext(ctx, map[string]string{"message": "Email validation failed", "tag": err.Tag(), "type": err.Type().String(), "issue": msg})
				errMsg = append(errMsg, msg)
			case "Role":
				msg = fmt.Sprintf("The value '%v' is incorrect. It should be either of these values %v", err.Value(), err.Param())
				loggers.ErrorWithContext(ctx, map[string]string{"message": "Role validation failed", "tag": err.Tag(), "type": err.Type().String(), "issue": msg})
				errMsg = append(errMsg, msg)
			case "Password":
				msg = fmt.Sprintf("The length for password is incorrect. The length is %v", len(err.Value().(string)))
				loggers.ErrorWithContext(ctx, map[string]string{"message": "Password validation failed", "tag": err.Tag(), "type": err.Type().String(), "issue": msg})
				errMsg = append(errMsg, msg)
			case "Address":
				msg = fmt.Sprintf("The value '%v' is incorrect for address", err.Value())
				loggers.ErrorWithContext(ctx, map[string]string{"message": "Address validation failed", "tag": err.Tag(), "type": err.Type().String(), "issue": msg})
				errMsg = append(errMsg, msg)
			case "Latitude":
				msg = fmt.Sprintf("The value '%v' is incorrect for latitude", err.Value())
				loggers.ErrorWithContext(ctx, map[string]string{"message": "Latitude validation failed", "tag": err.Tag(), "type": err.Type().String(), "issue": msg})
				errMsg = append(errMsg, msg)
			case "Longitude":
				msg = fmt.Sprintf("The value '%v' is incorrect for longitude", err.Value())
				loggers.ErrorWithContext(ctx, map[string]string{"message": "Longitude validation failed", "tag": err.Tag(), "type": err.Type().String(), "issue": msg})
				errMsg = append(errMsg, msg)
			case "Price":
				msg = fmt.Sprintf("The value '%v' is incorrect for Price", err.Value())
				loggers.ErrorWithContext(ctx, map[string]string{"message": "Price validation failed", "tag": err.Tag(), "type": err.Type().String(), "issue": msg})
				errMsg = append(errMsg, msg)
			}
		}
		//LogError("Payload's required validation failed", err, "payload", fmt.Sprintf("%#v", err.(validator.FieldError).Field()))
		return errMsg, false
	}
	return nil, true
}
