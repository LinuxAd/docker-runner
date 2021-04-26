package runner

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
}

type GeneralErr interface {
	Error() string
}

type Err struct {
	Code int
	Msg  string
}

type Health struct {
	Msg string
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

func newDBError(msg string) *DockerError {
	return &DockerError{
		Err{
			Code: 0,
			Msg:  msg,
		},
	}
}

func (ae *AppError) Error() string {
	return ae.Msg
}

func (de *DockerError) Error() string {
	return de.Msg
}

func (a *App) Init() {
	a.Router = mux.NewRouter()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/healthz", a.healthz).Methods("GET")
	a.Router.HandleFunc("/services", a.getServices).Methods("GET")

}

func (a *App) healthz(w http.ResponseWriter, r *http.Request) {
	h := Health{Msg: "app responding"}
	respondWithJSON(w, http.StatusOK, h)
}

func (a *App) getServices(w http.ResponseWriter, r *http.Request) {
	if len(Running) <= 0 {
		respondWithJSON(w, http.StatusOK, Response{
			Error: Error{
				Msg:  "no services running",
				Body: "no services running",
			},
		})
	}
	respondWithJSON(w, http.StatusOK, Response{
		Services: Running,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
