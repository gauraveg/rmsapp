package utils

import (
	"regexp"
	"slices"
	"strconv"

	"github.com/gauraveg/rmsapp/models"
)

func AlphaNumRegexCheck(value string) bool {
	isValid, _ := regexp.MatchString("^[A-Za-z0-9 ]+$", value)
	return isValid
}

func NumRegexCheck(value string) bool {
	isValid, _ := regexp.MatchString("^[0-9.]+$", value)
	return isValid
}

func ValidateUserPayload(payload models.UserData) bool {
	roles := []string{"user", "admin", "sub-admin"}

	matchName := AlphaNumRegexCheck(payload.Name)
	matchEmail, _ := regexp.MatchString("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$", payload.Email)
	matchRole := slices.Contains(roles, payload.Role)

	if matchEmail && matchName && matchRole {
		return true
	} else {
		return false
	}
}

func ValidateUserAddress(payload models.UserData) models.UserData {
	addresses := payload.Addresses
	regEx, _ := regexp.Compile("^[A-Za-z0-9 ]+$")
	newAddr := make([]models.AddressData, 0)
	for i := range addresses {
		matchAddr := regEx.MatchString(addresses[i].Address)
		if matchAddr {
			newAddr = append(newAddr, addresses[i])
		}
	}

	payload.Addresses = newAddr
	return payload
}

func ValidateRestPayload(payload models.RestaurantsRequest) bool {
	matchName := AlphaNumRegexCheck(payload.Name)

	if matchName {
		return true
	} else {
		return false
	}
}

func ValidateRestAddress(payload models.RestaurantsRequest) models.RestaurantsRequest {
	matchAddr := AlphaNumRegexCheck(payload.Address)
	matchLat := NumRegexCheck(strconv.FormatFloat(payload.Latitude, 'f', -1, 64))
	matchLong := NumRegexCheck(strconv.FormatFloat(payload.Longitude, 'f', -1, 64))
	if !matchAddr {
		payload.Address = ""
	}
	if !matchLat {
		payload.Latitude = 0
	}
	if !matchLong {
		payload.Longitude = 0
	}
	return payload
}

func ValidateDishPayload(payload models.DishRequest) bool {
	matchName := AlphaNumRegexCheck(payload.Name)
	matchPrice := NumRegexCheck(strconv.Itoa(payload.Price))
	if matchName && matchPrice {
		return true
	} else {
		return false
	}
}
