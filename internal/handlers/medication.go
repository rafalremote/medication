package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"medication/internal/middleware"
	"medication/internal/models"
	"medication/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// @title Medication API
// @version 1.0
// @description RESTful API for managing medications.
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /medications

type MedicationHandler struct {
	Service *services.MedicationService
	Logger  *logrus.Logger
}

// RegisterMedicationRoutes registers all the routes for medication-related operations.
// @Tags medications
func RegisterMedicationRoutes(r chi.Router, service *services.MedicationService, logger *logrus.Logger) {
	handler := &MedicationHandler{
		Service: service,
		Logger:  logger,
	}

	env := os.Getenv("TARGET_RELEASE")
	if env == "" {
		env = "DEV"
	}

	// Apply authentication middleware only if TARGET_RELEASE is PROD
	if env == "PROD" {
		r.Use(middleware.Authenticate)
		logger.Info("Authentication middleware enabled for PROD")
	} else {
		logger.Warn("Authentication middleware disabled for non-PROD environment")
	}

	r.Get("/", handler.GetAllMedications)       // @Router / [get]
	r.Post("/", handler.CreateMedication)       // @Router / [post]
	r.Get("/{id}", handler.GetMedication)       // @Router /{id} [get]
	r.Put("/{id}", handler.UpdateMedication)    // @Router /{id} [put]
	r.Delete("/{id}", handler.DeleteMedication) // @Router /{id} [delete]
}

// GetAllMedications handles the GET /medications endpoint to retrieve all medications.
// @Summary Get all medications
// @Description Retrieve all medications from the database with pagination support.
// @Tags medications
// @Accept json
// @Produce json
// @Param limit query int false "Number of records to fetch"
// @Param offset query int false "Offset to start fetching records"
// @Success 200 {array} models.Medication
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router / [get]
func (h *MedicationHandler) GetAllMedications(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Fetching all medications with pagination")

	// Parse query parameters for pagination
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	medications, err := h.Service.GetMedicationsWithPagination(limit, offset)
	if err != nil {
		h.Logger.Errorf("Error fetching medications: %v", err)
		http.Error(w, "Failed to fetch medications", http.StatusInternalServerError)
		return
	}

	if medications == nil {
		medications = []models.Medication{}
	}

	h.Logger.Infof("Fetched %d medications with offset %d and limit %d", len(medications), offset, limit)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(medications)
}

// CreateMedication handles the POST /medications endpoint to create a new medication.
// @Summary Create a new medication
// @Description Add a new medication to the database.
// @Tags medications
// @Accept json
// @Produce json
// @Param medication body models.Medication true "Medication object"
// @Success 201 {object} models.Medication
// @Failure 400 {object} string "Invalid request payload"
// @Failure 401 {object} string "Unauthorized"
// @Router / [post]
func (h *MedicationHandler) CreateMedication(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Creating a new medication")

	var medication models.Medication
	if err := json.NewDecoder(r.Body).Decode(&medication); err != nil {
		h.Logger.Errorf("Invalid request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateMedication(&medication); err != nil {
		h.Logger.Errorf("Error creating medication: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Created new medication: %v", medication)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(medication)
}

// GetMedication handles the GET /medications/{id} endpoint to retrieve a specific medication by ID.
// @Summary Get medication by ID
// @Description Retrieve a medication by its ID.
// @Tags medications
// @Accept json
// @Produce json
// @Param id path int true "Medication ID"
// @Success 200 {object} models.Medication
// @Failure 400 {object} string "Invalid medication ID"
// @Failure 404 {object} string "Medication not found"
// @Failure 401 {object} string "Unauthorized"
// @Router /{id} [get]
func (h *MedicationHandler) GetMedication(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Fetching medication by ID")

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.Logger.Errorf("Invalid medication ID: %v", err)
		http.Error(w, "Invalid medication ID", http.StatusBadRequest)
		return
	}

	medication, err := h.Service.GetMedicationByID(id)
	if err != nil {
		h.Logger.Errorf("Error fetching medication: %v", err)
		http.Error(w, "Medication not found", http.StatusNotFound)
		return
	}

	h.Logger.Infof("Fetched medication with ID %d", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(medication)
}

// UpdateMedication handles the PUT /medications/{id} endpoint to update a medication.
// @Summary Update medication by ID
// @Description Update the `name`, `dosage`, and `form` of a medication entry in the database.
// @Tags medications
// @Accept json
// @Produce json
// @Param id path int true "Medication ID"
// @Param medication body models.Medication true "Updated Medication object (only `name`, `dosage`, and `form` fields are allowed)"
// @Success 200 {object} map[string]string "success message"
// @Failure 400 {object} string "Invalid medication ID"
// @Failure 404 {object} string "Medication not found"
// @Failure 401 {object} string "Unauthorized"
// @Router /{id} [put]
func (h *MedicationHandler) UpdateMedication(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Updating medication by ID")

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.Logger.Errorf("Invalid medication ID: %v", err)
		http.Error(w, "Invalid medication ID", http.StatusBadRequest)
		return
	}

	var medication models.Medication
	if err := json.NewDecoder(r.Body).Decode(&medication); err != nil {
		h.Logger.Errorf("Invalid request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateMedication(id, &medication); err != nil {
		h.Logger.Errorf("Error updating medication: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Updated medication with ID %d", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "successfully updated"})
}

// DeleteMedication handles the DELETE /medications/{id} endpoint to delete a medication.
// @Summary Delete medication by ID
// @Description Delete a medication entry from the database.
// @Tags medications
// @Accept json
// @Produce json
// @Param id path int true "Medication ID"
// @Success 200 {object} map[string]string "success message"
// @Failure 400 {object} string "Invalid medication ID"
// @Failure 404 {object} string "Medication not found"
// @Failure 401 {object} string "Unauthorized"
// @Router /{id} [delete]
func (h *MedicationHandler) DeleteMedication(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Deleting medication by ID")

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.Logger.Errorf("Invalid medication ID: %v", err)
		http.Error(w, "Invalid medication ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteMedication(id); err != nil {
		h.Logger.Errorf("Error deleting medication: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Deleted medication with ID %d", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "successfully deleted"})
}
