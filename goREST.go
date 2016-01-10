//Import required packages
package main
import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"net/http"
	"github.com/gocql/gocql"
	"log"
	"github.com/gorilla/securecookie"
)

//Declare MUX variables
var router = mux.NewRouter()
var decoder = schema.NewDecoder()

//Routes
func main() {
	router.HandleFunc("/", indexPage)
	router.HandleFunc("/login", clientLogin)
	
	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}

//Create a struct
type Login struct {
	Email  string
   	 Password string
 }

//PROCESS THE FORM DATA FUNCTION
func clientLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")


	//PARSE FROM VALUES
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}

	//EXTRACT FORM VALUES & DECLARE VARIABLES
	email := r.FormValue("email")
	password := r.FormValue("password")

	login := new(Login)
	// r.PostForm is a map of our POST form values
	err = decoder.Decode(login, r.PostForm)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(login)

	
	//CASSANDRA DATABASE CONNECTION
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "imontessori"
	cqlsession, _ := cluster.CreateSession()
	defer cqlsession.Close()

	//Query data from cluster
	var firstname string
	if err := cqlsession.Query("SELECT firstname FROM users3 WHERE email='" + email + "' AND password='" + password + "'").Scan(&firstname); err != nil {
		log.Fatal(err)
	}

	fmt.Println(email)
}
