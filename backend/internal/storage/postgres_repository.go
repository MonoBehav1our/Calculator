package storage

import (
	"Calculator/internal/model"
	"database/sql"
	"errors"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{
		db: db,
	}
}

func (s *PostgresStorage) GetAllCalculations() ([]model.ExpressionResult, error) {
	rows, err := s.db.Query("SELECT id, expression, result FROM calculations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.ExpressionResult
	for rows.Next() {
		var result model.ExpressionResult
		if err := rows.Scan(&result.Id, &result.Expression, &result.Result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *PostgresStorage) GetCalculation(id int) (model.ExpressionResult, error) {
	var result model.ExpressionResult

	err := s.db.
		QueryRow("SELECT id, expression, result FROM calculations WHERE id = $1", id).
		Scan(&result.Id, &result.Expression, &result.Result)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ExpressionResult{}, errors.New("calculation not found")
		}
		return model.ExpressionResult{}, err
	}

	return result, nil
}

func (s *PostgresStorage) CreateCalculation(expression string, result string) (int, error) {
	var id int

	err := s.db.
		QueryRow("INSERT INTO calculations (expression, result) VALUES ($1, $2) RETURNING id",
			expression, result,
		).Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s *PostgresStorage) UpdateCalculation(m model.ExpressionResult) error {
	r, err := s.db.Exec("UPDATE calculations SET expression=$1, result=$2 WHERE id=$3", m.Expression, m.Result, m.Id)
	if err != nil {
		return err
	}

	affected, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (s *PostgresStorage) DeleteCalculation(id int) error {
	r, err := s.db.Exec("DELETE FROM calculations WHERE id=$1", id)
	if err != nil {
		return err
	}

	affected, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
