package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type settings struct {
	DBPassword string `xml:"db-password"`
}

type Package struct {
	Name     string   `xml:"name,attr"`
	Title      string `xml:"title"`
	Description      string `xml:"description"`
	Category []string `xml:"category"`
	Tag      []string `xml:"tag"`
}

type PackageVersion struct {
	Name    string `xml:"name,attr"`
	Package string `xml:"package,attr"`
	Url     string `xml:"url"`
}

type Repository struct {
	PackageVersion []PackageVersion `xml:"version"`
	Package        []Package        `xml:"package"`
}

var programSettings settings

func parseRepository() (*Repository, error) {
	// read Rep.xml
	dat, err := ioutil.ReadFile("Rep.xml")
	if err != nil {
		return nil, err
	}

	// parse the repository XML
	v := Repository{}
	err = xml.Unmarshal(dat, &v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func insertPackages(packages []Package) error {
	db, err := openDB()
	if err != nil {
		return err
	}

	defer db.Close()

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO PACKAGE(NAME, TITLE, DESCRIPTION, FULLTEXT_, TITLE_FULLTEXT) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	for i := 0; i < len(packages); i++ {
		p := packages[i]
		_, err = stmtIns.Exec(p.Name, p.Title, p.Description, p.Description, p.Title)
		if err != nil {
			return err
		}
	}

	return nil
}

func insertPackageVersions(pvs []PackageVersion) error {
	db, err := openDB()
	if err != nil {
		return err
	}

	defer db.Close()

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO PACKAGE_VERSION(NAME, PACKAGE) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	for i := 0; i < len(pvs); i++ {
		p := pvs[i]

		// TODO: normalize version number

		_, err = stmtIns.Exec(p.Name, p.Package)
		if err != nil {
			return err
		}
	}

	return nil
}

func readSettings() error {
	// read Rep.xml
	dat, err := ioutil.ReadFile("settings.xml")
	if err != nil {
		return err
	}

	// parse the repository XML
	err = xml.Unmarshal(dat, &programSettings)
	if err != nil {
		return err
	}

	return nil
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "npackd:"+programSettings.DBPassword+"@/npackd")
	if err != nil {
		return nil, err
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func updateDB() error {
	db, err := openDB()
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec(`
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
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE if not exists PACKAGE_VERSION(
			NAME varchar(255) NOT NULL, 
            PACKAGE varchar(255) NOT NULL, URL varchar(2048), 
            CONTENT BLOB)
			`)
	if err != nil {
		return err
	}
				
	return nil
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	db, err := openDB()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Execute the query
	rows, err := db.Query("SELECT NAME, TITLE, DESCRIPTION FROM PACKAGE")
	defer rows.Close()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Get column names
	/*colNames, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	*/

	// Get column names
	/*
			colTypes, err := rows.ColumnTypes()
			if err != nil {
				panic(err.Error()) // proper error handling instead of panic in your app
		    }
	*/

	/*
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
	*/

	var p Package
	for rows.Next() {
		err = rows.Scan(&p.Name, &p.Title, &p.Description)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(w, "Name:", p.Name, "Title:", p.Title)
	}

	fmt.Fprintf(w, "Hello!")
}

func prog() error {
	err := readSettings()
	if err != nil {
		return err
	}

	err = updateDB()
	if err != nil {
		return err
	}

	
	/*
	rep, err := parseRepository()
	if err != nil {
		return err
	}

	err = insertPackages(rep.Package)
	if err != nil {
		return err
	}

	err = insertPackageVersions(rep.PackageVersion)
	if err != nil {
		return err
	}
	*/

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	http.HandleFunc("/api", handleFunc)
	log.Println("Starting the server on :", 3001)

	err = http.ListenAndServe(":3001", nil)

	return err
}

func main() {
	err := prog()
	if err != nil {
		panic(err)
	}
}
