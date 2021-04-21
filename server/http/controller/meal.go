package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
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

type CreateMealRequest struct {
	Text     string    `json:"text"`
	Calories int       `json:"calories"`
	Time     time.Time `json:"time"`
}

func (u *Meal) Create(w http.ResponseWriter, r *http.Request) {
	args := &CreateMealRequest{}
	err := json.NewDecoder(r.Body).Decode(args)
	if err != nil {
		u.logger.WithError(err).Error("Failed to parse signup payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = u.usecase.CreateMeal(r.Context(), uuid.New(), "", -1)
	if err != nil {
		u.logger.WithError(err).Error("Failed to create user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	u.logger.Info("User created successfully")
}
