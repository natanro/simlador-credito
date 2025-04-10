package rest

import (
	"encoding/json"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/natanro/simlador-credito/motor-simulacao/entity"
	"github.com/natanro/simlador-credito/motor-simulacao/interactor"
)

type motorHandler struct {
	simulationRegister interactor.SimulationRegister
}

func NewMotorHandler(simulationRegister interactor.SimulationRegister) ServerInterface {
	return &motorHandler{
		simulationRegister: simulationRegister,
	}
}

func (h *motorHandler) CreateSimulation(ctx echo.Context) error {
	json_map := make(map[string]interface{})
	err := json.NewDecoder(ctx.Request().Body).Decode(&json_map)
	if err != nil {
		return err
	}

	request, err := h.validateRequest(json_map)
	if err != nil {
		return err
	}

	return h.simulationRegister.Create(&entity.Simulation{
		RequestedAmount: request.Amount,
		Installments:    request.Installments,
		Age:             request.Age,
	})
}

func (h *motorHandler) validateRequest(json_map map[string]interface{}) (SimulationRequest, error) {
	request := SimulationRequest{
		Amount:       json_map["amount"].(float64),
		Installments: int(json_map["installments"].(float64)),
		Age:          int(json_map["age"].(float64)),
	}

	if request.Amount <= 0 || request.Installments <= 0 || request.Age <= 0 {
		return SimulationRequest{}, errors.New("invalid request")
	}

	return request, nil
}
