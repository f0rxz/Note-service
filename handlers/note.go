package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"note-service/models"
	"note-service/storage"
	"note-service/utils"
	"time"

	"github.com/gorilla/mux"
)

var noteStore = storage.NewMemoryStorage()

func GetNotes(w http.ResponseWriter, r *http.Request) {
	notes := noteStore.GetAll()
	if notes == nil {
		notes = []*models.Note{}
	}
	utils.RespondWithJSON(w, http.StatusOK, notes)
}

func CreateNote(w http.ResponseWriter, r *http.Request) {
	var note models.Note
	json.NewDecoder(r.Body).Decode(&note)
	note.ID = generateID() // Generate a unique ID for the note
	noteStore.Create(&note)
	utils.RespondWithJSON(w, http.StatusCreated, note)
}

func GetNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	note, found := noteStore.Get(id)
	if !found {
		utils.RespondWithError(w, http.StatusNotFound, "Note not found")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, note)
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var note models.Note
	json.NewDecoder(r.Body).Decode(&note)
	updatedNote, err := noteStore.Update(id, &note)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Note not found")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, updatedNote)
}

func DeleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := noteStore.Delete(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Note not found")
		return
	}
	utils.RespondWithJSON(w, http.StatusNoContent, nil)
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
