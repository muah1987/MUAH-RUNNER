package api

import (
"encoding/json"
"net/http"
)

type WebhookPayload struct {
Event string                 `json:"event"`
Data  map[string]interface{} `json:"data"`
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
var payload WebhookPayload
if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
http.Error(w, "bad request", http.StatusBadRequest)
return
}
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(map[string]string{"status": "received"})
}
