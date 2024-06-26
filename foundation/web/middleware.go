package web

type Middleware func(Handler) Handler

func wrapMiddleware(middlewares []Middleware, handler Handler) Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		middleware := middlewares[i]
		if middleware != nil {
			handler = middleware(handler)
		}
	}

	return handler
}
