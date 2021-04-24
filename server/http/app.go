package http

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/guilhermeCoutinho/api-studies/dal"
	"github.com/guilhermeCoutinho/api-studies/server/http/controller"
	"github.com/guilhermeCoutinho/api-studies/server/http/wrapper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type App struct {
	address string
	config  *viper.Viper
	logger  logrus.FieldLogger
	router  *mux.Router
	wrapper *wrapper.HTTPWrapper
}

func NewApp(
	config *viper.Viper,
	logger logrus.FieldLogger,
	dal *dal.DAL,
) (*App, error) {
	app := &App{
		config:  config,
		logger:  logger,
		wrapper: wrapper.NewHTTPWrapper(logger),
	}

	app.buildRoutes(dal)
	app.configureAddress()

	return app, nil
}

func (a *App) configureAddress() {
	a.logger.Info("configuring http address")
	a.address = a.config.GetString("http.address")
}

func (a *App) buildRoutes(dal *dal.DAL) {
	authMiddleware := Middleware{
		logger: a.logger,
	}

	router := mux.NewRouter()
	authRouter := router.PathPrefix("/").Subrouter()
	authRouter.Use(authMiddleware.Authenticate)

	healthCheckController := controller.NewHealthcheck()
	userController := controller.NewUser(dal, a.config)
	authController := controller.NewAuth(dal, a.config)
	mealController := controller.NewMeal(dal, a.config)

	a.wrapper.Register(router, "/users", userController)
	a.wrapper.Register(router, "/auth", authController)
	a.wrapper.Register(router, "/healthcheck", healthCheckController)

	a.wrapper.Register(authRouter, "/users/{userID}/meals", mealController)

	a.router = router
}

func (a *App) ListenAndServe() {
	a.logger.Infof("Starting listening on %s", a.address)
	err := http.ListenAndServe(a.address, a.router)
	if err != nil {
		a.logger.WithError(err).Error("Error on starting server")
		os.Exit(1)
	}
}
