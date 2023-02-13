package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Date     string `json:"date"`
	Phone    string `json:"phone"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verifica se o usuário já existe no banco de dados.
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	row := db.QueryRow("SELECT id FROM users WHERE username=? OR email=?", user.Username, user.Email)
	var existingID string
	err = row.Scan(&existingID)
	if err == nil {
		// Um usuário com o mesmo nome de usuário ou e-mail já existe.
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Aqui você pode fazer o que precisar com os dados do usuário recebido.
	// Cria a tabela "users", se ela não existir.
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (id TEXT, username TEXT, password TEXT, email TEXT, date TEXT, phone TEXT)`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	insertuser(db, user, w)
}

func insertuser(db *sql.DB, user User, w http.ResponseWriter) {
	// Gera uma UUID única para o ID do usuário.
	id := uuid.New().String()
	// Adiciona o ID do usuário à estrutura do usuário.
	user.ID = id
	// Insere o usuário no banco de dados.
	_, err := db.Exec("INSERT INTO users (id, username, password, email, date, phone) VALUES (?, ?, ?, ?, ?, ?)", user.ID, user.Username, user.Password, user.Email, user.Date, user.Phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Envia uma resposta de sucesso com o JSON contendo o ID do usuário.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
