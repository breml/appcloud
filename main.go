package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/breml/appcloud/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/breml/appcloud/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"

	"github.com/cloudfoundry-community/go-cfenv"
)

type Temperature struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Temperature string        `bson:"temperature"`
}

var (
	mongouri = "mongodb://localhost:27017/weatherDB"
	mongodb  = "weatherDB"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial(mongouri)
	if err != nil {
		fmt.Printf("MongoDB dial err %v\n", err)
		return
	}
	defer session.Close()

	c := session.DB(mongodb).C("weatherCOLL")

	var result Temperature
	curTemp := "undef"

	err = c.Find(nil).One(&result)
	if err != nil {
		fmt.Printf("MongoDB find err %v\n", err)
	} else {
		curTemp = result.Temperature
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<h1>How is the weather ? </h1>")
	fmt.Fprintf(w, "Hello world! <br><hr> The weather is "+curTemp+" today")
	fmt.Fprintf(w, "<br><hr><form action='/temp'><h3>Change weather</h3><br><input type=text name=temp value=hot><br><input type=submit value='Change the weather...'></form>")
}

func tempHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	queryValues := r.URL.Query()

	if len(queryValues["temp"]) > 0 && queryValues["temp"][0] != "" {
		re := regexp.MustCompile("[^a-zA-Z0-9\\ ]")
		temperature := re.ReplaceAllString(queryValues["temp"][0], "")

		session, err := mgo.Dial(mongouri)
		if err != nil {
			fmt.Printf("MongoDB dial err %v\n", err)
			return
		}
		defer session.Close()

		c := session.DB(mongodb).C("weatherCOLL")

		err = c.DropCollection()
		if err != nil {
			fmt.Printf("MongoDB drop collection err %v\n", err)
		}

		err = c.Insert(Temperature{Temperature: temperature})
		if err != nil {
			fmt.Printf("MongoDB insert err %v\n", err)
		} else {
			fmt.Fprintf(w, "The weather is now: %s\n", temperature)
		}
		fmt.Fprintf(w, "done <hr> <a href='/'>back</a>")
	} else {
		fmt.Fprintf(w, "error: no new temp set <hr> <a href='/'>back</a>")
	}
}

func main() {
	// Get MongoDB settings from ENV
	appEnv, err := cfenv.Current()
	if err == nil {
		// In App Cloud
		mgoService, err := appEnv.Services.WithName("mongodb")
		if err == nil {
			var ok bool
			mongouri, ok = mgoService.Credentials["uri"].(string)
			if !ok {
				fmt.Printf("No valid MongoDB uri\n")
				os.Exit(1)
			}
			mongodb, ok = mgoService.Credentials["database"].(string)
			if !ok {
				fmt.Printf("No valid MongoDB database\n")
				os.Exit(1)
			}
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Use port: %s\n", port)

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/temp", tempHandler)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("ListenAndServe error: %v\n", err)
		os.Exit(1)
	}
}
