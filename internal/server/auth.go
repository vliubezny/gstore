package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vliubezny/gstore/internal/auth"
	"github.com/vliubezny/gstore/internal/model"
)

func (s *server) registerHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	var req registrationForm
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	u := model.User{
		Email: req.Email,
	}

	u, err := s.a.Register(r.Context(), u, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrEmailIsTaken) {
			writeError(l.WithError(err), w, http.StatusBadRequest, "email address has been already taken")
			return
		}

		writeInternalError(l.WithError(err), w, "fail to register user")
		return
	}

	writeOK(l, w, fromUserModel(u))
}
