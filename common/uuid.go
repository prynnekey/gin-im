package common

import (
	uuid "github.com/satori/go.uuid"
)

func GenerateUUID() string {
	u := uuid.NewV4()
	return u.String()
}
