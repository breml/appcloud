package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/breml/appcloud/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/breml/appcloud/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"

	"github.com/cloudfoundry-community/go-cfenv"
)

const prefixDelim = ": "

type Temperature struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Temperature string        `bson:"temperature"`
}

// Initialize global service variables with default values for local
// deployment. If executed in App Cloud, these variables will be replaced
// with the values from ENV.
var (
	name = "appcloud"
	port = "3000"

	mongouri = "mongodb://localhost:27017/weatherDB"
	mongodb  = "weatherDB"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial(mongouri)
	if err != nil {
		log.Printf("MongoDB dial err %v\n", err)
		return
	}
	defer session.Close()

	c := session.DB(mongodb).C("weatherCOLL")

	var result Temperature
	curTemp := "undef"

	err = c.Find(nil).One(&result)
	if err != nil {
		log.Printf("MongoDB find err %v\n", err)
	} else {
		curTemp = result.Temperature
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintf(w, "<h1>How is the weather ? </h1>")
	fmt.Fprintf(w, "Hello world! <br><hr> The weather is %s today.", curTemp)
	fmt.Fprintf(w, "<br><hr><form action='/temp'><h3>Change weather</h3><br><input type=text name=temp value=hot><br><input type=submit value='Change the weather...'></form>")
}

func tempHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	queryValues := r.URL.Query()

	if len(queryValues["temp"]) > 0 && queryValues["temp"][0] != "" {
		re := regexp.MustCompile("[^a-zA-Z0-9\\ ]")
		temperature := re.ReplaceAllString(queryValues["temp"][0], "")

		session, err := mgo.Dial(mongouri)
		if err != nil {
			log.Printf("MongoDB dial err %v\n", err)
			return
		}
		defer session.Close()

		c := session.DB(mongodb).C("weatherCOLL")

		err = c.DropCollection()
		if err != nil {
			log.Printf("MongoDB drop collection err %v\n", err)
		}

		err = c.Insert(Temperature{Temperature: temperature})
		if err != nil {
			log.Printf("MongoDB insert err %v\n", err)
		} else {
			fmt.Fprintf(w, "The weather is now: %s\n", temperature)
		}
	} else {
		fmt.Fprintf(w, "error: no new temp set <hr> <a href='/'>back</a>")
	}
	fmt.Fprintf(w, "<hr> <a href='/'>back</a>")
}

func main() {
	// Setup logging
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	log.SetPrefix(name + prefixDelim)

	// Get settings from ENV if present (e. g. if in App Cloud)
	appEnv, err := cfenv.Current()
	if err == nil {
		// Port to bind web app
		port = strconv.Itoa(appEnv.Port)

		name = appEnv.Name
		log.SetPrefix(name + prefixDelim)

		// MongoDB Service
		mgoService, err := appEnv.Services.WithName("mongodb")
		if err == nil {
			var ok bool
			mongouri, ok = mgoService.Credentials["uri"].(string)
			if !ok {
				log.Fatalf("No valid MongoDB uri\n")
			}
			mongodb, ok = mgoService.Credentials["database"].(string)
			if !ok {
				log.Fatalf("No valid MongoDB database\n")
			}
		}
	}

	log.Printf("Use port: %s\n", port)

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/temp", tempHandler)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("ListenAndServe error: %v\n", err)
	}
}
