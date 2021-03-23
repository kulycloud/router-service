package routing

import (
	"context"
	commonHttp "github.com/kulycloud/common/http"
	"log"
	"net/http"
)

func HandleRequest(ctx context.Context, req *commonHttp.Request) *commonHttp.Response {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recovered panic: %e", err)
		}
	}()
	conf, err := ConfigFromRequest(req)
	if err != nil {
		return buildErrorResponse(err)
	}

	routerRequest := RouterRequestFromRequest(req)

	routerResult, err := conf.Routes.Route(routerRequest)
	if err != nil {
		// TODO extra 404 handler
		return buildErrorResponse(err)
	}

	return forwardRoute(ctx, req, routerResult)
}

func buildErrorResponse(err error) *commonHttp.Response {
	resp := commonHttp.NewResponse()
	if err == ErrNoRoute {
		resp.Status = http.StatusNotFound
	} else {
		resp.Status = http.StatusInternalServerError
	}

	resp.Body.Write([]byte(err.Error()))
	resp.Headers.Set("Content-Type", "text/plain")

	return resp
}
