package handlers

import (
	"encoding/json"
	"net/http"

	"awesomeProject/internal/app/counter"
	"awesomeProject/internal/config"
	"github.com/gorilla/mux"
)

type Handlers struct {
	cfg     *config.Config
	counter counter.PageViewCounter
}

func New(cfg *config.Config, counter counter.PageViewCounter) *Handlers {
	return &Handlers{cfg: cfg, counter: counter}
}

// POST /increment/{page_id}: увеличивает счетчик просмотров для страницы с идентификатором page_id на 1.
func (s *Handlers) IncrementReq(w http.ResponseWriter, r *http.Request) {
	// Валидация

	// Обработка
	vars := mux.Vars(r)
	pageID := vars["pageID"]
	if err := s.counter.Increment(pageID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

	// Ответ
}

// GET /count/{page_id}: возвращает текущее количество просмотров для страницы page_id в формате JSON { "page_id": string, "count": int }.
func (s *Handlers) GetCountReq(w http.ResponseWriter, r *http.Request) {
	// Валидация

	// Обработка
	vars := mux.Vars(r)
	pageID := vars["pageID"]
	count, err := s.counter.GetCount(pageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	response := map[string]interface{}{
		"page_id": pageID,
		"count":   count,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	// Ответ
}
