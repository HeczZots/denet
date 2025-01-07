package server

import (
	"denet/db"
	"denet/server/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	h    *mux.Router
	addr string
	db   *db.DB
	mw   *middleware.AuthMiddleware
}

func New(addr string, db *db.DB, jwtSecret []byte) *Server {
	s := &Server{addr: addr, db: db, mw: middleware.NewAuthMiddleware(jwtSecret)}
	r := mux.NewRouter()
	s.h = r
	// Middleware

	r.Use(s.mw.Auth)

	// Routes
	r.HandleFunc("/users/{id}/status", s.getUserStatus).Methods("GET")
	r.HandleFunc("/users/leaderboard", s.getLeaderboard).Methods("GET")
	r.HandleFunc("/users/{id}/task/complete", s.completeTask).Methods("POST")
	r.HandleFunc("/users/{id}/referrer", s.setReferrer).Methods("POST")
	r.HandleFunc("/login", s.login).Methods("POST")
	r.HandleFunc("/registration", s.registration).Methods("POST")

	return s
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.addr, s.h)
}
