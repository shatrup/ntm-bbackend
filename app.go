package main

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
	"ntm-backend/api/apps"
	"os"
)

type App struct {
	ExternalRouter *mux.Router
}

func (a *App) Initialize() {
	a.InitializeExternalRoutes()
}

func (a *App) InitializeExternalRoutes() {
	a.ExternalRouter = mux.NewRouter()
	a.ExternalRouter.Use(commonMiddleware)

	registerApiAuthRoutes(a.ExternalRouter)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func registerApiAuthRoutes(router *mux.Router) {
	apiContainer := &apps.ApiContainer{Logger: logger}
	apiContainer.Register(router, "/api")
}

func (a *App) RunExternal() {
	externalServer := &http.Server{
		Addr:    os.Getenv("NTM_SERVER_IP") + ":" + os.Getenv("NTM_SERVER_PORT"),
		Handler: middlewareHandler(a.ExternalRouter),
	}
	logger.Printf("Starting external server at IP: %v and PORT: %v ", os.Getenv("NTM_SERVER_IP"), os.Getenv("NTM_SERVER_PORT"))
	logger.Fatal(externalServer.ListenAndServe())
}

func middlewareHandler(router *mux.Router) *negroni.Negroni {
	/* Keeping the commented code, moving from classic implementation for self defined so that we can pass logrus logger and
	the request/response logging will be done in the log file. Keeping the middleware logging commented for debugging purpose */
	//n := negroni.Classic()
	intLogger := negroni.NewLogger()

	var loggerDefaultFormat = "{{.Status}} | {{.Duration}} | {{.Hostname}} | {{.Method}} {{.Path}}"

	intLogger.ALogger = logger
	intLogger.SetFormat(loggerDefaultFormat)
	n := negroni.New(negroni.NewRecovery(), intLogger, negroni.NewStatic(http.Dir("public")))
	//n.Use(negronilogrus.NewMiddlewareFromLogger(logger, "web"))
	n.UseHandler(router)
	return n
}
