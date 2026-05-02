CREATE TABLE IF NOT EXISTS notes (
	version INT NOT NULL,
	title TEXT UNIQUE NOT NULL,
	content TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS note_versions (
	version INT NOT NULL,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	changed_at TIMESTAMP NOT NULL
);