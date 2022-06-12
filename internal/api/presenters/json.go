package presenters

import (
	"gitlab.com/g6834/team17/auth-service/internal/utils"
	"net/http"
)

func (p *presenters) JSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	err := utils.WriteJson(w, v)
	if err != nil {
		p.Error(w, r, err)
	}
}
