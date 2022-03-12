package main

import (
	"projector/controllers"
	//	_ "projector/controllers/destiny/model"
)

func main() {
	controllers.Start()
	//user := model.User{}
	//user.InitUser()
	//endpoints
	/*
		router.HandleFunc("/api/generatemanifest/", generate.generateManifest).Methods("GET")
		router.HandleFunc("/api/destinyquery/", query.query).Methods("GET")*/

}
