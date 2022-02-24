package user

import (
	"encoding/json"
	"net/http"
)

func RegisterHandler() {
	http.HandleFunc("/", findAllUser)
}

func findAllUser(w http.ResponseWriter, r *http.Request) {
	marshal, err := json.Marshal(users)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(marshal)
	w.WriteHeader(http.StatusOK)
}
