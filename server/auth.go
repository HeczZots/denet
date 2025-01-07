package server

import (
	"encoding/json"

	"log/slog"
	"net/http"
	"strconv"
)

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	// Парсим данные из запроса
	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверяем
	u, err := s.db.GetUserByLoginAndPassword(r.Context(), creds.Login, creds.Password)
	if err != nil {
		slog.Info("login", "err", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	ut, err := s.mw.GenerateJWT(creds.Login)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User ID " + strconv.Itoa(u.ID) + "\nAccess token: " + ut))
}

func (s *Server) registration(w http.ResponseWriter, r *http.Request) {
	// Парсим данные из запроса
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}
	var creds Credentials
	err = json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Проверяем, что логин и пароль не пустые
	if creds.Login == "" || creds.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Хэшируем пароль
	hashedPassword, err := hashPassword(creds.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Сохраняем пользователя в базу данных
	err = s.db.CreateUser(r.Context(), hashedPassword, creds.Login)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	ut, err := s.mw.GenerateJWT(creds.Login)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully Access token: \n" + ut))
}
