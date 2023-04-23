package users

import "net/http"

func GetHealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("welcome"))
	}
}
