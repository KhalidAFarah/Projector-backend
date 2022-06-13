package youtube

import (
	"net/http"
	//"github.com/gorilla/mux"
	"io/ioutil"
	"fmt"
	"log"
	"encoding/json"
)


type Playlist struct {
	Etag  string `json:"etag"`
	Items []struct {
		ContentDetails struct {
			RelatedPlaylists struct {
				Likes   string `json:"likes"`
				Uploads string `json:"uploads"`
			} `json:"relatedPlaylists"`
		} `json:"contentDetails"`
		Etag string `json:"etag"`
		ID   string `json:"id"`
		Kind string `json:"kind"`
	} `json:"items"`
	Kind     string `json:"kind"`
	PageInfo struct {
		ResultsPerPage int `json:"resultsPerPage"`
		TotalResults   int `json:"totalResults"`
	} `json:"pageInfo"`
}

func GetPlaylist(w http.ResponseWriter, router *http.Request){
	w.Header().Set("Content-Type", "application/json")

	token := router.URL.Query().Get("token")
	playlist := router.URL.Query().Get("playlist")
	next := router.URL.Query().Get("next")
	//token := mux.Vars(router)["token"]
	//next := mux.Vars(router)["next"]
	


	client := http.Client{}

	if playlist == "" {
		request, error := http.NewRequest("GET", "https://content-youtube.googleapis.com/youtube/v3/channels?part=contentDetails&mine=true", nil)
		if error != nil {
			log.Fatal("Unable to create request")
		}
		request.Header.Add("Accept", "application/json")
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer " + token)
		response, error := client.Do(request)

		if error != nil {
			log.Fatal("Unable to send request to bungie")
		}

		defer response.Body.Close()
		body, error := ioutil.ReadAll(response.Body)
		if error != nil {
			log.Fatal("Unable to read mobileworld data")
		}

		var data Playlist
		json.Unmarshal(body, &data)
		if error != nil {
			log.Fatal("Unable to read mobileworld data")
		}

		playlist = data.Items[0].ContentDetails.RelatedPlaylists.Uploads
	}
	
	var url string

	
	if next == "" {
		url = "https://youtube.googleapis.com/youtube/v3/playlistItems?part=contentDetails&part=contentDetails,snippet&maxResults=50&playlistId="+playlist
	}else{
		url = "https://youtube.googleapis.com/youtube/v3/playlistItems?part=contentDetails&part=contentDetails,snippet&maxResults=50&playlistId="+playlist + "&pageToken="+next
	}
	fmt.Println(url)

	request, error := http.NewRequest("GET", url, nil)
	if error != nil {
		log.Fatal("Unable to create request")
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer " + token)
	response, error := client.Do(request)

	if error != nil {
		log.Fatal("Unable to send request to bungie")
	}

	defer response.Body.Close()
	body, error := ioutil.ReadAll(response.Body)
	if error != nil {
		log.Fatal("Unable to read mobileworld data")
	}

	videos := make(map[string]interface{})
	json.Unmarshal(body, &videos)
	if error != nil {
		log.Fatal("Unable to read mobileworld data")
	}

	//fmt.Println(videos)

	json.NewEncoder(w).Encode(videos)

	


}