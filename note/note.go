package note

import "time"

type Note struct {
	Version int
	Title string
	Content string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func NewNote(title string, content string) Note {
	return Note{
		Version: 0,
		Title: title,
		Content: content,

		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}
}
