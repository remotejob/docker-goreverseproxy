//go:generate  /home/juno/neonworkspace/gowork/bin/statik -src=./public

package main // import "github.com/remotejob/docker-goreverseproxy"

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"
	// _ "github.com/remotejob/godocker/statik"
)

var mongodbuser string
var mongodbpass string

//Employees title name
type Employees struct {
	Title string
	Name  string
}

func init() {

	log.Println("Start init")
	mongodbuser = os.Getenv("SECRET_USERNAME")
	mongodbpass = os.Getenv("SECRET_PASSWORD")

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		fmt.Println(pair[0], pair[1])
	}

	if _, err := os.Stat("/usr/share/nginx"); os.IsNotExist(err) {
		// path/to/whatever does not exist
		log.Println("/usr/share/nginx not exit ")

	} else {

		log.Println("/usr/share/nginx exist delete ")
		os.RemoveAll("/usr/share/nginx")

	}

}

func testhandler(w http.ResponseWriter, r *http.Request) {

	log.Println("testhandler 2")

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:     []string{"mymongo-controller"},
		Timeout:   60 * time.Second,
		Database:  "admin",
		Username:  mongodbuser,
		Password:  mongodbpass,
		Mechanism: "SCRAM-SHA-1",
	}

	session, err := mgo.DialWithInfo(mongoDBDialInfo)

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("node-mongo-employee").C("employees")

	result := []Employees{}
	err = c.Find(nil).All(&result)
	if err != nil {
		log.Fatal(err)
	}

	for _, empl := range result {
		fmt.Fprintf(w, "Hi  %s %s", empl.Name, empl.Title)
	}

}

func main() {
	// statikFS, err := fs.New()
	// if err != nil {
	// 	log.Fatalf(err.Error())
	// }

	http.HandleFunc("/test", testhandler)

	// // fs := http.FileServer(http.Dir("/home/juno/neonworkspace/gowork/src/github.com/remotejob/godocker/assets"))
	// fs := http.FileServer(http.Dir("assets"))

	// http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	// // http.Handle("/assets", http.FileServer(http.Dir("/home/juno/neonworkspace/gowork/src/github.com/remotejob/godocker/assets")))
	// http.Handle("/", http.FileServer(statikFS))
	http.ListenAndServe(":8080", nil)
}
