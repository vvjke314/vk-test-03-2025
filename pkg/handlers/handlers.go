package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/vvjke314/vk-test-03-2025/internal/custom_errors"
	"github.com/vvjke314/vk-test-03-2025/internal/entities"
	"github.com/vvjke314/vk-test-03-2025/internal/usecases"
)

type KVHandler struct {
	uc     *usecases.KeyValueUseCase
	logger *log.Logger
}

func NewKVHandler(uc *usecases.KeyValueUseCase, logger *log.Logger) *KVHandler {
	return &KVHandler{
		uc:     uc,
		logger: logger,
	}
}

// CreateKeyHandler handle POST /kv
func (h *KVHandler) CreateKeyHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Printf("Request to create key-value pair: %s %s", r.Method, r.URL.Path)

	var req struct {
		Key   string          `json:"key"`
		Value json.RawMessage `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Bad JSON body: %v", err)
		http.Error(w, `{"error": "bad json"}`, http.StatusBadRequest)
		return
	}

	if req.Key == "" {
		h.logger.Printf("empty key")
		http.Error(w, `{"error": "key can not be empty"}`, http.StatusBadRequest)
		return
	}

	if !json.Valid(req.Value) {
		h.logger.Printf("json did not validate %s", req.Key)
		http.Error(w, `{"error": "bad value JSON"}`, http.StatusBadRequest)
		return
	}

	item := entities.VaultItem{
		Key:   req.Key,
		Value: string(req.Value),
	}

	err := h.uc.InsertValue(item)
	if err != nil {
		h.logger.Printf("error while creating key %s: %v", req.Key, err)

		switch {
		case err.Error() == fmt.Sprintf("key '%s' already exists", item.Key):
			h.logger.Printf("key already exists: %s", req.Key)
			http.Error(w, `{"error": "key already exists"}`, http.StatusConflict)
		default:
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	h.logger.Printf("successfully created key: %s", req.Key)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// UpdateKeyHandler handles PUT /kv/{id}
func (h *KVHandler) UpdateKeyHandler(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("id")
	h.logger.Printf("Request to update key: %s %s", r.Method, r.URL.Path)

	var req struct {
		Value json.RawMessage `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("JSON decode error for key %s: %v", key, err)
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// Validate that value contains proper JSON
	if !json.Valid(req.Value) {
		h.logger.Printf("Invalid JSON in value for key %s", key)
		http.Error(w, `{"error": "Invalid JSON in value"}`, http.StatusBadRequest)
		return
	}

	err := h.uc.UpdateValue(key, string(req.Value))
	if err != nil {
		if errors.Is(err, custom_errors.ErrKeyNotExists) {
			h.logger.Printf("Key not found: %s", key)
			http.Error(w, `{"error": "Key not found"}`, http.StatusNotFound)
			return
		}

		h.logger.Printf("Error updating key %s: %v", key, err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	h.logger.Printf("Successfully updated key: %s", key)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GetKeyHandler handles GET /kv/{id}
func (h *KVHandler) GetKeyHandler(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("id")
	h.logger.Printf("Request to get key: %s %s", r.Method, r.URL.Path)

	item, err := h.uc.Get(key)
	if err != nil {
		if errors.Is(err, custom_errors.ErrKeyNotExists) {
			h.logger.Printf("Key not found: %s", key)
			http.Error(w, `{"error": "Key not found"}`, http.StatusNotFound)
			return
		}

		h.logger.Printf("Error getting key %s: %v", key, err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	h.logger.Printf("Successfully retrieved key: %s", key)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"key":   item.Key,
		"value": json.RawMessage(item.Value),
	})
}

// DeleteKeyHandler handles DELETE /kv/{id}
func (h *KVHandler) DeleteKeyHandler(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("id")
	h.logger.Printf("Request to delete key: %s %s", r.Method, r.URL.Path)

	err := h.uc.DeleteRow(key)
	if err != nil {
		if errors.Is(err, custom_errors.ErrKeyNotExists) {
			h.logger.Printf("Key not found: %s", key)
			http.Error(w, `{"error": "Key not found"}`, http.StatusNotFound)
			return
		}

		h.logger.Printf("Error deleting key %s: %v", key, err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	h.logger.Printf("Successfully deleted key: %s", key)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
