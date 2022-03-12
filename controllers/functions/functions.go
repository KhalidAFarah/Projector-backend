package functions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type Message struct {
	Response string `json:"response"`
}

type QnA []struct {
	Question string   `json:"Question"`
	Choices  []string `json:"Choices"`
	Answer   int      `json:"Answer"`
}

type Cards struct {
	Proj []struct {
		Img   string `json:"img"`
		Title string `json:"title"`
		Txt   string `json:"txt"`
		Link  string `json:"link"`
	} `json:"proj"`
	Exp []struct {
		Img   string `json:"img"`
		Title string `json:"title"`
		Txt   string `json:"txt"`
		Link  string `json:"link"`
	} `json:"exp"`
}

func Front(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	m := Message{Response: "This is the frontpage"}
	json.NewEncoder(w).Encode(m)
}
func Sup(w http.ResponseWriter, router *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	m := Message{Response: "Sup ✋"}
	json.NewEncoder(w).Encode(m)
}

func Gamesshow(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonFile, err := os.Open("../resources/QnA.json")
	if err != nil {
		m := Message{Response: "Unable to open file"}
		json.NewEncoder(w).Encode(m)
		return
	}
	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		m := Message{Response: "Unable to read file"}
		json.NewEncoder(w).Encode(m)
		return
	}
	data := QnA{}
	err2 := json.Unmarshal(jsonData, &data)
	if err2 != nil {
		m := Message{Response: "Unable to read json data"}
		json.NewEncoder(w).Encode(m)
		return
	}
	json.NewEncoder(w).Encode(data)
}

/*func init() {
	fmt.Println("sd")
	controllers.addEndpoint("/api/", "GET", front)
	controllers.addEndpoint("/api/sup/", "GET", sup)
	controllers.addEndpoint("/api/gamesshow/", "GET", gamesshow)
}*/
