package main

import (
	"Crud_RestApi_With_Mux_GORM/controllers"
	"Crud_RestApi_With_Mux_GORM/database"
	"fmt"

	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	//"gorm.io/gorm"
)

var (
	router *mux.Router
	//secretkey string = "secretkeyjwt"
)

func main() {

	// Load Configurations from config.json using Viper
	LoadAppConfig()
	fmt.Println("LoadAppConfig is done")

	// Initialize Database
	InitializeDatabase()
	fmt.Println("InitializeDatabase is done")
	//CreateRouter

	CreateRouter()
	fmt.Println("CreateRouter is done")

	// Register Routes
	InitializeRoute()
	fmt.Println("InitializeRoute is done")

	ServerStart()

}

/////////////////////////////////////////////////////////////////////////////////////////////////////

func InitializeDatabase() {
	database.Connect(AppConfig.ConnectionString)
	database.Migrate()
}

func CreateRouter() {
	router = mux.NewRouter().StrictSlash(true)
}

func InitializeRoute() {

	//router.HandleFunc("/", controllers.HomePage)

	router.HandleFunc("/signup", controllers.SignUp).Methods("POST")
	router.HandleFunc("/signin", controllers.SignIn).Methods("POST")
	router.HandleFunc("/admin", controllers.IsAuthorized(controllers.AdminIndex)).Methods("GET")
	router.HandleFunc("/user", controllers.IsAuthorized(controllers.UserIndex)).Methods("GET")
	router.HandleFunc("/", controllers.Index).Methods("GET")
	

	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	})

	//func routerForAdmins()

	router.HandleFunc("/api/contacts", controllers.AddContact).Methods("POST")
	router.HandleFunc("/api/contacts", controllers.IsAuthorized(controllers.GetContacts)).Methods("GET")
	router.HandleFunc("/api/contacts/{id}", controllers.GetContactByIDs).Methods("GET")
	router.HandleFunc("/api/contacts/{id}", controllers.UpdateContact).Methods("PUT")
	router.HandleFunc("/api/contacts/{id}", controllers.DeleteContact).Methods("DELETE")

}

// start the server
func ServerStart() {
	log.Println(fmt.Sprintf("Starting Server on port %s", AppConfig.Port))
	err := http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))
	if err != nil {
		log.Fatal(err)
		fmt.Println("error in run server")
	}

	// old
	// log.Println(fmt.Sprintf("Starting Server on port %s", AppConfig.Port))
	// log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", AppConfig.Port), router))
}
