package main

import (
	"encoding/xml"
	"database/sql"
	"fmt"
	"log"
	"net/http"
    "time"
    "reflect"
    "io/ioutil"

	_ "github.com/go-sql-driver/mysql"
)

// Settings is for settings
type Settings struct {
	DBPassword string `xml:"db-password"`
}

var settings Settings

func readSettings() error {
    // read Rep.xml
    dat, err := ioutil.ReadFile("settings.xml")
    if err != nil {
        return err
    }

    // parse the repository XML
    err = xml.Unmarshal(dat, &settings)
    if err != nil {
        return err
    }

    return nil
}

func openDB() *sql.DB {
	db, err := sql.Open("mysql", "npackd:" + settings.DBPassword + "@/npackd")
	if err != nil {
		panic(err)
    }
    
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(10)

    err = db.Ping()
	if err != nil {
		panic(err)
    }

    return db
}

func updateDB() {
	db := openDB()
    
    defer db.Close()

    db.Exec(`
    CREATE TABLE if not exists PACKAGE (
        NAME varchar(255) NOT NULL,
        TITLE varchar(1024) NOT NULL,
        URL varchar(2048) DEFAULT NULL,
        ICON varchar(2048) DEFAULT NULL,
        DESCRIPTION varchar(4096) NOT NULL,
        LICENSE varchar(255) DEFAULT NULL,
        FULLTEXT_ varchar(4096) NOT NULL,
        STATUS int DEFAULT NULL,
        SHORT_NAME varchar(255) DEFAULT NULL,
        REPOSITORY int DEFAULT NULL,
        CATEGORY0 int DEFAULT NULL,
        CATEGORY1 int DEFAULT NULL,
        CATEGORY2 int DEFAULT NULL,
        CATEGORY3 int DEFAULT NULL,
        CATEGORY4 int DEFAULT NULL,
        TITLE_FULLTEXT varchar(1024) NOT NULL,
        STARS int DEFAULT NULL
      )    
    `)
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	db := openDB()
    
    defer db.Close()

	// Execute the query
    rows, err := db.Query("SELECT * FROM PACKAGE")
    defer rows.Close()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Get column names
	colNames, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

    // Get column names
    /*
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
    }
    */

    var myMap = make(map[string]interface{})
    cols := make([]interface{}, len(colNames))
    colPtrs := make([]interface{}, len(colNames))
    for i := 0; i < len(colNames); i++ {
        colPtrs[i] = &cols[i]
    }
    for rows.Next() {
        err = rows.Scan(colPtrs...)
        if err != nil {
            log.Fatal(err)
        }
        for i, col := range cols {
            myMap[colNames[i]] = col
        }
        // Do something with the map
        for key, val := range myMap {
            fmt.Fprintln(w, "Key:", key, "Value Type:", reflect.TypeOf(val))
        }
    }

	fmt.Fprintf(w, "Hello!")
}

func main() {
    err := readSettings()
    if err != nil {
        panic(err)
    }

    updateDB()

	http.HandleFunc("/api", handleFunc)
	log.Println("Server is running:", 3001)
	log.Fatal(http.ListenAndServe(":3001", nil))
}
