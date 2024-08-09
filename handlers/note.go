package handlers

import (
	"encoding/json"
	"net/http"
	"note-service/models"
	"note-service/storage"
	"note-service/utils"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var noteStorage, _ = storage.NewStorage("notes.db", time.Second*5)

func GetNotes(w http.ResponseWriter, r *http.Request) {
	const limit = 25
	vars := mux.Vars(r)
	page, _ := strconv.ParseInt(vars["page"], 10, 64)
	notes := noteStorage.GetNotesRange(page*limit, limit)
	utils.RespondWithJSON(w, http.StatusOK, notes)
}

func CreateNote(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.FormValue("title")
	content := r.FormValue("content")

	noteStorage.CreateNote(title, content)
	w.WriteHeader(http.StatusCreated)
}

func GetNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	note := noteStorage.GetNote(id)
	if note != nil {
		utils.RespondWithJSON(w, http.StatusOK, note)
		return
	}
	utils.RespondWithError(w, http.StatusNotFound, "Note not found")
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	noteStorage.UpdateNote(id, &note)
	utils.RespondWithJSON(w, http.StatusOK, note)
}

func DeleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	noteStorage.DeleteNote(id)
	utils.RespondWithJSON(w, http.StatusNoContent, nil)
}
