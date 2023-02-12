package utils

import uuid "github.com/nu7hatch/gouuid"

func GenerateUUIDV4() *uuid.UUID {
	id, _ := uuid.NewV4()
	return id
}
