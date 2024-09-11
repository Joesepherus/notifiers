package errorUtils

import (
	"log"
	"net/http"
	"tradingalerts/services/loggingService"
)

func MethodNotAllowed_error(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Method not allowed")
		loggingService.LogToDB("ERROR", "Method not allowed", r)
		http.Redirect(w, r, "/error?message=method+not+allowed", http.StatusSeeOther)
		return
	}
}
