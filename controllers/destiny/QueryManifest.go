package destiny

import (
	"database/sql"
	"encoding/json"

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
	Type     string `json:"type"`
	Response string `json:"response"`
}

func DestinyManifestQuery(id, tablename string) (Message, Item) {
	db, error := sql.Open("sqlite3", "./controllers/destiny/manifest/manifest.db")
	if error != nil {
		m := Message{Type: "Error", Response: "Unable to load destiny manifest file"}
		return m, Item{}
	}

	rows, error := db.Query("SELECT * FROM " + tablename + " WHERE hash='" + id + "';")
	if error != nil {

		m := Message{Type: "Error", Response: "Unable to query the destiny manifest"}
		return m, Item{}
	}

	var idd string
	var jsondata string
	for rows.Next() {
		rows.Scan(&idd, &jsondata)
		var data Item
		json.Unmarshal([]byte(jsondata), &data)

		return Message{}, data
	}
	return Message{}, Item{}

}
