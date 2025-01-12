package router

import (
	"github.com/gofiber/fiber/v2"
	"rest-api-go/internal/controller"
)

func RegisterElectionRoutes(r *fiber.App, electionCtrl *controller.ElectionController) {
	route := r.Group("/elections")
	route.Post("/create", electionCtrl.RegisterElection)
	route.Get("/", electionCtrl.GetAllElections)
	route.Get("/:id", electionCtrl.GetElectionById)
	route.Post("/candidates", electionCtrl.RegisterCandidates)
	route.Get("/candidates/:electionId", electionCtrl.GetCandidatesByElectionId)
}
