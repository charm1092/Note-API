package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpHandlers *HTTPHandlers 
}

func NewHTTPServer(httpHandler *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: httpHandler,
	}
}

func (s *HTTPServer) StartServer() error {
	router := mux.NewRouter()

	router.Path("/notes").Methods("POST").HandlerFunc(s.httpHandlers.HandleCreateNote)
	router.Path("/notes/{title}").Methods("DELETE").HandlerFunc(s.httpHandlers.HandleDeleteNote)
	router.Path("/notes/{title}").Methods("GET").HandlerFunc(s.httpHandlers.HandleGetNote)
	router.Path("/notes").Methods("GET").HandlerFunc(s.httpHandlers.HandleGetAllNotes)
	router.Path("/notes/{title}").Methods("PATCH").HandlerFunc(s.httpHandlers.HandleChangeNote)
	router.Path("/notes/{title}/history").Methods("GET").HandlerFunc(s.httpHandlers.HandleGetHistoryVersionsOfNote)
	router.Path("/notes/{title}/restore/{version}").Methods("POST").HandlerFunc(s.httpHandlers.HandleRestoreVersion)

	if err := http.ListenAndServe(":9091", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	}

	return nil
}