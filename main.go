package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// start server, listening port 9090
func startServer() {
	// index file
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/", http.StatusFound)
	}) //设置访问的路由

	// static file
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	// other logic handlers
	http.HandleFunc("/rank", rank)
	http.HandleFunc("/top", top)
	//	http.HandleFunc("/update", update)

	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func rank(w http.ResponseWriter, r *http.Request) {
	log.Println("rank")

	if r.Method == "POST" {
		data := r.FormValue("data")

		str, err := base64Decode(data)
		if err != nil {
			fmt.Fprintf(w, "[]")
			return
		}

		values, err := url.ParseQuery(str)
		if err != nil {
			fmt.Fprintf(w, "[]")
			return
		}
		level, err := strconv.Atoi(values.Get("level"))
		if err != nil {
			fmt.Fprintf(w, "[]")
			return
		}

		secs, err := strconv.ParseFloat(values.Get("secs"), 64)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "{}")
			return
		}
		name := values.Get("name")

		r := getRank(secs, level)
		id := insertScore(secs, level, name)

		type RankRet struct {
			Id   int `json:"id"`
			Rank int `json:"rank"`
		}
		json_str, err := json.Marshal(RankRet{Id: id, Rank: r})
		if err != nil {
			log.Println(err)
		}
		//		log.Println(str)
		fmt.Fprintf(w, string(json_str))
	}

}

func top(w http.ResponseWriter, r *http.Request) {
	log.Println("top")

	if r.Method == "POST" {
		data := r.FormValue("data")

		str, err := base64Decode(data)
		if err != nil {
			fmt.Fprintf(w, "[]")
			return
		}

		values, err := url.ParseQuery(str)
		if err != nil {
			fmt.Fprintf(w, "[]")
			return
		}

		level, err := strconv.Atoi(values.Get("level"))
		if err != nil {
			fmt.Fprintf(w, "[]")
			return
		}
		id, err := strconv.Atoi(values.Get("id"))
		name := values.Get("name")
		if err == nil && name != "" {
			if len(name) > 16 {
				name = name[0:16]
			}
			updateName(name, id)
		}

		fmt.Fprintf(w, getTop(40, level))
	}
}

func update(w http.ResponseWriter, r *http.Request) {
	log.Println("top")
	if r.Method == "POST" {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			fmt.Fprintf(w, "{}")
			return
		}
		name := r.FormValue("name")
		//		log.Println(name)
		if len(name) > 0 {
			updateName(name, id)
		}
		fmt.Fprintf(w, "{}")
	}
}

const (
	base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_.="
)

func base64Decode(src string) (string, error) {
	var coder = base64.NewEncoding(base64Table)
	res, err := coder.DecodeString(string(src))
	return string(res), err
}

func main() {
	str, err := base64Decode("bGV2ZWw9MyZpZD03NyZuYW1lPeWwj_S8mQ==")
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	values, err := url.ParseQuery(str)

	fmt.Println(str)
	fmt.Println(values.Get("secs1"))
	fmt.Println(values.Get("name"))
	fmt.Println(values.Get("level"))
	//	createTable()
	//	id := insertScore(3.33, 3)
	//	fmt.Println(id)

	//	fmt.Println(getRank(35, 3))
	//	fmt.Println(getTop(10, 3))
	//	updateName("邝昌浪", 9)
	startServer()
}
