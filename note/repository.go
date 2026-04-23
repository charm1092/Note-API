package note

import (
	"sync"
	"time"
)

// * создать заметку
// * получить заметку
// * изменить заметку
// * посмотреть историю
// * восстановить старую версию
// * удалить заметку
type Repository struct {
	notes map[string]*Note
	versions map[string][]NoteVersion
	mtx sync.RWMutex
}

func NewList() *Repository {
	return &Repository{
		notes: make(map[string]*Note),
		versions: make(map[string][]NoteVersion),
	}
}

func (r *Repository) AddNote(note Note) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.notes[note.Title]; ok {
		return ErrNoteAlreadyExists
	}

	r.notes[note.Title] = &note
	return nil
}

func (r *Repository) GetNote(title string) (Note, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	note, ok := r.notes[title]
	if !ok {
		return Note{}, ErrNoteNotFound
	}

	return *note, nil
}

func (r *Repository) ListNotes() map[string]*Note {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	tmp := make(map[string]*Note, len(r.notes))

	for k, v := range r.notes {
		tmp[k] = v
	}
	return tmp
}

func (r *Repository) DeleteNote(title string) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.notes[title]; !ok {
		return ErrNoteNotFound
	}

	delete(r.notes, title)

	return nil
}

func (r *Repository) RenameNote(title string, newTitle string) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	note, ok := r.notes[title]
	if !ok {
		return ErrNoteNotFound
	}

	if _, ok := r.notes[newTitle]; ok {
		return ErrNoteAlreadyExists
	}

	note.Title = newTitle
	r.notes[newTitle] = note
	delete(r.notes, title)

	return nil
}

func (r *Repository) ChangeContentNote(title string, newContent string) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.notes[title]; !ok {
		return ErrNoteNotFound
	}

	note := r.notes[title]
	newNote := NewNoteVersion(*note, len(r.versions[title])+1)

	r.versions[title] = append(r.versions[title], newNote)
	note.Content = newContent
	note.UpdatedAt = time.Now()

	return nil
}

func (r *Repository) GetNoteHistory(title string) ([]NoteVersion, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if _, ok := r.notes[title]; !ok {
		return nil, ErrNoteNotFound
	}

	versions := r.versions[title]
	tmp := make([]NoteVersion, len(versions))
	copy(tmp, versions)

	return tmp, nil
}

func (r *Repository) RestoreVersion(title string, version int) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	note, ok := r.notes[title]
	if !ok {
		return ErrNoteNotFound
	}

	versions := r.versions[title]
	for _, v := range versions {
		if v.Version == version {
			currentVersion := NewNoteVersion(*note, len(versions)+1)
			r.versions[title] = append(r.versions[title], currentVersion)

			note.Content = v.Content

			note.UpdatedAt = time.Now()
			return nil
		}
	}

	return ErrVersionNotFound
	
}