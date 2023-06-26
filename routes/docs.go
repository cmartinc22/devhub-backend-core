//go:build !test
// +build !test

package routes

import (
	"net/http"
	"os"

	"github.com/cmartinc22/devhub-backend-core/handlers"
	"github.com/cmartinc22/devhub-backend-core/models"
	"github.com/pedidosya/peya-go/logs"
	"github.com/pedidosya/peya-go/server"
)

// Associate route to expose swagger apidoc on non-productive endpoints
func AddApiDocRoutes(cfg *models.DocsConfiguration, s *server.Server) {
	if cfg.Enabled {
		// Expose swagger compiled HTML view
		s.AddRouteWithOptions(
			"/doc",
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html")
				http.ServeFile(w, r, "api"+string(os.PathSeparator)+"api.html")
			},
			&server.RouteOptions{
				CORSEnabled: true,
				CORSOptions: []server.CORSOption{
					server.AllowedMethods([]string{http.MethodGet}),
				},
				TimeOutSeconds: 0,
			},
			http.MethodGet)
		logs.Info("[main] documentation enabled")
	} else {
		logs.Info("[main] documentation disabled")
	}
}

func AddOASRoutes(cfg *models.DocsConfiguration, s *server.Server) {
	// Expose OAS definition
	s.AddRouteWithOptions(
		"/api/{apiVersion}/oas.yaml",
		handlers.HandleGetOAS("api"),
		&server.RouteOptions{
			CORSEnabled: true,
			CORSOptions: []server.CORSOption{
				server.AllowedMethods([]string{http.MethodGet, http.MethodOptions}),
			},
			TimeOutSeconds: 0,
		},
		http.MethodGet,
		http.MethodOptions,
	)
}
