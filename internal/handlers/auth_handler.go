package handlers

import (
	"encoding/json"
	"net/http"
	"log"
	"filestore/internal/auth"
	"filestore/internal/db"
	"filestore/internal/utils"
)

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user db.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}
	user.Password = hashedPassword

	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	err = db.GetDB().QueryRow(query, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, user)
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user db.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var storedUser db.User
	query := `SELECT id, password FROM users WHERE email = $1`
	err = db.GetDB().QueryRow(query, user.Email).Scan(&storedUser.ID, &storedUser.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if !auth.CheckPasswordHash(user.Password, storedUser.Password) {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := auth.GenerateJWT(storedUser.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	response := map[string]string{"token": token}
	utils.RespondWithJSON(w, http.StatusOK, response)
}