package http

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/guilhermeCoutinho/api-studies/server/http/controller"
	"github.com/guilhermeCoutinho/api-studies/usecase"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type App struct {
	address string
	config  *viper.Viper
	logger  logrus.FieldLogger
	router  *mux.Router
}

func NewApp(
	config *viper.Viper,
	logger logrus.FieldLogger,
	usecase *usecase.Usecase,
) (*App, error) {
	app := &App{
		config: config,
		logger: logger,
	}

	err := app.configureRoutes(usecase)
	if err != nil {
		return nil, err
	}
	app.configureAddress()

	return app, nil
}

func (a *App) configureAddress() {
	a.logger.Info("configuring http address")
	a.address = a.config.GetString("http.address")
}

func (a *App) configureRoutes(usecase *usecase.Usecase) error {
	a.logger.Info("configuring http routes")
	var err error
	a.router, err = a.buildRoutes(usecase)
	return err
}

func (a *App) buildRoutes(usecase *usecase.Usecase) (*mux.Router, error) {
	r := mux.NewRouter()

	healthCheckController := controller.NewHealthcheck(a.logger)
	authController := controller.NewAuth(a.logger)
	userController := controller.NewUser(a.logger, usecase)

	r.HandleFunc("/healthcheck", healthCheckController.HealthCheck).Methods("GET")

	r.HandleFunc("/users", userController.Create).Methods("POST")
	r.HandleFunc("/users", userController.Login).Methods("GET")

	r.HandleFunc("/auth", authController.Authenticate).Methods("POST")
	return r, nil
}

func (a *App) ListenAndServe() {
	a.logger.Infof("Starting listening on %s", a.address)
	err := http.ListenAndServe(a.address, a.router)
	if err != nil {
		a.logger.WithError(err).Error("Error on starting server")
		os.Exit(1)
	}
}
