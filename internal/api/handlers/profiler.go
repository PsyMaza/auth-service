package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"gitlab.com/g6834/team17/auth-service/internal/interfaces"
	"gitlab.com/g6834/team17/auth-service/internal/models"
	"gitlab.com/g6834/team17/auth-service/internal/utils"
	"net/http"
	"net/http/pprof"
	"runtime"
	"strconv"
)

type profilerHandlers struct {
	logger     *zerolog.Logger
	presenters interfaces.Presenters
}

func newProfilerHandlers(logger *zerolog.Logger, presenters interfaces.Presenters) *profilerHandlers {
	return &profilerHandlers{
		logger:     logger,
		presenters: presenters,
	}
}

func ProfilerRouter(logger *zerolog.Logger, presenters interfaces.Presenters) http.Handler {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/pprof/", http.StatusMovedPermanently)
	})
	r.HandleFunc("/pprof", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/", http.StatusMovedPermanently)
	})

	handlers := newProfilerHandlers(logger, presenters)
	r.Post("/profile/status", handlers.profileStatus)

	r.HandleFunc("/pprof/*", pprof.Index)
	r.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/pprof/profile", pprof.Profile)
	r.HandleFunc("/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/pprof/trace", pprof.Trace)

	return r
}

func (h *profilerHandlers) profileStatus(w http.ResponseWriter, r *http.Request) {
	_, span := utils.StartSpan(r.Context())
	defer span.End()

	status := r.URL.Query().Get("status")
	if len(status) == 0 {
		h.presenters.Error(w, r, models.ErrorBadRequest(errors.New(`param "staus" not set`)))
		return
	}

	var rate int
	switch status {
	case "enable":
		rate = 1
	case "disable":
		rate = 0
	}

	profile := r.URL.Query().Get("profile")

	if len(profile) > 0 {
		switch profile {
		case "memory":
			queryRate := r.URL.Query().Get("rate")
			if len(queryRate) > 0 {
				parseRate, err := strconv.Atoi(queryRate)
				if err != nil {
					h.presenters.Error(w, r, models.ErrorBadRequest(err))
					return
				}
				rate = parseRate
			}
			setMemProfileRate(rate)
		case "block":
			setBlockProfileRate(rate)
		case "mutex":
			setMutexProfileRate(rate)
		}
	} else {
		setMemProfileRate(rate)
		setBlockProfileRate(rate)
		setMutexProfileRate(rate)
	}
}

func setMemProfileRate(rate int) {
	runtime.MemProfileRate = rate
}

func setBlockProfileRate(rate int) {
	runtime.SetBlockProfileRate(rate)
}

func setMutexProfileRate(rate int) {
	runtime.SetMutexProfileFraction(rate)
}
