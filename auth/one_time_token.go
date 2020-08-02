package auth

import (
	"time"

	"github.com/odmishien/go-auth-playground/models"
	uuid "github.com/satori/go.uuid"
)

func GetNewOneTimeToken(email string, expiresAt time.Time) (models.OneTimeScript, error) {
	ots := models.OneTimeScript{
		Token:  uuid.NewV4().String(),
		Expire: expiresAt,
		Email:  email,
	}
	return ots, nil
}
