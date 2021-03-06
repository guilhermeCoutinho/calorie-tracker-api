package http

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/guilhermeCoutinho/api-studies/controller/auth"
	"github.com/guilhermeCoutinho/api-studies/controller/healthcheck"
	"github.com/guilhermeCoutinho/api-studies/controller/meal"
	"github.com/guilhermeCoutinho/api-studies/controller/noauthuser"
	"github.com/guilhermeCoutinho/api-studies/controller/user"
	"github.com/guilhermeCoutinho/api-studies/dal"
	"github.com/guilhermeCoutinho/api-studies/server/http/wrapper"
	"github.com/guilhermeCoutinho/api-studies/services/calorieprovider"
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

	healthCheckController := healthcheck.NewHealthcheck()
	userNoAuthController := noauthuser.NewUserNoAuth(dal, a.config)
	userController := user.NewUser(dal, a.config)
	authController := auth.NewAuth(dal, a.config)
	mealController := meal.NewMeal(dal, a.config, &calorieprovider.ProviderImpl{})

	a.wrapper.Register(router, "/users", userNoAuthController)
	a.wrapper.Register(router, "/auth", authController)
	a.wrapper.Register(router, "/healthcheck", healthCheckController)

	a.wrapper.Register(authRouter, "/users/{userID}/meals", mealController)
	a.wrapper.Register(authRouter, "/users/{userID}/meals/{mealID}", mealController)
	a.wrapper.Register(authRouter, "/users/{userID}", userController)

	a.wrapper.Register(authRouter, "/meals/{mealID}", mealController)
	a.wrapper.Register(authRouter, "/meals", mealController)
	a.wrapper.Register(authRouter, "/users", userController)

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
