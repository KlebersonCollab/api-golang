package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/update/"):]

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Abre o banco de dados SQLite.
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Executa uma consulta SQL para atualizar o usuário no banco de dados.
	result, err := db.Exec("UPDATE users SET username=?, password=?, email=?, date=?, phone=? WHERE id=?", user.Username, user.Password, user.Email, user.Date, user.Phone, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verifica se o usuário foi atualizado com sucesso.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Retorna uma resposta de sucesso com o JSON contendo os dados atualizados do usuário.
	user.ID = id
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
