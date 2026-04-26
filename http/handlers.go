package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"note_api/note"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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
	ctx := r.Context()
	var noteDTO NoteDTO
	if err := json.NewDecoder(r.Body).Decode(&noteDTO); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time: time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := noteDTO.ValidateForCreate(); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time: time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	noteVar := note.NewNote(noteDTO.Title, noteDTO.Content)
	if err := h.repo.AddNote(ctx, noteVar); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, note.ErrNoteAlreadyExists) {
			http.Error(w, errDTO.ToString(), http.StatusConflict)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	b, err := json.MarshalIndent(noteVar, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}


}

// /notes/{title}
// DELETE
// info: pattern
func (h *HTTPHandlers) HandleDeleteNote(w http.ResponseWriter, r *http.Request)  {
	ctx := r.Context()
	title := mux.Vars(r)["title"]

	if err := h.repo.DeleteNote(ctx, title); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, note.ErrNoteNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// /notes/{title}
// GET
// info: pattern
func (h *HTTPHandlers) HandleGetNote(w http.ResponseWriter, r *http.Request)  {
	ctx := r.Context()
	title := mux.Vars(r)["title"]

	noteVar, err := h.repo.GetNote(ctx, title)
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, note.ErrNoteNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	b, err := json.MarshalIndent(noteVar, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

// /notes
// GET
// info: -
func (h *HTTPHandlers) HandleGetAllNotes(w http.ResponseWriter, r *http.Request)  {
	ctx := r.Context()
	notes, err := h.repo.ListNotes(ctx)
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, note.ErrNoteNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}
	b, err := json.MarshalIndent(notes, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

// /notes/{title}
// PATCH
// info: pattern + JSON in request body
func (h *HTTPHandlers) HandleChangeNote(w http.ResponseWriter, r *http.Request)  {
	ctx := r.Context()
	title := mux.Vars(r)["title"]
	currentTitle := title

	var dto UpdateNoteDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time: time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if dto.NewTitle == "" && dto.NewContent == "" {
		errDTO := ErrorDTO{

			Message: "nothing to update",
			Time:    time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if dto.NewTitle != "" {
		if err := h.repo.RenameNote(ctx, title, dto.NewTitle); err != nil {
			errDTO := ErrorDTO{
				Message: err.Error(),
				Time:    time.Now(),
			}

			if errors.Is(err, note.ErrNoteNotFound) {
				http.Error(w, errDTO.ToString(), http.StatusNotFound)
			} else if errors.Is(err, note.ErrNoteAlreadyExists) {
				http.Error(w, errDTO.ToString(), http.StatusConflict)
			} else {
				http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			}
			return
		}

		currentTitle = dto.NewTitle
	}

	if dto.NewContent != "" {
		if err := h.repo.ChangeContentNote(ctx, currentTitle, dto.NewContent); err != nil {
			errDTO := ErrorDTO{
				Message: err.Error(),
				Time:    time.Now(),
			}

			if errors.Is(err, note.ErrNoteNotFound) {
				http.Error(w, errDTO.ToString(), http.StatusNotFound)
			} else {
				http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			}
			return
		}
	}

	updatedNote, err := h.repo.GetNote(ctx, currentTitle)
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := json.MarshalIndent(updatedNote, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

// /notes/{title}/history
// GET
// info: pattern
func (h *HTTPHandlers) HandleGetHistoryVersionsOfNote(w http.ResponseWriter, r *http.Request)  {
	ctx := r.Context()
	title := mux.Vars(r)["title"]

	history, err := h.repo.GetNoteHistory(ctx, title)
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time: time.Now(),
		}

		if errors.Is(err, note.ErrNoteNotFound) {
				http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}

	b, err := json.MarshalIndent(history, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

// /notes/{title}/restore/{version}
// POST
// info: pattern
func (h *HTTPHandlers) HandleRestoreVersion(w http.ResponseWriter, r *http.Request)  {
	ctx := r.Context()
	title := mux.Vars(r)["title"]
	versionStr := mux.Vars(r)["version"]

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	restoredTitle, err := h.repo.RestoreVersion(ctx, title, version);
	if err != nil {
		errDTO := ErrorDTO {
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, note.ErrNoteNotFound) ||
			errors.Is(err, note.ErrVersionNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}

	updatedNote, err := h.repo.GetNote(ctx, restoredTitle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(updatedNote, "", "    ")
	if err != nil {
		panic(err)
	}
	
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}