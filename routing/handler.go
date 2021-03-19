package routing

import (
	"context"
	commonHttp "github.com/kulycloud/common/http"
	"net/http"
)

func HandleRequest(_ context.Context, req *commonHttp.Request) *commonHttp.Response {
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

	return forwardRoute(req, routerResult)
}

func buildErrorResponse(err error) *commonHttp.Response {
	resp := commonHttp.NewResponse()
	resp.Status = http.StatusInternalServerError
	resp.Body.Write([]byte(err.Error()))
	resp.Headers.Set("Content-Type", "text/plain")

	return resp
}
