package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hillview.tv/coreAPI/env"
	"github.com/hillview.tv/coreAPI/middleware"
	"github.com/hillview.tv/coreAPI/routers"
)

func main() {
	primary := mux.NewRouter()

	// Healthcheck Endpoint

	primary.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	// Define the API Endpoints

	r := primary.PathPrefix("/core/v1.1").Subrouter()

	// Logging of requests
	r.Use(middleware.LoggingMiddleware)

	// Adding response headers
	r.Use(middleware.MuxHeaderMiddleware)

	// Track & Update Last Active
	r.Use(middleware.TokenHandlers)

	// Core Admin

	admin := r.PathPrefix("/admin").Subrouter()

	// Admin Creators

	create := admin.PathPrefix("/create").Subrouter()

	create.Handle("/link", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.CreateLink))).Methods(http.MethodPost)

	// Admin Edits

	edit := admin.PathPrefix("/edit").Subrouter()

	edit.Handle("/asset", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleEditAsset))).Methods(http.MethodPost)
	edit.Handle("/video", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleEditVideo))).Methods(http.MethodPost)

	// Admin Deletes

	delete := admin.PathPrefix("/delete").Subrouter()

	delete.Handle("/video", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleDeleteVideo))).Methods(http.MethodPost)

	// Admin Lists

	list := admin.PathPrefix("/list").Subrouter()

	list.Handle("/adminUsers/{limit}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListAdminUsers))).Methods(http.MethodGet)
	list.Handle("/mobileUsers/{limit}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListMobileUsers))).Methods(http.MethodGet)

	list.Handle("/assets/{limit}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListAssets))).Methods(http.MethodGet)
	list.Handle("/checkouts", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListCheckouts))).Methods(http.MethodGet)

	list.Handle("/openCheckouts", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListOpenCheckouts))).Methods(http.MethodGet)

	list.Handle("/videos/{limit}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListVideos))).Methods(http.MethodGet)
	list.Handle("/links", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListLinks))).Methods(http.MethodGet)

	// Admin Searcher

	search := r.PathPrefix("/search").Subrouter()

	search.HandleFunc("/mobileUsers", routers.HandleSearchAdminUsers).Methods(http.MethodGet)

	// Launch API Listener
	fmt.Printf("✅ Hillview Core API running on port %s\n", env.Port)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Origin", "Authorization", "Accept", "X-CSRF-Token"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	log.Fatal(http.ListenAndServe(":"+env.Port, handlers.CORS(originsOk, headersOk, methodsOk)(primary)))
}
