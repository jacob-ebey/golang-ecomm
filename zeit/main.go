package zeit

import (
	"fmt"
	"net/http"

	core "github.com/jacob-ebey/graphql-core"
	httphandler "github.com/jacob-ebey/graphql-httphandler"

	"github.com/jacob-ebey/golang-ecomm/auth"
	"github.com/jacob-ebey/golang-ecomm/runtime"
)

var handler *httphandler.GraphQLHttpHandler
var err error

func initialize() bool {
	if handler == nil && err == nil {
		executor, executorErr := runtime.NewExecutor(runtime.NewExecutorOpts{
			RunBefore: []core.PreExecuteHook{
				&auth.HttpHeaderHook{
					Source: "Authorization",
					Dest:   "authorization",
				},
			},
		})

		if executorErr != nil {
			runtime.PrintError(err)
			return false
		}

		handler = &httphandler.GraphQLHttpHandler{
			Executor:   *executor,
			Playground: true,
		}
	}

	return true
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if initialize() == false {
		fmt.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	handler.ServeHTTP(w, r)
}
