package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Temperature struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Temperature string        `bson:"temperature"`
}

func debugHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<h1>Debug</h1>\n")

	fmt.Fprintf(w, "<h2>HW Specs</h2>")
	fmt.Fprintf(w, "<h3>CPU</h3>")
	fmt.Fprintf(w, "<table>")
	fmt.Fprintf(w, "<tr><td>NumCPU</td><td>%d</td></tr>\n", runtime.NumCPU())

	cpuinfo, err := cpu.CPUInfo()
	if err == nil {
		for _, cpus := range cpuinfo {
			fmt.Fprintf(w, "<tr><td>CPU</td><td>%v</td></tr>\n", cpus.CPU)
			fmt.Fprintf(w, "<tr><td>VendorID</td><td>%v</td></tr>\n", cpus.VendorID)
			fmt.Fprintf(w, "<tr><td>Family</td><td>%v</td></tr>\n", cpus.Family)
			fmt.Fprintf(w, "<tr><td>Model</td><td>%v</td></tr>\n", cpus.Model)
			fmt.Fprintf(w, "<tr><td>Stepping</td><td>%v</td></tr>\n", cpus.Stepping)
			fmt.Fprintf(w, "<tr><td>PhysicalID</td><td>%v</td></tr>\n", cpus.PhysicalID)
			fmt.Fprintf(w, "<tr><td>CoreID</td><td>%v</td></tr>\n", cpus.CoreID)
			fmt.Fprintf(w, "<tr><td>Cores</td><td>%v</td></tr>\n", cpus.Cores)
			fmt.Fprintf(w, "<tr><td>ModelName</td><td>%v</td></tr>\n", cpus.ModelName)
			fmt.Fprintf(w, "<tr><td>Mhz</td><td>%v</td></tr>\n", cpus.Mhz)
			fmt.Fprintf(w, "<tr><td>CacheSize</td><td>%v</td></tr>\n", cpus.CacheSize)
			fmt.Fprintf(w, "<tr><td>Flags</td><td>%v</td></tr>\n", cpus.Flags)
		}
		fmt.Fprintf(w, "</table>")
	} else {
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "Error while getting CPU info!\n")
	}

	fmt.Fprintf(w, "<h3>Virtual Memory</h3>")
	vmem, err := mem.VirtualMemory()
	if err == nil {
		fmt.Fprintf(w, "<table>")
		fmt.Fprintf(w, "<tr><td>Total</td><td>%v</td></tr>\n", vmem.Total)
		fmt.Fprintf(w, "<tr><td>Used</td><td>%v</td></tr>\n", vmem.Used)
		fmt.Fprintf(w, "<tr><td>UsedPercent</td><td>%v</td></tr>\n", vmem.UsedPercent)
		fmt.Fprintf(w, "<tr><td>Free</td><td>%v</td></tr>\n", vmem.Free)
		fmt.Fprintf(w, "<tr><td>Active</td><td>%v</td></tr>\n", vmem.Active)
		fmt.Fprintf(w, "<tr><td>Inactive</td><td>%v</td></tr>\n", vmem.Inactive)
		fmt.Fprintf(w, "<tr><td>Buffers</td><td>%v</td></tr>\n", vmem.Buffers)
		fmt.Fprintf(w, "<tr><td>Cached</td><td>%v</td></tr>\n", vmem.Cached)
		fmt.Fprintf(w, "<tr><td>Wired</td><td>%v</td></tr>\n", vmem.Wired)
		fmt.Fprintf(w, "<tr><td>Shared</td><td>%v</td></tr>\n", vmem.Shared)
		fmt.Fprintf(w, "</table>")
	} else {
		fmt.Fprintf(w, "Error while getting mem info!\n")
	}

	fmt.Fprintf(w, "<h3>Swap Memory</h3>")
	swapmem, err := mem.SwapMemory()
	if err == nil {
		fmt.Fprintf(w, "<table>")
		fmt.Fprintf(w, "<tr><td>Total</td><td>%v</td></tr>\n", swapmem.Total)
		fmt.Fprintf(w, "<tr><td>Used</td><td>%v</td></tr>\n", swapmem.Used)
		fmt.Fprintf(w, "<tr><td>Free</td><td>%v</td></tr>\n", swapmem.Free)
		fmt.Fprintf(w, "<tr><td>UsedPercent</td><td>%v</td></tr>\n", swapmem.UsedPercent)
		fmt.Fprintf(w, "<tr><td>Sin</td><td>%v</td></tr>\n", swapmem.Sin)
		fmt.Fprintf(w, "<tr><td>Sout</td><td>%v</td></tr>\n", swapmem.Sout)
		fmt.Fprintf(w, "</table>")
	} else {
		fmt.Fprintf(w, "Error while getting swap mem info!\n")
	}

	fmt.Fprintf(w, "</table>")

	fmt.Fprintf(w, "<h2>Environment</h2>\n")
	fmt.Fprintf(w, "<table>")
	env := os.Environ()
	for _, envvar := range env {
		pair := strings.Split(envvar, "=")
		fmt.Fprintf(w, "<tr><td>%s</td><td>%s</td></tr>\n", pair[0], pair[1])
	}
	fmt.Fprintf(w, "</table>")
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
	http.HandleFunc("/debug", debugHandler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("ListenAndServe error: %v\n", err)
		os.Exit(1)
	}
}
