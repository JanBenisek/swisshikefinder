package main

import (
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	// set various http headers (see details in elsewhere)
	// also need to allow content from schweizmobil.ch for the images
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self' https://images.schweizmobil.ch; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)

	})
}
