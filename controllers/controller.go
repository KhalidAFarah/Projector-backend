package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"


	"github.com/Khalidium/projector/controllers/functions"
	"github.com/Khalidium/projector/controllers/destiny"
)

//router
var router mux.Router


func Start() {
	router := mux.NewRouter()
	//cors
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"https://proteje.netlify.app/*"})

	//endpoints
	router.HandleFunc("/api/", functions.Front).Methods("GET")
	router.HandleFunc("/api/sup/", functions.Sup).Methods("GET")
	router.HandleFunc("/api/gamesshow", functions.Gamesshow).Methods("GET")

	router.HandleFunc("/api/destiny/generatemanifest/", destiny.GenerateManifest).Methods("GET")
	router.HandleFunc("/api/destiny/query/", destiny.DestinyManifestQuery).Methods("GET")
	


	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(credentials, methods, origins)(router))) //+os.Getenv("PORT")
	
}

func addEndpoint(endpoint string, menthod string, function func(w http.ResponseWriter, router *http.Request)) {
	router.HandleFunc(endpoint, function).Methods(menthod)
}
