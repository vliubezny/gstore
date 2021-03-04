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

	var req credentials
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

func (s *server) loginHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	var req credentials
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	l = l.WithField("email", req.Email)

	tokens, err := s.a.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			writeError(l.WithError(err), w, http.StatusUnauthorized, "invalid username or password")
			return
		}

		writeInternalError(l.WithError(err), w, "fail to login user")
		return
	}

	l.Info("logged in successfully")

	writeOK(l, w, fromTokenPairModel(tokens))
}

func (s *server) refreshHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	token := extractBearer(r)
	if token == "" {
		writeError(l, w, http.StatusUnauthorized, "missing token")
		return
	}

	tokens, err := s.a.Refresh(r.Context(), token)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) {
			writeError(l.WithError(err), w, http.StatusUnauthorized, "invalid refresh token")
			return
		}

		writeInternalError(l.WithError(err), w, "fail to refresh tokens")
		return
	}

	l.Info("refreshed tokens successfully")

	writeOK(l, w, fromTokenPairModel(tokens))
}

func (s *server) revokeHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	token := extractBearer(r)
	if token == "" {
		writeError(l, w, http.StatusUnauthorized, "missing token")
		return
	}

	err := s.a.Revoke(r.Context(), token)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) {
			writeError(l.WithError(err), w, http.StatusUnauthorized, "invalid refresh token")
			return
		}

		writeInternalError(l.WithError(err), w, "fail to revoke token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) updateUserPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	id, err := getIDFromURL(r, "id")
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid user ID")
		return
	}

	var req userPermissions
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	user := req.toModel()
	user.ID = id

	if err := s.a.UpdateUserPermissions(r.Context(), user); err != nil {
		switch {
		case errors.Is(err, auth.ErrNotFound):
			writeError(l.WithError(err), w, http.StatusNotFound, "user not found")
		default:
			writeInternalError(l.WithError(err), w, "fail to update user permissions")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
