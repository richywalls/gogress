package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"portals"
)

func init() {
	r := mux.NewRouter()
	//r := s.HandleFunc("/*", portals.Test).Subrouter()
	r.Handle("/api/portal/", portals.AuthHandler(portals.GetPortals)).Methods("GET")
	r.Handle("/api/portal/{key}", portals.AuthHandler(portals.SavePortal)).Methods("POST")
	r.Handle("/api/portal/{key}", portals.AuthHandler(portals.GetPortal)).Methods("GET")

	r.Handle("/api/agent/{key}", portals.AuthHandler(portals.SaveAgent)).Methods("POST")
	r.Handle("/api/agent/", portals.AuthHandler(portals.GetAgents)).Methods("GET")

	r.Handle("/api/op/{key}", portals.AuthHandler(portals.SaveOperation)).Methods("POST")
	r.Handle("/api/op/", portals.AuthHandler(portals.GetOperations)).Methods("GET")

	r.HandleFunc("/auth/google", portals.Authenticate).Methods("POST")
	http.Handle("/", r)
}