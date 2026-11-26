package routes

import (
	"conductor.com.br/pier/hooks/cmd/http/helpers"
	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/domains/manager"
	"github.com/go-chi/chi"
	"net/http"
)

type Management struct {
	service manager.Service
}

func NewManagementRouter(service manager.Service) Management {
	return Management{service: service}
}

func (m Management) Route(r chi.Router) {
	r.Get("/health", m.checkHealth)
	r.Get("/info", m.getInfo)
}

// @Title checkHealth
// @Tags Management
// @Summary Check health
// @Description Checks the API and it's components health
// @Success 200 {object} manager.ManagerHealth
// @Accept json
// @Router /manager/health [get]
func (m Management) checkHealth(w http.ResponseWriter, r *http.Request) {
	resp := m.service.CheckHealth()
	statusCode := http.StatusOK

	if resp.Status != "UP" {
		statusCode = http.StatusInternalServerError
	}
	helpers.JsonResponse(w, resp, statusCode)
}

// @Title checkInfo
// @Tags Management, Internals
// @Summary Check version information
// @Description Checks the API and it's components versioning information
// @Success 200 {object} manager.ManagerInfoResponse
// @Accept json
// @Router /manager/info [get]
func (m Management) getInfo(w http.ResponseWriter, r *http.Request) {
	var id, author, email, time string

	if id == "" {

		id, author, email, time = manager.GetGitInfo()
	}

	git := manager.ManagerInfoResponse{
		App: &manager.ManagerInfoResponseApp{
			Profile:     global.Viper.GetString("app.profile"),
			Name:        global.Viper.GetString("app.name"),
			Description: global.Viper.GetString("app.description"),
			Version:     global.VersionInfo.Current,
		},
		Git: &manager.ManagerInfoResponseGitCommit{
			Id:     id,
			Author: author,
			Email:  email,
			Time:   time,
		},
	}
	helpers.JsonResponse(w, git, http.StatusOK)
}
