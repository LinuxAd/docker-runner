package runner

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
}

type GeneralErr interface {
	Error() string
}

type Err struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type Health struct {
	Msg string `json:"msg"`
}

type AppError struct {
	Err
}

type DockerError struct {
	Err
}

func newAppErr(msg string) *AppError {
	return &AppError{
		Err{
			Code: 1,
			Msg:  msg,
		}}
}

func (ae *AppError) Error() string {
	return ae.Msg
}

func (de *DockerError) Error() string {
	return de.Msg
}

func (a *App) Run(addr string) {
	go runningTicker()
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

func runningTicker() {
	uptimeTicker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-uptimeTicker.C:
			go CheckRunning()
		}
	}
}

func (a *App) Init() {
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/healthz", a.healthz).Methods("GET")
	a.Router.HandleFunc("/services", a.getServices).Methods("GET")
	a.Router.HandleFunc("/services", a.newService).Methods("POST")
}

func (a *App) healthz(w http.ResponseWriter, r *http.Request) {
	h := Health{Msg: "app responding"}
	respondWithJSON(w, http.StatusOK, h)
}

func (a *App) getServices(w http.ResponseWriter, r *http.Request) {
	var err Error

	if len(Running) <= 0 {
		err = Error{
			Msg:  "no services running",
			Body: "no services running",
		}
	}
	respondWithJSON(w, http.StatusOK, Response{
		Services: Running,
		Error:    err,
	})
}

func (a *App) newService(w http.ResponseWriter, r *http.Request) {
	var s Service
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&s); err != nil {
		respondWithError(w, http.StatusBadRequest, newAppErr("Invalid Request Paylod"))
		return
	}
	defer r.Body.Close()

	if err := s.newService(); err != nil {
		respondWithError(w, http.StatusBadRequest, newAppErr("Invalid request payload"))
	}

	respondWithJSON(w, http.StatusOK, s)
}

func respondWithError(w http.ResponseWriter, code int, err GeneralErr) {
	respondWithJSON(w, code, err)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
