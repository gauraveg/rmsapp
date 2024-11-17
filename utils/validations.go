package utils

import (
	"regexp"
	"slices"

	"github.com/gauraveg/rmsapp/models"
)

func ValidateUserPayload(payload models.UserData) bool {
	roles := []string{"user", "admin", "sub-admin"}

	matchName, _ := regexp.MatchString("^[A-Za-z0-9 ]+$", payload.Name)
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
