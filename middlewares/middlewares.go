package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/deliveryhero/devhub-backend-core/clients"
	"github.com/deliveryhero/devhub-backend-core/models"
	"github.com/pedidosya/peya-go/logs"
	"github.com/pedidosya/peya-go/server"
)

type (
	ServerHandler func(e interface{}, w http.ResponseWriter, r *http.Request)
	Middleware    func(interface{}, ServerHandler) ServerHandler
)

func Handler(engineSpec interface{}, handler ServerHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(engineSpec, w, r)
	}
}

func AuthCheckMiddleware(stsClient clients.STSClient, cfClient clients.CFClient, scope string) Middleware {
	return func(engineSpec interface{}, handler ServerHandler) ServerHandler {
		return func(engineSpec interface{}, w http.ResponseWriter, r *http.Request) {
			var authResult *models.AuthResult
			var err, cfErr, stsErr *models.CustomError

			authResult, stsErr = clients.STSAuthCheckMiddleware(stsClient, r, scope)
			stsEnabled := authResult.Enabled
			stsValid := authResult.IsValid

			if !stsEnabled || !stsValid {
				authResult, cfErr = clients.CFAuthCheckMiddleware(cfClient, r, scope)
				if !authResult.Enabled {
					err = stsErr
				} else if !authResult.IsValid {
					if cfErr == nil {
						err = &models.CustomError{
							Code:     models.NotFound,
							Messages: []string{fmt.Sprintf("Resource not found: %s", r.URL.Path)},
						}
					} else {
						err = cfErr
					}
				} else if authResult.IsValid {
					err = nil
				}
			}

			if err != nil {
				logs.Errorf("[middleware] auth: %s", strings.Join(err.Messages, ","))
				switch err.Code {
				case models.NotFound:
					server.NotFound(w, r, err.Messages...)
					return
				case models.Forbidden:
					server.Forbidden(w, r, err.Messages...)
					return
				}
			}

			ctx_id_value := ""
			if authResult.Identity != nil {
				ctx_id_value = *authResult.Identity
			}
			ctx := context.WithValue(r.Context(), models.IdentityContext, ctx_id_value)

			handler(engineSpec, w, r.WithContext(ctx))
		}
	}
}

func ApiVersionCheckMiddleware() Middleware {
	return func(engineSpec interface{}, handler ServerHandler) ServerHandler {
		return func(engineSpec interface{}, w http.ResponseWriter, r *http.Request) {
			apiVersion := server.GetStringFromPath(r, "apiVersion", "")

			// TODO get valid api version from DB or env
			if models.ToApiVersion(apiVersion) == models.NotSupported {
				server.Render(w, r, "Unsupported API version", 404)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), models.ApiVersionContext, apiVersion))
			handler(engineSpec, w, r)
		}
	}
}

func ContentType(contentType string) Middleware {
	return func(engineSpec interface{}, h ServerHandler) ServerHandler {
		return func(engineSpec interface{}, w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", contentType)
			h(engineSpec, w, r)
		}
	}
}

func ApplyMiddlewares(e interface{}, h ServerHandler, middlewares ...Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		h = middleware(e, h)
	}
	return Handler(e, h)
}
