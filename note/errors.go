package note

import "errors"

var ErrNoteNotFound = errors.New("Note not found")
var ErrNoteAlreadyExists = errors.New("Note already exists")
var ErrVersionNotFound = errors.New("Such note's version not found")