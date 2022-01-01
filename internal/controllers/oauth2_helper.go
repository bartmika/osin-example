package controllers

import (
	"context"
	"errors"

	"github.com/bartmika/osin-example/internal/models"
	"github.com/bartmika/osin-example/internal/utils"
)

func (h *Controller) authenticatedUser(ctx context.Context, email string, password string) (*models.User, error) {
	// Lookup the user in our database, else return a `400 Bad Request` error.
	user, err := h.UserRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("Incorrect email or password")
	}

	// Verify the inputted password and hashed password match.
	passwordMatch := utils.CheckPasswordHash(password, user.PasswordHash)
	if passwordMatch == false {
		return nil, errors.New("Incorrect email or password")
	}
	return user, nil
}
