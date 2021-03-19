package routing

import (
	"errors"
	"fmt"
	commonHttp "github.com/kulycloud/common/http"
)

type RouterRequest struct {
	Path string
	Method string
}

type RouterResult struct {
	RewrittenPath *string
	Destination string
}

var ErrNoRoute = errors.New("no matching route found")
var ErrInvalidConfig = errors.New("invalid config")

func RouterRequestFromRequest(req *commonHttp.Request) *RouterRequest {
	return &RouterRequest{
		Path:   req.Path,
		Method: req.Method,
	}
}

func (req *RouterRequest) Clone() *RouterRequest {
	return &RouterRequest{
		Path:   req.Path,
		Method: req.Method,
	}
}

func (routeList *RouteList) Route(req *RouterRequest) (*RouterResult, error) {
	for i := range *routeList {
		route := &((*routeList)[i])
		res, err := route.tryRoute(req)
		if err == nil {
			return res, nil // We got our route!
		}
		if !errors.Is(err, ErrNoRoute) {
			return nil, err // Some other error occurred! Want to report that!
		}
	}

	return nil,ErrNoRoute
}

// Try to route to a specific Route (return ErrNoRoute if impossible)
func (route *Route) tryRoute(req *RouterRequest) (*RouterResult, error) {
	rewrittenPath, err := route.tryMatch(req)
	if err != nil {
		return nil, err
	}

	return route.getDestination(req, rewrittenPath)
}

func (route *Route) getDestination(req *RouterRequest, rewrittenPath *string) (*RouterResult, error) {
	if route.Destination != nil && route.Subroutes != nil {
		return nil, fmt.Errorf("%w: can only specify either path.destination or path.subroutes", ErrInvalidConfig)
	}

	if route.Destination != nil {
		return &RouterResult{
			RewrittenPath: rewrittenPath,
			Destination: *route.Destination,
		}, nil
	}

	if route.Subroutes != nil {
		forwardedReq := req
		if forwardedReq.Path != *rewrittenPath { // Path was rewritten! Create new request
			forwardedReq = forwardedReq.Clone()
			forwardedReq.Path = *rewrittenPath
		}

		return route.Subroutes.Route(forwardedReq)
	}

	return nil, fmt.Errorf("%w: have to specify either path.destination or path.subroutes", ErrInvalidConfig)
}

func (route *Route) tryMatch(req *RouterRequest) (*string, error) {
	rewrittenPath := &req.Path

	if route.Path != nil {
		var err error
		var match bool
		match, rewrittenPath, err = route.Path.matches(req)
		if err != nil {
			return nil, err
		}

		if !match {
			return nil, ErrNoRoute
		}
	}

	if route.Methods != nil {
		match := false
		for _, method := range *route.Methods {
			if req.Method == method {
				match = true
				break
			}
		}
		if !match {
			return nil, ErrNoRoute
		}
	}

	return rewrittenPath, nil
}

func (path *PathSpecification) matches(req *RouterRequest) (bool, *string, error) { // matches, rewrittenPath, error
	if path.Exact != nil && path.Prefix != nil {
		return false, nil, fmt.Errorf("%w: can only specify either path.exact or path.prefix", ErrInvalidConfig)
	}

	if path.Exact != nil {
		matches := req.Path == *path.Exact
		if matches {
			if path.Rewrite != nil && *path.Rewrite {
				newPathData := "/"
				return matches, &newPathData, nil
			} else {
				return matches, &req.Path, nil
			}
		}
		return false, nil, nil
	}

	if path.Prefix != nil {
		prefixLen := len(*path.Prefix)
		matches := req.Path[:prefixLen] == *path.Prefix

		if matches {
			if path.Rewrite != nil && *path.Rewrite {
				newPathData := req.Path[prefixLen:]
				// Add leading "/" if missing it in rewritten path
				if newPathData[0] != '/' {
					newPathData = "/" + newPathData
				}
				return matches, &newPathData, nil
			} else {
				return matches, &req.Path, nil
			}
		}
		return false, nil, nil
	}

	return false, nil, fmt.Errorf("%w: have to specify either path.exact or path.prefix", ErrInvalidConfig)
}
