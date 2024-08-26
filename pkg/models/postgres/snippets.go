package postgres

import (
	"database/sql"
	"errors"
	"snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	query := `
	INSERT INTO snippets (title, content, created_at, expires_at) 
	VALUES ($1, $2, NOW(), NOW() + $3);`

	res, err := m.DB.Exec(query, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, nil
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	query := `
	SELECT * 
	FROM snippets 
	WHERE id=$1 and expires_at > NOW();`

	s := &models.Snippet{}

	err := m.DB.QueryRow(query, id).Scan(&s.Id, &s.Title, &s.Content, &s.CreatedAt, &s.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	query := `
	SELECT *
	FROM snippets
	WHERE expires_at > NOW()
	ORDER BY created_at DESC
	FETCH FIRST 10 rows ONLY;`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}
	for rows.Next() {
		s := &models.Snippet{}

		err = rows.Scan(&s.Id, &s.Title, &s.Content, &s.CreatedAt, &s.ExpiresAt)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	return snippets, nil
}
