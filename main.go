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

	// Public Lists
	pubList := r.PathPrefix("/list").Subrouter()
	pubList.HandleFunc("/mobileUsers", routers.V1HandleListMobileUsers).Methods(http.MethodGet)

	// Admin Handlers
	admin := r.PathPrefix("/admin").Subrouter()

	// V2.1 Handlers
	// videos
	admin.Handle("/video/{query}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleGetVideo))).Methods(http.MethodGet)
	admin.Handle("/video/{query}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleEditVideo))).Methods(http.MethodPut)
	admin.Handle("/videos", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListVideo))).Methods(http.MethodGet)
	admin.Handle("/video", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleCreateVideo))).Methods(http.MethodPost)

	// assets
	admin.Handle("/asset/{query}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleGetAsset))).Methods(http.MethodGet)
	admin.Handle("/asset/{query}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleEditAsset))).Methods(http.MethodPut)
	admin.Handle("/asset", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleCreateAsset))).Methods(http.MethodPost)
	admin.Handle("/assets", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListAsset))).Methods(http.MethodGet)

	// playlists
	admin.Handle("/playlists", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListPlaylists))).Methods(http.MethodGet)
	admin.Handle("/playlist/{query}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleUpdatePlaylist))).Methods(http.MethodPut)
	admin.Handle("/playlist/{query}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleGetPlaylist))).Methods(http.MethodGet)
	admin.Handle("/playlist", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleCreatePlaylist))).Methods(http.MethodPost)

	// links
	admin.Handle("/links", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListLinks))).Methods(http.MethodGet)
	admin.Handle("/link/{query}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleUpdateLink))).Methods(http.MethodPut)
	admin.Handle("/link/{query}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleGetLink))).Methods(http.MethodGet)
	admin.Handle("/link", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleCreateLink))).Methods(http.MethodPost)

	// checkouts
	admin.Handle("/checkouts", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListCheckouts))).Methods(http.MethodGet)
	admin.Handle("/checkout/{query}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleUpdateCheckout))).Methods(http.MethodPut)

	// users
	admin.Handle("/users", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListUsers))).Methods(http.MethodGet)
	admin.Handle("/user/{query}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleUpdateUser))).Methods(http.MethodPut)
	admin.Handle("/user/{query}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleGetUser))).Methods(http.MethodGet)

	// mobile users
	admin.Handle("/mobileUsers", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleListMobileUsers))).Methods(http.MethodGet)
	admin.Handle("/mobileUser/{query}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleGetMobileUser))).Methods(http.MethodGet)
	admin.Handle("/mobileUser/{query}", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleUpdateMobileUser))).Methods(http.MethodPut)
	admin.Handle("/mobileUser", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleCreateMobileUser))).Methods(http.MethodPost)

	// Upload Handler
	admin.Handle("/upload", middleware.AccessTokenMiddleware(http.HandlerFunc(routers.HandleUpload))).Methods(http.MethodPost)

	// Launch API Listener
	fmt.Printf("✅ Hillview Core API running on port %s\n", env.Port)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Origin", "Authorization", "Accept", "X-CSRF-Token"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	log.Fatal(http.ListenAndServe(":"+env.Port, handlers.CORS(originsOk, headersOk, methodsOk)(primary)))
}
