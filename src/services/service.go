package services

import (
	"errors"
)

var ErrTemplateNotFound = errors.New("template not found")

type ServiceInfo struct {
	Name string
	Note string
}

// StaticService временно хранит шаблоны в памяти.
// На следующем этапе его можно заменить реализацией, которая читает Portainer API.
// type StaticService struct {
// 	templates []StackTemplate
// }

// func NewStaticService() *StaticService {
// 	return &StaticService{templates: defaultTemplates()}
// }

// func (s *StaticService) List(_ context.Context) ([]StackTemplate, error) {
// 	items := make([]StackTemplate, len(s.templates))
// 	copy(items, s.templates)
// 	return items, nil
// }

// func (s *StaticService) GetByID(_ context.Context, id string) (StackTemplate, error) {
// 	for _, item := range s.templates {
// 		if item.ID == id {
// 			return item, nil
// 		}
// 	}

// 	return StackTemplate{}, ErrTemplateNotFound
// }

// func defaultTemplates() []StackTemplate {
// 	return []StackTemplate{
// 		{
// 			ID:          "go-api-starter",
// 			Name:        "Go API Starter",
// 			Category:    "Backend",
// 			Status:      "Готов к запуску",
// 			Description: "Типовой backend-стенд для разработки API-сервиса. Шаблон развернет приложение, PostgreSQL и Redis.",
// 			Purpose:     "Быстрый запуск изолированного dev-стенда для backend-разработки.",
// 			Fit:         "Подходит для тестирования API, миграций и интеграций.",
// 			Parameters:  []string{"имя стенда", "namespace/project", "branch/tag", "endpoint"},
// 			Services: []ServiceInfo{
// 				{Name: "go-api", Note: "HTTP service · порт 1323"},
// 				{Name: "postgres", Note: "Database · internal"},
// 				{Name: "redis", Note: "Cache · internal"},
// 			},
// 		},
// 		{
// 			ID:          "web-app-stack",
// 			Name:        "Web App Stack",
// 			Category:    "Frontend",
// 			Status:      "Готов к запуску",
// 			Description: "Стенд для frontend-приложения с nginx, preview-доменом и cache-сервисом.",
// 			Purpose:     "Проверка UI, статических assets и интеграции с backend API.",
// 			Fit:         "Подходит для review веток и демонстрации интерфейса внутри dev-среды.",
// 			Parameters:  []string{"имя стенда", "repository", "branch/tag", "preview domain"},
// 			Services: []ServiceInfo{
// 				{Name: "web-app", Note: "Frontend bundle · nginx"},
// 				{Name: "nginx", Note: "Reverse proxy · public"},
// 				{Name: "redis", Note: "Cache · internal"},
// 			},
// 		},
// 		{
// 			ID:          "worker-pipeline",
// 			Name:        "Worker Pipeline",
// 			Category:    "Jobs",
// 			Status:      "Готов к запуску",
// 			Description: "Шаблон для фоновых обработчиков, очереди задач и PostgreSQL.",
// 			Purpose:     "Изолированная проверка jobs, очередей и миграций без влияния на общий dev-контур.",
// 			Fit:         "Подходит для воркеров, scheduled jobs и асинхронных интеграций.",
// 			Parameters:  []string{"имя стенда", "worker image", "branch/tag", "queue name"},
// 			Services: []ServiceInfo{
// 				{Name: "worker", Note: "Background service"},
// 				{Name: "postgres", Note: "Database · internal"},
// 				{Name: "queue", Note: "Message broker · internal"},
// 			},
// 		},
// 	}
// }
