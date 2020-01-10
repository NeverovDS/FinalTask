package Apiserver

import (
	"FinalTask/Internal/App/Store"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  Store.Store
}

func New(config *Config) APIServer {
	return APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}
func (s *APIServer) Open() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		s.logger.Println(err)
		return err
	}
	if err := db.Ping(); err != nil {
		s.logger.Println(err)
		return err
	}

	store := Store.Store{Db: db}
	s.store = store
	fmt.Println(store)
	fmt.Println("Successfully connected!")

	return nil
}
func (s APIServer) Start() error {

	_ = s.Open()

	if err := s.configureLogger(); err != nil {
		s.logger.Println(err)
		return err
	}

	s.configureRouter()

	s.logger.Info("starting api server")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}
func (s APIServer) configureRouter() {
	s.router.HandleFunc("/tweets", s.handleCreateTweet())
}

func (s APIServer) handleCreateTweet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "GET":
			account_id, ok := r.URL.Query()["account_id"]
			if !ok || len(account_id[0]) < 1 {
				http.ServeFile(w, r, "Internal/App/Apiserver/form.html")
				return
			}
			fmt.Println(account_id[0])
			accIDinINT, err := strconv.Atoi(account_id[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			UserTweet, err := s.store.GetByAccountID(accIDinINT)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Fprintf(w, "UserTweet = %v\n", UserTweet)

		case "POST":

			// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			_, err := fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
			if err != nil {
				fmt.Println(err)
				return
			}
			message := r.FormValue("message")
			account_id := r.FormValue("account_id")
			account_idINint, err := strconv.Atoi(account_id)
			if err != nil {
				return
			}
			SomeTweet := Store.NewTweet{
				Message:   message,
				AccountId: account_idINint,
			}

			SaveTweet, err := s.store.Save(SomeTweet)

			if err != nil {
				s.logger.Println(err)
				return
			}

			fmt.Fprintf(w, "ID = %d\n", SaveTweet.Id)
			fmt.Fprintf(w, "Message = %s\n", SaveTweet.Message)
			fmt.Fprintf(w, "AccountId = %d\n", SaveTweet.AccountId)
		default:
			fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
		}
	}
}
