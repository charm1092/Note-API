package note

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) *Repository {
	return &Repository{
		pool: conn,
	}
}

func (r *Repository) AddNote(ctx context.Context, note Note) error {

	sqlQuery := `
	INSERT INTO notes(version, title, content, created_at)
	VALUES($1, $2, $3, $4)
	ON CONFLICT (title) DO NOTHING
	`

	_, err := r.pool.Exec(ctx, sqlQuery, note.Version, note.Title, note.Content, note.CreatedAt)
	return err
}

func (r *Repository) GetNote(ctx context.Context, title string) (Note, error) {
	sqlQuery := `
	SELECT version, title, content, created_at, updated_at
	FROM notes
	WHERE title = $1
	`

	var note Note

	err := r.pool.QueryRow(ctx, sqlQuery, title).Scan(
		&note.Version,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Note{}, ErrNoteNotFound
		}
		return Note{}, err
	}
	return note, nil
}

func (r *Repository) ListNotes(ctx context.Context) (map[string]*Note, error) {
	sqlQuery := `
	SELECT version, title, content, created_at, updated_at
	FROM notes
	ORDER BY created_at DESC 
	`
	rows, err := r.pool.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make(map[string]*Note, 0)

	for rows.Next() {
		var note Note

		err := rows.Scan(
			&note.Version,
			&note.Title,
			&note.Content,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		n := note
		notes[n.Title] = &n
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *Repository) DeleteNote(ctx context.Context, title string) error {
	sqlQuery := `
	DELETE FROM notes
	WHERE title = $1
	`
	_, err := r.pool.Exec(ctx,
		sqlQuery,
		title,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrNoteNotFound
		}
		return err
	}
	return nil
}

func (r *Repository) RenameNote(ctx context.Context, title string, newTitle string) error {
	sqlQueryForVersionTable := `
	INSERT INTO note_versions(version, title, content, changed_at)
	VALUES($1, $2, $3, $4)
	`

	sqlQueryForUpdateNotes := `
	UPDATE notes SET
	version = version+1,
	title = $1,
	updated_at = $2
	WHERE title = $3 
	`

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	currentNote, err := r.GetNote(ctx, title)
	if err != nil {
		return err
	}

	version := NewNoteVersion(currentNote, currentNote.Version+1)
	currentTime := time.Now()

	_, err = tx.Exec(ctx,
		sqlQueryForVersionTable,
		version.Version,
		version.Title,
		version.Content,
		version.ChangedAt,
	)
	if err != nil {
		return err
	}


	commandTag, err := tx.Exec(ctx,
		 sqlQueryForUpdateNotes,
		 newTitle, currentTime, title,
		)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ErrNoteNotFound
	}

	return tx.Commit(ctx)
	
}

func (r *Repository) ChangeContentNote(ctx context.Context, title string, newContent string) error {
	sqlQueryForVersionTable := `
	INSERT INTO note_versions(version, title, content, changed_at)
	VALUES($1, $2, $3, $4)
	`

	sqlQueryForUpdateNotes := `
	UPDATE notes SET
	version = version+1,
	content = $1,
	updated_at = $2
	WHERE title = $3 
	`

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	currentNote, err := r.GetNote(ctx, title)
	if err != nil {
		return ErrNoteNotFound
	}

	version := NewNoteVersion(currentNote, currentNote.Version+1)
	_, err = tx.Exec(ctx, 
		sqlQueryForVersionTable,
		version.Version,
		version.Title,
		version.Content,
		version.ChangedAt,
	)

	if err != nil {
		return err
	}

	currentTime := time.Now()

	commandTag, err := tx.Exec(ctx,
		 sqlQueryForUpdateNotes,
		 newContent, currentTime, title,
		)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ErrNoteNotFound
	}

	return tx.Commit(ctx)
}

func (r *Repository) GetNoteHistory(ctx context.Context, title string) ([]NoteVersion, error) {
	sqlQuery := `
	SELECT version, title, content, changed_at
	FROM note_versions
	WHERE title = $1
	ORDER BY changed_at DESC
	`
	rows, err := r.pool.Query(ctx, sqlQuery, title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	noteVersions := make([]NoteVersion, 0)

	for rows.Next() {
		var version NoteVersion

		err := rows.Scan(
			&version.Version,
			&version.Title,
			&version.Content,
			&version.ChangedAt,
		)
		if err != nil {
			return nil, err
		}

		noteVersions = append(noteVersions, version)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return noteVersions, nil

}

func (r *Repository) RestoreVersion(ctx context.Context, title string, version int) (string, error) {
	// берем нужную нам версию заметки для восстановления
	sqlQuery1 := `
	SELECT version, title, content, changed_at
	FROM note_versions
	WHERE title = $1 AND version = $2
	`

	// запрос для вставки в список версий текущей заметки
	sqlQuery2 := `
	INSERT INTO note_versions(version, title, content, changed_at)
	VALUES($1, $2, $3, $4)
	`

	// перезапись текущей заметки на ту кот получили в 1 запросе
	sqlQuery3 := `
	UPDATE notes
	SET title = $1,
		content = $2,
		version = version+1
		updated_at = $3
	WHERE title = $4
	`

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	currentNote, err := r.GetNote(ctx, title)
	if err != nil {
		return "", err
	}

	var oldVersion NoteVersion
	currentVersion := NewNoteVersion(currentNote, currentNote.Version+1)
	currentTime := time.Now()

	err = tx.QueryRow(ctx, sqlQuery1, title, version).Scan(
		&oldVersion.Title,
		&oldVersion.Version,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return "", ErrVersionNotFound
		}
		return "", err
	}

	_, err = tx.Exec(ctx, sqlQuery2,
		currentVersion.Version,
		currentVersion.Title,
		currentVersion.Content,
		currentVersion.ChangedAt,
	)
	if err != nil {
		return "", err
	}

	commandTag, err := tx.Exec(ctx, sqlQuery3,
		oldVersion.Title,
		oldVersion.Content,
		currentTime,
		title,
	)
	if err != nil {
		return "", err
	}

	if commandTag.RowsAffected() == 0 {
		return "", ErrNoteNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return oldVersion.Title, nil
}