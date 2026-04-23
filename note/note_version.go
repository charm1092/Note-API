package note

import "time"

type NoteVersion struct {
	Version int
	Title string
	Content string
	ChangedAt time.Time
}

func NewNoteVersion(note Note, version int) NoteVersion {
	return NoteVersion{
        Version:   version,
        Title:     note.Title,
        Content:   note.Content,
        ChangedAt: time.Now(),
	}
}