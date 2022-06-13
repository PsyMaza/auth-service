package middlewares

import (
	"context"
	"gitlab.com/g6834/team17/auth-service/internal/constants"
	"gitlab.com/g6834/team17/auth-service/internal/interfaces"
	"gitlab.com/g6834/team17/auth-service/internal/models"
	"net/http"
)

func Validate(presenters interfaces.Presenters, authService interfaces.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			at, err := r.Cookie(constants.ACCESS_TOKEN)
			if err != nil {
				presenters.Error(rw, r, models.ErrorForbidden(err))
				return
			}

			rt, err := r.Cookie(constants.REFRESH_TOKEN)
			if err != nil {
				presenters.Error(rw, r, models.ErrorForbidden(err))
				return
			}

			ok, err := authService.VerifyToken(r.Context(), at.Value)
			if !ok || err != nil {
				presenters.Error(rw, r, models.ErrorForbidden(err))
				return
			}
			ok, err = authService.VerifyToken(r.Context(), rt.Value)
			if !ok || err != nil {
				presenters.Error(rw, r, models.ErrorForbidden(err))
				return
			}

			user, _, _ := authService.ParseToken(r.Context(), at.Value)

			ctx := context.WithValue(r.Context(), constants.CTX_USER, user)

			next.ServeHTTP(rw, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
