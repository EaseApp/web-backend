package main

import(
	// "golang.org/x/net/websocket"
	// "io"
	// "net"
	// "io/ioutil"
	"html/template"
	"fmt"
	"log"
	"net/http"
	"github.com/EaseApp/web-backend/src/db"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	user "github.com/EaseApp/web-backend/src/app/controllers/user"
	application "github.com/EaseApp/web-backend/src/app/controllers/application"
	dao "github.com/EaseApp/web-backend/src/app/dao"
)

func render(w http.ResponseWriter, tmpl string) {
    tmpl = fmt.Sprintf("src/app/views/%s", tmpl)
    t, err := template.ParseFiles(tmpl)
    if err != nil {
        log.Print("template parsing error: ", err)
    }
    err = t.Execute(w, "hello")
    if err != nil {
        log.Print("template executing error: ", err)
    }
}

func NewStaticUserHandler(w http.ResponseWriter, req *http.Request){
	_ = user.InsertStaticUser()
	http.Redirect(w, req, "/users", http.StatusFound)
}

func DBCountHandler(w http.ResponseWriter, req *http.Request){
	vars := mux.Vars(req)
	db := vars["db"]

	// Move RecordCount out of user DAO and into generic DAO
	fmt.Fprintf(w, "Db (%v) has (%v) objects.", db, dao.RecordCount(db))
}

func main() {
	log.Println("Starting server...")

	// TODO: Use command line flag credentials.
	client, err := db.NewClient("localhost:28015")
	if err != nil {
		log.Fatal("Couldn't initialize database: ", err.Error())
	}
	db.CreateEaseDb(client)
	db.CreateUserTable(client)
	db.CreateDbTable(client)
	// CreatePubSubServer()

	dao.Init(client.Session)
	user.Init(client.Session)
	application.Init(client.Session)


	router := mux.NewRouter()
	// router.Host("{listen}.domain.com").Path("/").HandlerFunc(PubSubHandler).Name("root")

	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		render(w, "index.html")
	})

	router.HandleFunc("/static/user/new", NewStaticUserHandler)

	// These should be POST, but it's easier to test with GET
	router.HandleFunc("/users/sign_in", user.SignInHandler)
	router.HandleFunc("/users/sign_up", user.SignUpHandler)

	router.HandleFunc("/count/{db}", DBCountHandler)
	router.HandleFunc("/{db}", user.FetchAllHandler)

	router.HandleFunc("/{client}/{application}", application.QueryApplicationHandler)
	router.HandleFunc("/{client}/{application}/new", application.CreateApplicationHandler).Methods("POST")
	// router.HandleFunc("/{client}/{application}/pubsub", websocket.Handler(EchoServer))
	router.HandleFunc("/{client}/{application}/{id}", application.UpdateApplicationHandler).Methods("PUT")
	router.HandleFunc("/{client}/{application}/{id}", application.DeleteApplicationHandler).Methods("DELETE")


	defer client.Close()

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3000")
}
