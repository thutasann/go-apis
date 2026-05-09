package main

import (
	"fmt"
	"net/http"
	"time"
)

// - Open browser → http://localhost:8080
// - Refresh or close tab quickly
// - No wasted CPU
// - No leaked goroutines
// - Massive cost savings at scale
func http_context_cancel_handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Println("Request started")

	select {
	case <-time.After(5 * time.Second):
		fmt.Fprintln(w, "Work completed")
	case <-ctx.Done():
		fmt.Println("Request cancelled: ", ctx.Err())
		return
	}
}

func HTTP_Context_Cancellation() {
	http.HandleFunc("/", http_context_cancel_handler)

	fmt.Println("Serer is running on :8080")
	http.ListenAndServe(":8080", nil)
}
