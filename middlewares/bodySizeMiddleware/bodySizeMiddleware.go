package bodySizeMiddleware

import "net/http"

const LIMIT int64 = 10 * 1024 * 1024

func LimitRequestBodySize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > LIMIT {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}

		// Wrap the request body reader to enforce the limit
		limitedReader := http.MaxBytesReader(w, r.Body, LIMIT)
		r.Body = limitedReader
		next.ServeHTTP(w, r)
	})
}
