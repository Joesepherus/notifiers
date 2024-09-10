package errorUtils

import (
	"net/http"
)

func MethodNotAllowed_error(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/error?message=method+not+allowed", http.StatusSeeOther)
		return
	}
}
