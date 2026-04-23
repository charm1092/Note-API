package http

import (
	"net/http"
	"note_api/note"
)

type HTTPHandlers struct {
	repo *note.Repository 
}

func NewHTTPHandlers(repo *note.Repository) *HTTPHandlers {
	return &HTTPHandlers{
		repo: repo,
	}
}

// /notes
// POST
// info: JSON in HTTP request body
func (h *HTTPHandlers) HandleCreateNote(w http.ResponseWriter, r *http.Request)  {
	
}

// /notes/{title}
// DELETE
// info: pattern
func (h *HTTPHandlers) HandleDeleteNote(w http.ResponseWriter, r *http.Request)  {
	
}

// /notes/{title}
// GET
// info: pattern
func (h *HTTPHandlers) HandleGetNote(w http.ResponseWriter, r *http.Request)  {
	
}

// /notes
// GET
// info: -
func (h *HTTPHandlers) HandleGetAllNotes(w http.ResponseWriter, r *http.Request)  {
	
}

// /notes/{title}
// PATCH
// info: pattern + JSON in request body
func (h *HTTPHandlers) HandleChangeNote(w http.ResponseWriter, r *http.Request)  {
	
}

// /notes/{title}/history
// GET
// info: pattern
func (h *HTTPHandlers) HandleGetHistoryVersionsOfNote(w http.ResponseWriter, r *http.Request)  {
	
}

// /notes/{title}/restore/{version}
// POST
// info: pattern
func (h *HTTPHandlers) HandleRestoreVersion(w http.ResponseWriter, r *http.Request)  {
	
}