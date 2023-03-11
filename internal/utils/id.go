package utils

import gonanoid "github.com/matoous/go-nanoid/v2"

const (
	RANDOM_SUBDOMAIN_LENGTH = 6
	CONNECTION_ID_LENGTH    = 10
	SECRET_KEY_LENGTH       = 16
	SESSION_ID_LENGTH       = 8

	NANOID_ALPHABETS = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func generateNanoid(size int) string {
	id, _ := gonanoid.Generate(NANOID_ALPHABETS, size)
	return id
}

func GenerateRandomSubdomain() string {
	return generateNanoid(RANDOM_SUBDOMAIN_LENGTH)
}

func GenerateConnectionId() string {
	return generateNanoid(CONNECTION_ID_LENGTH)
}

func GenerateSecretKey() string {
	return generateNanoid(SECRET_KEY_LENGTH)
}

func GenerateSessionToken() string {
	return generateNanoid(SESSION_ID_LENGTH)
}
