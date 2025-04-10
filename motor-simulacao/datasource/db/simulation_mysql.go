package db

import (
	"errors"
	"fmt"

	"github.com/natanro/simlador-credito/motor-simulacao/datasource/db/model"
	"gorm.io/gorm"
)

type simulationRepository struct {
	db              *gorm.DB
}

type SimulationRepository interface {
	Create(simulation *model.Simulation) (int, error)
	Update(simulation *model.Simulation) error
}

func NewSimulationRepository(db *gorm.DB) SimulationRepository {
	return &simulationRepository{
		db:              db,
	}
}

func (s *simulationRepository) Create(simulation *model.Simulation) (int, error) {
	if simulation == nil {
		return -1, errors.New("simulation cannot be nil")
	}

	if err := s.db.Create(simulation).Error; err != nil {
		return -1, fmt.Errorf("error creating simulation: %w", err)
	}

	return simulation.ID, nil
}

func (s *simulationRepository) Update(simulation *model.Simulation) error {
	if simulation == nil {
		return errors.New("simulation cannot be nil")
	}

	return s.db.Save(simulation).Error
}
