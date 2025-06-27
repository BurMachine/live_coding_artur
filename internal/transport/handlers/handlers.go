package handlers

import (
	"awesomeProject/internal/app/counter"
	"awesomeProject/internal/config"
	"net/http"
)

type Handlers struct {
	cfg *config.Config
	cnt *counter.Counter
}

func New(cfg *config.Config, counter2 *counter.Counter) *Handlers {
	return &Handlers{cfg: cfg, cnt: counter2}
}

// POST /increment/{page_id}: увеличивает счетчик просмотров для страницы с идентификатором page_id на 1.
func (s *Handlers) IncrementReq(w http.ResponseWriter, r *http.Request) {
	// Валидация

	// Обработка
	s.cnt.Increment("1")

	// Ответ
}

// GET /count/{page_id}: возвращает текущее количество просмотров для страницы page_id в формате JSON { "page_id": string, "count": int }.
func (s *Handlers) GetCountReq(w http.ResponseWriter, r *http.Request) {
	// Валидация

	// Обработка
	s.cnt.GetCount("1")

	// Ответ
}
