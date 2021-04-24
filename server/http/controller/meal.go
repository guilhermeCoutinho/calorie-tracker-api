package controller

import (
	"encoding/json"
	"net/http"

	"github.com/guilhermeCoutinho/api-studies/messages"
	"github.com/guilhermeCoutinho/api-studies/usecase"
	"github.com/sirupsen/logrus"
)

type Meal struct {
	logger  logrus.FieldLogger
	usecase *usecase.Usecase
}

func NewMeal(
	logger logrus.FieldLogger,
	usecase *usecase.Usecase,
) *Meal {
	return &Meal{
		logger:  logger,
		usecase: usecase,
	}
}

func (u *Meal) Create(w http.ResponseWriter, r *http.Request) {
	logger := u.logger.WithFields(logrus.Fields{
		"methodName": "createMeal",
	})

	args := &messages.CreateMealRequest{}
	err := json.NewDecoder(r.Body).Decode(args)
	if err != nil {
		logger.WithError(err).Error("Failed to parse payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = u.usecase.CreateMeal(r.Context(), args)
	if err != nil {
		logger.WithError(err).Error("Failed to create meal")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.Info("Meal created successfully")
}
