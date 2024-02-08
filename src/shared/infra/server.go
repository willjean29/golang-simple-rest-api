package server

import (
	"app/db"
	"app/middlewares"
	taskModels "app/src/modules/tasks/infra/gorm/entities"
	userModels "app/src/modules/users/infra/gorm/entities"
	error "app/src/shared/errors"
	"app/src/shared/infra/http/routes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Server struct {
	port   string
	router *mux.Router
}

func NewServer() (*Server, error.Error) {
	var defaultPort int64 = 4000
	port := strconv.FormatInt(defaultPort, 10)
	router := mux.NewRouter().StrictSlash(true)
	server := &Server{
		port:   port,
		router: router,
	}
	server.LoadEnv()
	server.DBConnection()
	server.Middlewares()
	server.Routes()
	return server, error.Error{}
}

func (server *Server) LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func (s *Server) DBConnection() {
	// connection database
	db.DBConnection()
	db.DB.AutoMigrate(&taskModels.Task{}, &userModels.User{})
}

func (s *Server) Routes() {
	// statics files to filesystem
	fs := http.FileServer(http.Dir("./uploads/"))
	s.router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", fs))

	// routes for application
	routes.Routes(s.router)
	s.router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(error.New("Endpoint not found", http.StatusNotFound, errors.New("Not found - "+r.URL.Path)))
	})

}

func (s *Server) Middlewares() {
	s.router.Use(middlewares.ContentType)
}

func (s *Server) Run() {
	log.Println("Running on port ", s.port)
	log.Fatal(http.ListenAndServe(":"+s.port, s.router))
}
