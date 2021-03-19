package routing

import (
	"fmt"
	commonHttp "github.com/kulycloud/common/http"
)

func forwardRoute(req *commonHttp.Request, routerResult *RouterResult) *commonHttp.Response {
	ref, ok := req.KulyData.Step.References[routerResult.Destination]
	if !ok || ref == nil {
		return buildErrorResponse(fmt.Errorf("%w: reference %s not found", ErrInvalidConfig, routerResult.Destination))
	}

	// TODO Forward
	//comm := commonHttp.NewCommunicator(context.Background(), ref.Endpoints)
	//comm.Ping()

	return nil
}

