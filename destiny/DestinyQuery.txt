package main

import (
	"database/sql"
	_"os"
	"fmt"
	"strconv"
	_"github.com/mattn/go-sqlite3"
)

func main()  {
	db, error := sql.Open("sqlite3", "Manifest/manifest.content/world_sql_content_3d029e66883b2c5765b6e4848f1c2965.content")
	if error != nil {
		fmt.Println(error)
	}
	fmt.Println("------fetching from sqlite file------")
	rows, error := db.Query("SELECT * FROM DestinyInventoryItemDefinition WHERE id='-2146672205';")
	if error != nil {
		fmt.Println(error)
	}
	var id int
	var json string
	for rows.Next() {
		rows.Scan(&id, &json)

		fmt.Println("ID: " + strconv.Itoa(id))
		fmt.Println(json)
	}

}

