package main

import (
	"net/http"
	"path"
	"time"

	"github.com/gorilla/mux"
	core "github.com/jacob-ebey/graphql-core"
	httphandler "github.com/jacob-ebey/graphql-httphandler"
	"github.com/joho/godotenv"

	"github.com/jacob-ebey/golang-ecomm/auth"
	"github.com/jacob-ebey/golang-ecomm/runtime"
)

func main() {
	godotenv.Load(".env")

	executor, err := runtime.NewExecutor(runtime.NewExecutorOpts{
		RunBefore: []core.PreExecuteHook{
			&auth.HttpHeaderHook{
				Source: "Authorization",
				Dest:   "authorization",
			},
		},
	})

	if err != nil {
		runtime.PrintError(err)
		panic(err)
	}

	handler := httphandler.GraphQLHttpHandler{
		Executor:   *executor,
		Playground: true,
	}

	router := mux.NewRouter()
	router.HandleFunc("/graphql", handler.ServeHTTP)

	if runtime.ShouldServeStaticFiles() {
		fileServer := http.FileServer(http.Dir(path.Clean("./frontend/build")))
		router.PathPrefix("/images").Handler(fileServer)
		router.PathPrefix("/static").Handler(fileServer)
		router.Handle("/favicon.ico", fileServer)
		router.Handle("/manifest.json", fileServer)

		router.PathPrefix("/").HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, path.Clean("./frontend/build/index.html"))
		}))
	}

	server := &http.Server{
		Handler:      router,
		Addr:         runtime.GetAddress(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
