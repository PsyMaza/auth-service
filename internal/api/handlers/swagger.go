package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "gitlab.com/g6834/team17/auth-service/api/swagger/public"
	"net/http"
)

func SwaggerRouter(basePath string) http.Handler {
	r := chi.NewRouter()

	httpSwagger.UIConfig(map[string]string{
		"showExtensions":        "true",
		"onComplete":            `() => { window.ui.setBasePath('v3'); }`,
		"defaultModelRendering": `"model"`,
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", basePath))))

	return r
}
