package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"myapp/internal/app/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{svc: s}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	users := h.svc.ListUsers()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(users)
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name,omitempty"`
}

type registerResponse struct {
    ID       int64  `json:"id"`
    Success  bool   `json:"success"`
    Message  string `json:"message"`
    Token    string `json:"token"`
    UserData struct {
        ID       int64  `json:"id"`
        Username string `json:"username"`
        Email    string `json:"email"`
        Name     string `json:"name"`
    } `json:"user_data"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var req registerRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		log.Printf("register: invalid JSON: %v", err)
		http.Error(w, "Bad Request: invalid JSON", http.StatusBadRequest)
		return
	}
	u, err := h.svc.RegisterUser(service.RegisterParams{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
    var resp registerResponse
    resp.ID = u.ID
    resp.Success = true
    resp.Message = "Login berhasil"
    resp.Token = "123" // optional token, can be filled with JWT later
    resp.UserData.ID = u.ID
    resp.UserData.Username = u.Username
	resp.UserData.Email = u.Email
	resp.UserData.Name = u.Name
	_ = json.NewEncoder(w).Encode(resp)
}
