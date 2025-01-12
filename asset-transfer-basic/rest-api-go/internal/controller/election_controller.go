package controller

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"rest-api-go/internal/models"
	"rest-api-go/internal/service"
)

type RegisterCandidatesRequest struct {
	ElectionID uint               `json:"electionId"`
	Candidates []models.Candidate `json:"candidates"`
}

type ElectionController struct {
	electionService service.ElectionService
}

func NewElectionController(service service.ElectionService) *ElectionController {
	return &ElectionController{electionService: service}
}

func (ctrl *ElectionController) RegisterElection(ctx *fiber.Ctx) error {
	var election *models.Election

	err := ctx.BodyParser(&election)
	if err != nil {
		log.Printf("Invalid request format")
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	err = ctrl.electionService.RegisterElection(election)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	log.Printf("register-election request successful")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Election successfully created"})
}

func (ctrl *ElectionController) GetAllElections(ctx *fiber.Ctx) error {
	elections, err := ctrl.electionService.GetAllElections()
	if err != nil {
		log.Printf("Failed to fetch elections: %v", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	log.Printf("get all elections request successful")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message":   "Elections retrieved successfully",
		"elections": elections,
	})
}

func (ctrl *ElectionController) GetElectionById(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid election ID",
			"error":   err.Error(),
		})
	}
	election, err := ctrl.electionService.GetElectionById(uint(id))
	if err != nil {
		log.Printf("Failed to fetch election by id: %v", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch election",
			"error":   err.Error(),
		})
	}

	if election == nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Election not found",
		})
	}

	log.Printf("get election by id request successful for ID: %d", id)
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message":  "Election fetched successfully",
		"election": election,
	})
}

func (ctrl *ElectionController) RegisterCandidates(ctx *fiber.Ctx) error {
	var req RegisterCandidatesRequest

	err := ctx.BodyParser(&req)
	if err != nil {
		log.Printf("Invalid request body")
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := ctrl.electionService.RegisterCandidates(req.ElectionID, req.Candidates); err != nil {
		log.Printf("Failed to register candidates %v", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to register candidates",
			"error":   err.Error(),
		})
	}
	return nil
}

func (ctrl *ElectionController) GetCandidatesByElectionId(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("electionId")
	if err != nil {
		log.Printf("invalid request body")
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	candidates, err := ctrl.electionService.GetCandidates(uint(id))
	if err != nil {
		log.Printf("can't get candidates for specific election")
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	log.Printf("get election by id request successful for ID: %d", id)
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message":    "Election candidates fetched successfully",
		"candidates": candidates,
	})
}
