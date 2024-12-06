package db

import "medication/internal/models"

type DBInterface interface {
	GetMedications() ([]models.Medication, error)
	GetMedicationsWithPagination(limit, offset int) ([]models.Medication, error)
	GetMedicationByID(id int) (*models.Medication, error)
	CreateMedication(medication *models.Medication) error
	UpdateMedication(id int, medication *models.Medication) error
	DeleteMedication(id int) error
	Close() error
}
