package http

import (
	"github.com/gin-gonic/gin"

	"go-interview/internal/biography/app/commands/change_life_area_goal"
	"go-interview/internal/biography/app/commands/change_life_area_parent"
	"go-interview/internal/biography/app/commands/create_criteria"
	"go-interview/internal/biography/app/commands/create_life_area"
	"go-interview/internal/biography/app/commands/delete_criteria"
	"go-interview/internal/biography/app/commands/delete_life_area"
	"go-interview/internal/biography/app/queries/get_life_area"
	list_live_area "go-interview/internal/biography/app/queries/list_life_areas"
	"go-interview/internal/biography/domain"
)

// AppRepositories groups repositories used by handlers.
type AppRepositories interface {
	domain.LifeAreaCreator
	domain.LifeAreaGetter
	domain.LifeAreaLister
	domain.LifeAreaDeleter
	domain.LifeAreaParentChanger
	domain.LifeAreaGoalChanger
	domain.CriteriaDeleter
	domain.CriteriaCreator
	domain.CriteriaNodeGetter
}

// RegisterRoutes attaches all biography HTTP endpoints to the provided Gin router.
func RegisterRoutes(router gin.IRoutes, repo AppRepositories, genID domain.IDGenerator) {
	createLifeAreaUseCase := create_life_area.NewCreateLifeAreaHandler(repo, genID)
	changeLifeAreaGoalUseCase := change_life_area_goal.NewChangeLifeAreaGoalHandler(repo)
	changeLifeAreaParentUseCase := change_life_area_parent.NewChangeLifeAreaParentHandler(repo)
	createCriteriaUseCase := create_criteria.NewCreateCriteriaHandler(repo, genID)
	deleteCriteriaUseCase := delete_criteria.NewDeleteCriteriaHandler(repo)
	deleteLifeAreaUseCase := delete_life_area.NewDeleteLifeAreaHandler(repo)
	getLifeAreaUseCase := get_life_area.NewGetLifeAreaHandler(repo)
	listLifeAreasUseCase := list_live_area.NewListLifeAreaHandler(repo, genID)

	createLifeAreaHandler := NewCreateLifeAreaHandlerHTTP(createLifeAreaUseCase)
	listLifeAreasHandler := NewListLifeAreasHandlerHTTP(listLifeAreasUseCase)
	getLifeAreaHandler := NewGetLifeAreaHandlerHTTP(getLifeAreaUseCase)
	changeLifeAreaGoalHandler := NewChangeLifeAreaGoalHandlerHTTP(changeLifeAreaGoalUseCase)
	changeLifeAreaParentHandler := NewChangeLifeAreaParentHandlerHTTP(changeLifeAreaParentUseCase)
	createCriteriaHandler := NewCreateCriteriaHandlerHTTP(createCriteriaUseCase)
	deleteCriteriaHandler := NewDeleteCriteriaHandlerHTTP(deleteCriteriaUseCase)
	deleteLifeAreaHandler := NewDeleteLifeAreaHandlerHTTP(deleteLifeAreaUseCase)

	router.GET("/life-areas", listLifeAreasHandler.Handle)
	router.POST("/life-areas", createLifeAreaHandler.Handle)
	router.GET("/life-areas/:id", getLifeAreaHandler.Handle)
	router.DELETE("/life-areas/:id", deleteLifeAreaHandler.Handle)
	router.PATCH("/life-areas/:id/goal", changeLifeAreaGoalHandler.Handle)
	router.PATCH("/life-areas/:id/parent", changeLifeAreaParentHandler.Handle)
	router.POST("/life-areas/:id/criteria", createCriteriaHandler.Handle)
	router.DELETE("/criteria", deleteCriteriaHandler.Handle)
}
