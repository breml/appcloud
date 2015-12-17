package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/breml/appcloud/Godeps/_workspace/src/labix.org/v2/mgo"
	"github.com/breml/appcloud/Godeps/_workspace/src/labix.org/v2/mgo/bson"
)

type Temperature struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Temperature string        `bson:"temperature"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial("mongodb://localhost:27017/weatherDB")
	if err != nil {
		fmt.Printf("MongoDB dial err %v\n", err)
		return
	}
	defer session.Close()

	var result Temperature

	c := session.DB("weatherDB").C("weatherCOLL")
	err = c.Find(nil).One(&result)
	if err != nil {
		fmt.Printf("MongoDB find err %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<h1>How is the weather ? </h1>")
	fmt.Fprintf(w, "Hello world! <br><hr> The weather is "+result.Temperature+" today")
	// res.write('<br><hr><form action=\'/temp\'><h3>Change weather</h3><br><input type=text name=temp value=hot><br><input type=submit value=\'Change the weather...\'></form>');
}

func tempHandler(w http.ResponseWriter, r *http.Request) {
	/*
	   	app.get('/temp', function (req, res) { // telling nodeJs to get all commands from /temp into this function

	   	mongoClient.connect( "mongodb://localhost:27017/weatherDB" , function(err, db) { // connect to the local database
	   	  	if(err) { return console.dir(err); } // check if connection is ok, else output

	   	  	db.collection('weatherCOLL').drop(); // drop the collection if existing
	     		db.collection('weatherCOLL').insert( {"temperature": req.query.temp } ) // add a new object with "temperature"
	   		res.send("done <hr> <a href=\'/\'>back</a>");
	   		res.end();
	     		db.close(); // close the Database connection
	     	});
	   });
	*/
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	queryValues := r.URL.Query()
	// TODO: check if weather is only chars: [a-zA-Z0-9]
	if len(queryValues["temp"]) > 0 && queryValues["temp"][0] != "" {
		fmt.Fprintf(w, "The weather is now: %s\n", queryValues["temp"][0])
		fmt.Fprintf(w, "done <hr> <a href='/'>back</a>")
	} else {
		fmt.Fprintf(w, "error: no new temp set <hr> <a href='/'>back</a>")
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Use port: %s\n", port)

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/temp", tempHandler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("ListenAndServe error: %v\n", err)
		os.Exit(1)
	}
}
