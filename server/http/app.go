package http

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
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
) (*App, error) {
	app := &App{
		config: config,
		logger: logger,
	}

	err := app.configureRoutes()
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

func (a *App) configureRoutes() error {
	a.logger.Info("configuring http routes")
	var err error
	a.router, err = a.buildRoutes()
	return err
}

func (a *App) buildRoutes() (*mux.Router, error) {
	r := mux.NewRouter()

	authenticatedRouter := r.PathPrefix("/").Subrouter()
	authenticatedRouter.Handle("/healthcheck", newHealthcheckController()).Methods("GET")
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
