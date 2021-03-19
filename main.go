package main

import (
	commonHttp "github.com/kulycloud/common/http"
	"github.com/kulycloud/common/logging"
	"github.com/kulycloud/router-service/routing"
)

var logger = logging.GetForComponent("service")

func main() {
	srv, err := commonHttp.NewServer(30000, routing.HandleRequest)
	if err != nil {
		logger.Panicw("could not create server", "error", err)
	}

	err = srv.Serve()

	if err != nil {
		logger.Panicw("could not serve", "error", err)
	}
}

