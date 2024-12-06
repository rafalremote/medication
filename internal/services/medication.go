package services

import (
	"errors"

	"medication/internal/db"
	"medication/internal/models"
)

type MedicationService struct {
	DB db.DBInterface
}

func NewMedicationService(database db.DBInterface) *MedicationService {
	return &MedicationService{DB: database}
}

func (s *MedicationService) GetMedications() ([]models.Medication, error) {
	return s.DB.GetMedications()
}

func (s *MedicationService) GetMedicationsWithPagination(limit, offset int) ([]models.Medication, error) {
	return s.DB.GetMedicationsWithPagination(limit, offset)
}

func (s *MedicationService) GetMedicationByID(id int) (*models.Medication, error) {
	medication, err := s.DB.GetMedicationByID(id)
	if err != nil {
		return nil, err
	}
	if medication == nil {
		return nil, errors.New("medication not found")
	}
	return medication, nil
}

func (s *MedicationService) CreateMedication(medication *models.Medication) error {
	if medication.Name == "" {
		return errors.New("medication name cannot be empty")
	}
	return s.DB.CreateMedication(medication)
}

// UpdateMedication updates an existing medication in the database.
func (s *MedicationService) UpdateMedication(id int, medication *models.Medication) error {
	existing, err := s.DB.GetMedicationByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("medication not found")
	}

	existing.Name = medication.Name
	existing.Dosage = medication.Dosage
	existing.Form = medication.Form

	return s.DB.UpdateMedication(id, existing)
}

func (s *MedicationService) DeleteMedication(id int) error {
	existing, err := s.DB.GetMedicationByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("medication not found")
	}
	return s.DB.DeleteMedication(id)
}
