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

	edit.Handle("/video", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleEditVideo))).Methods(http.MethodPost)
	edit.Handle("/adminUser", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleEditAdminAccount))).Methods(http.MethodPost)

	// Admin Deletes

	delete := admin.PathPrefix("/delete").Subrouter()

	delete.Handle("/video", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleDeleteVideo))).Methods(http.MethodPost)

	// Admin Lists

	list := admin.PathPrefix("/list").Subrouter()

	list.Handle("/adminUsers/{limit}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListAdminUsers))).Methods(http.MethodGet)
	list.Handle("/mobileUsers", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListMobileUsers))).Methods(http.MethodGet)

	list.Handle("/openCheckouts", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListOpenCheckouts))).Methods(http.MethodGet)

	list.Handle("/videos/{limit}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListVideos))).Methods(http.MethodGet)
	list.Handle("/playlists", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListPlaylists))).Methods(http.MethodGet)
	list.Handle("/links", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListLinks))).Methods(http.MethodGet)

	// Public Lists
	pubList := r.PathPrefix("/list").Subrouter()
	pubList.HandleFunc("/mobileUsers", routers.HandleListMobileUsers).Methods(http.MethodGet)

	// V2.1 Handlers
	admin.Handle("/video/{query}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleGetVideo))).Methods(http.MethodGet)
	admin.Handle("/videos", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListVideo))).Methods(http.MethodGet)

	admin.Handle("/asset/{query}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleGetAsset))).Methods(http.MethodGet)
	admin.Handle("/asset/{query}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleEditAsset))).Methods(http.MethodPut)
	admin.Handle("/assets", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListAsset))).Methods(http.MethodGet)

	admin.Handle("/checkouts", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListCheckouts))).Methods(http.MethodGet)

	admin.Handle("/upload", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleUpload))).Methods(http.MethodPost)

	// Launch API Listener
	fmt.Printf("✅ Hillview Core API running on port %s\n", env.Port)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Origin", "Authorization", "Accept", "X-CSRF-Token"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	log.Fatal(http.ListenAndServe(":"+env.Port, handlers.CORS(originsOk, headersOk, methodsOk)(primary)))
}
