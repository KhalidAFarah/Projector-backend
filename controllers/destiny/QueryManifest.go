package destiny

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
)

type Request struct {
	Id int `json:"id"`
}

type Data struct {
	Id   string `json:"id"`
	Json string `json:"json"`
}

type Message struct {
	Type string `json:"type"`
	Response string `json:"response"`
}

func DestinyManifestQuery(id, tablename string) (Message, string){
	db, error := sql.Open("sqlite3", "./controllers/destiny/manifest/world_sql_content_c1d4ac435e5ce5b3046fe2d0e6190ce4.content")
	if error != nil {
		m := Message{Type: "Error", Response: "Unable to load destiny manifest file"}
		return m, ""
	}
	
	rows, error := db.Query("SELECT * FROM " + tablename + " WHERE id='" + id + "';")
	if error != nil {
		fmt.Println(error)
		m := Message{Type: "Error", Response: "Unable to query the destiny manifest"}
		return m, ""
	}

	var idd int
	var jsondata string
	for rows.Next() {
		rows.Scan(&idd, &jsondata)

		return Message{}, jsondata
	}
	return Message{}, ""

}