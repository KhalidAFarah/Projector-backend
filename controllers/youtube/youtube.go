package youtube

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
)

func GetPlaylist(w http.ResponseWriter, router *http.Request){
	w.Header().Set("Content-Type", "application/json")

	code := mux.Vars(router)["code"]

	

	fmt.Println(code)

	//https://content-youtube.googleapis.com/youtube/v3/channels?part=contentDetails&mine=true


}