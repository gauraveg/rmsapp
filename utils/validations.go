package utils

import (
	"regexp"
	"slices"

	"github.com/gauraveg/rmsapp/models"
)

func ValidateUserPayload(payload models.UserData) bool {
	roles := []string{"user", "admin", "sub-admin"}

	matchName, _ := regexp.MatchString("[A-Za-z ]", payload.Name)
	matchEmail, _ := regexp.MatchString("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$", payload.Email)
	matchRole := slices.Contains(roles, payload.Role)

	if matchEmail && matchName && matchRole {
		return true
	} else {
		return false
	}
}

func ValidateAddress(payload models.UserData) models.UserData {
	addresses := payload.Addresses
	for i := range addresses {
		matchAddr, _ := regexp.MatchString("[A-Za-z ]", addresses[i].Address)
		if !matchAddr {
			addresses[i] = addresses[len(addresses)-1]
			addresses = addresses[:len(addresses)-1]
		}
	}

	payload.Addresses = addresses
	return payload
}
