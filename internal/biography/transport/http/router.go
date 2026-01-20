package http

import (
	"context"
	"net/http"
	"strings"

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

type contextKey string

const paramsContextKey = contextKey("path_params")

// route defines a single route with its method, pattern, and handler.
type route struct {
	method  string
	pattern []string
	handler http.HandlerFunc
}

// Router holds the list of routes and dispatches requests.
type Router struct {
	routes []route
}

// AppRepositories is an interface that groups all the repository interfaces needed by the use cases.
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

// NewRouter creates and initializes a new Router.
func NewRouter(repo AppRepositories, genID domain.IDGenerator) *Router {
	// Command Handlers
	createLifeAreaUseCase := create_life_area.NewCreateLifeAreaHandler(repo, genID)
	changeLifeAreaGoalUseCase := change_life_area_goal.NewChangeLifeAreaGoalHandler(repo)
	changeLifeAreaParentUseCase := change_life_area_parent.NewChangeLifeAreaParentHandler(repo)
	createCriteriaUseCase := create_criteria.NewCreateCriteriaHandler(repo, genID)
	deleteCriteriaUseCase := delete_criteria.NewDeleteCriteriaHandler(repo)
	deleteLifeAreaUseCase := delete_life_area.NewDeleteLifeAreaHandler(repo)

	// Query Handlers
	getLifeAreaUseCase := get_life_area.NewGetLifeAreaHandler(repo)
	listLifeAreasUseCase := list_live_area.NewListLifeAreaHandler(repo, genID)

	// HTTP Handlers
	createLifeAreaHandler := NewCreateLifeAreaHandlerHTTP(createLifeAreaUseCase)
	listLifeAreasHandler := NewListLifeAreasHandlerHTTP(listLifeAreasUseCase)
	getLifeAreaHandler := NewGetLifeAreaHandlerHTTP(getLifeAreaUseCase)
	changeLifeAreaGoalHandler := NewChangeLifeAreaGoalHandlerHTTP(changeLifeAreaGoalUseCase)
	changeLifeAreaParentHandler := NewChangeLifeAreaParentHandlerHTTP(changeLifeAreaParentUseCase)
	createCriteriaHandler := NewCreateCriteriaHandlerHTTP(createCriteriaUseCase)
	deleteCriteriaHandler := NewDeleteCriteriaHandlerHTTP(deleteCriteriaUseCase)
	deleteLifeAreaHandler := NewDeleteLifeAreaHandlerHTTP(deleteLifeAreaUseCase)

	routes := []route{
		{http.MethodGet, []string{"life-areas"}, listLifeAreasHandler.Handle},
		{http.MethodPost, []string{"life-areas"}, createLifeAreaHandler.Handle},
		{http.MethodGet, []string{"life-areas", ":id"}, getLifeAreaHandler.Handle},
		{http.MethodDelete, []string{"life-areas", ":id"}, deleteLifeAreaHandler.Handle},
		{http.MethodPatch, []string{"life-areas", ":id", "goal"}, changeLifeAreaGoalHandler.Handle},
		{http.MethodPatch, []string{"life-areas", ":id", "parent"}, changeLifeAreaParentHandler.Handle},
		{http.MethodPost, []string{"life-areas", ":id", "criteria"}, createCriteriaHandler.Handle},
		{http.MethodDelete, []string{"criteria"}, deleteCriteriaHandler.Handle},
	}
	return &Router{routes: routes}
}

// ServeHTTP finds the matching route and calls its handler.
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Handle the case of the root path "/"
	if r.URL.Path == "/" {
		http.NotFound(w, r)
		return
	}
	cleanedPath := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(cleanedPath, "/")

	for _, route := range rt.routes {
		if route.method != r.Method {
			continue
		}

		if len(route.pattern) != len(parts) {
			continue
		}

		match := true
		params := make(map[string]string)
		for i, patternPart := range route.pattern {
			if strings.HasPrefix(patternPart, ":") {
				paramName := strings.TrimPrefix(patternPart, ":")
				params[paramName] = parts[i]
				continue
			}
			if patternPart != parts[i] {
				match = false
				break
			}
		}

		if match {
			if len(params) > 0 {
				ctx := context.WithValue(r.Context(), paramsContextKey, params)
				r = r.WithContext(ctx)
			}
			route.handler(w, r)
			return
		}
	}

	http.NotFound(w, r)
}
