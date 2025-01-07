package server

import (
	"denet/tasks"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) getUserStatus(w http.ResponseWriter, r *http.Request) {
	// Парсим данные из запроса
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	u, err := s.db.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func (s *Server) getLeaderboard(w http.ResponseWriter, r *http.Request) {
	// Парсим данные из запроса
	vars := mux.Vars(r)
	count, _ := strconv.Atoi(vars["count"])
	if count <= 0 {
		count = 10
	}

	users, err := s.db.GetLeaderboard(r.Context(), count)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	resp, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (s *Server) completeTask(w http.ResponseWriter, r *http.Request) {
	// Парсим данные из запроса
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Парсим тело запроса
	var requestBody struct {
		TelegramUsername int `json:"telegram_user_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// 38.180.148.174:8085
	// Проверяем подписку пользователя на канал
	subscribed, err := tasks.IsUserSubscribedToChannel(requestBody.TelegramUsername)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to check subscription: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !subscribed {
		// Не уверен по поводу корректности статус кода по RFC в таком случае
		http.Error(w, "User is not subscribed to the channel", http.StatusForbidden)
		return
	}
	// Увеличиваем количество поинтов у пользователя
	err = s.db.IncrementPoints(r.Context(), nil, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to increment points: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task completed successfully"))
}

// Мною неясный ендпоинт - логику можно было перенести в регистрацию
func (s *Server) setReferrer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// Парсим тело запроса
	var requestBody struct {
		Referer int `json:"referrer"`
	}

	// is middleware доставать username из контекста и по нему устанавливать id ?
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if requestBody.Referer == id {
		http.Error(w, "Invalid referer ID", http.StatusBadRequest)
		return
	}

	err = s.db.SetRefferer(r.Context(), id, requestBody.Referer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Здесь должна быть логика установки реферального кода для пользователя
	w.Write([]byte("Referrer set for user " + strconv.Itoa(id)))
}
