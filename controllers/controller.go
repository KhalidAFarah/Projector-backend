package controllers

import (
	"log"
	"net/http"
	//"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"projector/controllers/destiny"
	"projector/controllers/functions"
	"projector/controllers/youtube"
)

//router
var router mux.Router

func Start() {
	router := mux.NewRouter()
	//cors
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"https://proteje.netlify.app/*", "https://proteje.netlify.app", "https://proteje.herokuapp.com/*", "https://proteje.herokuapp.com", "*"})

	//endpoints
	router.HandleFunc("/api/", functions.Front).Methods("GET")
	router.HandleFunc("/api/sup/", functions.Sup).Methods("GET")
	router.HandleFunc("/api/gamesshow", functions.Gamesshow).Methods("GET")
	router.HandleFunc("/api/youtube/", youtube.GetPlaylist).Methods("GET")

	router.HandleFunc("/api/destiny/generatemanifest/", destiny.GenerateManifest).Methods("GET")
	router.HandleFunc("/api/destiny/builds/", destiny.GetBuilds).Methods("GET")
	//router.HandleFunc("/api/destiny/query/", destiny.DestinyManifestQuery).Methods("GET")

	log.Fatal(http.ListenAndServe(":9200", handlers.CORS(credentials, methods, origins)(router)))//

}

func addEndpoint(endpoint string, menthod string, function func(w http.ResponseWriter, router *http.Request)) {
	router.HandleFunc(endpoint, function).Methods(menthod)
}
