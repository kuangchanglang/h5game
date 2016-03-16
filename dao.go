// db
package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// global variable
var db *sql.DB = nil

// if db is open
func isOpen() bool {
	return db != nil
}

// open db connection
func openDb() {
	if isOpen() {
		return
	}
	tmpDb, err := sql.Open("mysql", "user:111@/haha")

	if err != nil {
		log.Fatalf("Open database error: %s\n", err)
	}
	// Open doesn't open a connection. Validate DSN data:
	err = tmpDb.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	db = tmpDb
}

// close db connection
func closeDb() {
	db.Close()
}

func createTable() {
	openDb()
	sql := `create table score(id int primary key auto_increment, 
	secs double not null comment "seconds", 
	level int not null default 3 comment "level",
	name varchar(16) not null default 'noname',
	createtime timestamp not null default CURRENT_TIMESTAMP, 
	key score_key(secs)) 
	ENGINE=InnoDb, Charset=utf8;`
	_, err := db.Query(sql)
	if err != nil {
		log.Println(err)
	}
}

func insertScore(secs float64, level int, name string) int {
	if name == "" {
		name = "noname"
	}
	openDb()
	// Prepare statement for inserting data
	stmt, err := db.Prepare("INSERT INTO score(secs, level, name) VALUES(?, ?, ?)") // ? = placeholder
	if err != nil {
		log.Println(err)
		//		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmt.Close() // Close the statement when we leave main() / the program terminates
	_, err = stmt.Exec(secs, level, name)
	if err != nil {
		log.Println(err)
		//		panic(err.Error()) // proper error handling instead of panic in your app
	}
	var id int = -1
	db.QueryRow("select last_insert_id() from score;").Scan(&id)

	return id
}

func getRank(secs float64, level int) int {
	openDb()

	var rank = 101
	db.QueryRow("select count(id) from score where secs<? and level=?;", secs, level).Scan(&rank)

	return rank + 1
}

func getTop(count int, level int) string {
	openDb()

	rows, err := db.Query("select name, secs, createtime from score where level=? order by secs limit 0, ?;", level, count)
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()
	var score Score
	var scores []Score
	for rows.Next() {
		err := rows.Scan(&score.Name, &score.Secs, &score.Createtime)
		if err != nil {
			log.Println(err)
		}
		scores = append(scores, score)
	}

	str, err := json.Marshal(scores)
	err = rows.Err()
	if err != nil {
		log.Println(err)
	}
	return string(str)
}

func updateName(name string, id int) {
	openDb()
	log.Println(string(id))
	_, err := db.Query("update score set name=? where id=?;", name, id)
	if err != nil {
		log.Println(err)
	}

	return
}
