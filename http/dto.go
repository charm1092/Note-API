package http

import (
	"encoding/json"
	"errors"
	"time"
)

type NoteDTO struct {
	Title string `json:"title"`
	Content string `json:"content"`
}

type UpdateNoteDTO struct {
	NewTitle string `json:"new_title"`
	NewContent string `json:"new_content"`
}

type ErrorDTO struct {
	Message string
	Time time.Time
}

func (n *NoteDTO) ValidateForCreate() error {
	if n.Title == "" {
		return errors.New("title is empty")
	}

	return nil
}

func (e *ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}

	return string(b)
}