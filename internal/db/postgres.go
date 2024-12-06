package db

import (
	"database/sql"
	"fmt"
	"medication/config"
	"medication/internal/models"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	conn *sql.DB
}

func NewPostgresDB(cfg *config.Config) (*PostgresDB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	conn, err := sql.Open(cfg.DBDriver, dsn)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Connected to PostgresDB at %s:%d\n", cfg.DBHost, cfg.DBPort)
	return &PostgresDB{conn: conn}, nil
}

func (db *PostgresDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.conn.Exec(query, args...)
}

func (db *PostgresDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.conn.Query(query, args...)
}

func (db *PostgresDB) GetMedications() ([]models.Medication, error) {
	rows, err := db.conn.Query("SELECT id, name, dosage, form, created_at, updated_at FROM medications")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var medications []models.Medication
	for rows.Next() {
		var med models.Medication
		if err := rows.Scan(&med.ID, &med.Name, &med.Dosage, &med.Form, &med.CreatedAt, &med.UpdatedAt); err != nil {
			return nil, err
		}
		medications = append(medications, med)
	}
	return medications, nil
}

func (db *PostgresDB) GetMedicationsWithPagination(limit, offset int) ([]models.Medication, error) {
	rows, err := db.conn.Query(
		"SELECT id, name, dosage, form, created_at, updated_at FROM medications ORDER BY id LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var medications []models.Medication
	for rows.Next() {
		var med models.Medication
		if err := rows.Scan(&med.ID, &med.Name, &med.Dosage, &med.Form, &med.CreatedAt, &med.UpdatedAt); err != nil {
			return nil, err
		}
		medications = append(medications, med)
	}
	return medications, nil
}

func (db *PostgresDB) GetMedicationByID(id int) (*models.Medication, error) {
	var med models.Medication
	err := db.conn.QueryRow(
		"SELECT id, name, dosage, form, created_at, updated_at FROM medications WHERE id = $1",
		id,
	).Scan(&med.ID, &med.Name, &med.Dosage, &med.Form, &med.CreatedAt, &med.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &med, nil
}

func (db *PostgresDB) CreateMedication(medication *models.Medication) error {
	err := db.conn.QueryRow(
		`INSERT INTO medications (name, dosage, form) 
		 VALUES ($1, $2, $3) 
		 RETURNING id, created_at, updated_at`,
		medication.Name, medication.Dosage, medication.Form,
	).Scan(&medication.ID, &medication.CreatedAt, &medication.UpdatedAt)
	return err
}

func (db *PostgresDB) UpdateMedication(id int, medication *models.Medication) error {
	_, err := db.conn.Exec(
		"UPDATE medications SET name = $1, dosage = $2, form = $3, updated_at = NOW() WHERE id = $4",
		medication.Name, medication.Dosage, medication.Form, id,
	)
	return err
}

func (db *PostgresDB) DeleteMedication(id int) error {
	_, err := db.conn.Exec("DELETE FROM medications WHERE id = $1", id)
	return err
}

func (db *PostgresDB) Close() error {
	return db.conn.Close()
}
