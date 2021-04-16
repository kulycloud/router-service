package routing

import (
	"context"
	"fmt"
	commonHttp "github.com/kulycloud/common/http"
)

func forwardRoute(ctx context.Context, req *commonHttp.Request, routerResult *RouterResult) *commonHttp.Response {
	ref, ok := req.KulyData.Step.References[routerResult.Destination]
	if !ok || ref == nil {
		return buildErrorResponse(fmt.Errorf("%w: reference %s not found", ErrInvalidConfig, routerResult.Destination))
	}

	if routerResult.RewrittenPath != nil {
		req.Path = *routerResult.RewrittenPath
	}

	req.KulyData.StepUid = ref.Step
	res, err := commonHttp.ProcessRequest(ctx, ref.Endpoints, req)
	if err != nil {
		return buildErrorResponse(fmt.Errorf("error during forwarding: %w", err))
	}

	return res
}

