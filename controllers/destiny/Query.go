package destiny

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Request struct {
	Id int `json:"id"`
}

type Data struct {
	Id   string `json:"id"`
	Json string `json:"json"`
}

type Message struct {
	Response string `json:"response"`
}

func DestinyManifestQuery(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request Request
	error := json.NewDecoder(router.Body).Decode(&request)
	if error != nil {
		m := Message{Response: "Please send a valid json object"}
		json.NewEncoder(w).Encode(m)
		return
	}
	db, error := sql.Open("sqlite3", "Manifest/manifest.content/world_sql_content_3d029e66883b2c5765b6e4848f1c2965.content")
	if error != nil {
		m := Message{Response: "Unable to load destiny manifest file"}
		json.NewEncoder(w).Encode(m)
		return
	}

	rows, error := db.Query("SELECT * FROM DestinyInventoryItemDefinition WHERE id='-2146672205';")
	if error != nil {
		m := Message{Response: "Unable to query the destiny manifest"}
		json.NewEncoder(w).Encode(m)
		return
	}

	var id int
	var jsondata string
	for rows.Next() {
		rows.Scan(&id, &jsondata)
		data := Data{Id: strconv.Itoa(id), Json: jsondata}
		json.NewEncoder(w).Encode(data)
		return
	}

}
