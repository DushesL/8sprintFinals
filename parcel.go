package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	query := `INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)`
	result, err := s.db.Exec(query, p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, fmt.Errorf("failed to add parcel: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка

	// заполните объект Parcel данными из таблицы
	query := `SELECT number, client, status, address, created_at FROM parcel WHERE number = ?`
	row := s.db.QueryRow(query, number)

	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, fmt.Errorf("parcel with number %d not found", number)
		}
		return p, fmt.Errorf("failed to scan parcel: %w", err)
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы
	query := `SELECT number, client, status, address, created_at FROM parcel WHERE client = ?`
	rows, err := s.db.Query(query, client)
	if err != nil {
		return nil, fmt.Errorf("failed to query parcels by client: %w", err)
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan parcel row: %w", err)
		}
		parcels = append(parcels, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over parcel rows: %w", err)
	}

	return parcels, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	query := `UPDATE parcel SET status = ? WHERE number = ?`
	result, err := s.db.Exec(query, status, number)
	if err != nil {
		return fmt.Errorf("failed to update parcel status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("parcel with number %d not found", number)
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	queryCheck := `SELECT status FROM parcel WHERE number = ?`
	var currentStatus string
	err := s.db.QueryRow(queryCheck, number).Scan(&currentStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("parcel with number %d not found", number)
		}
		return fmt.Errorf("failed to check parcel status: %w", err)
	}

	if currentStatus != ParcelStatusRegistered {
		return fmt.Errorf("cannot change address for parcel with status '%s'", currentStatus)
	}

	// Обновление адреса
	queryUpdate := `UPDATE parcel SET address = ? WHERE number = ?`
	result, err := s.db.Exec(queryUpdate, address, number)
	if err != nil {
		return fmt.Errorf("failed to update parcel address: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("parcel with number %d not found", number)
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	queryCheck := `SELECT status FROM parcel WHERE number = ?`
	var currentStatus string
	err := s.db.QueryRow(queryCheck, number).Scan(&currentStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("parcel with number %d not found", number)
		}
		return fmt.Errorf("failed to check parcel status: %w", err)
	}

	if currentStatus != ParcelStatusRegistered {
		return fmt.Errorf("cannot delete parcel with status '%s'", currentStatus)
	}

	// Удаление посылки
	queryDelete := `DELETE FROM parcel WHERE number = ?`
	result, err := s.db.Exec(queryDelete, number)
	if err != nil {
		return fmt.Errorf("failed to delete parcel: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("parcel with number %d not found", number)
	}

	return nil
}
