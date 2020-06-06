package apps

import (
	"github.com/elastic/go-elasticsearch"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"log"
	"ntm-backend/api/middleware"
	"ntm-backend/builder"
	"ntm-backend/repository"
	"ntm-backend/service"
)

type ApiContainer struct {
	Logger *logrus.Logger
}

func (a ApiContainer) Register(router *mux.Router, path string) {
	api := mux.NewRouter()

	// Api Middleware
	auth := middleware.ApiAuthMiddleware{a.Logger}

	esRepository := repository.ESRepository{a.Logger}

	evolutionService := service.TemporalEvolutionService{a.Logger, esRepository}


	snmpRepository := repository.SnmpESRepository{a.Logger, newESClient()}

	builderRepo := builder.SnmpBuilder{a.Logger}

	snmpService := service.SnmpService{a.Logger, snmpRepository, builderRepo}

	// Routers
	apiAuthRouter := ApiRouter{a.Logger, evolutionService, snmpService}

	apiAuthRouter.RegisterApiRoutes(api, path)

	router.PathPrefix(path).Handler(negroni.New(
		negroni.HandlerFunc(auth.ApiAuth),
		negroni.Wrap(api),
	))
}

func newESClient() *elasticsearch.Client  {
	cfg := elasticsearch.Config{}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating new elastic search client : %s", err)
		return &elasticsearch.Client{}
	}
	return es
}

