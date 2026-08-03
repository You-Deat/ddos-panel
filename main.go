package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "time"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    cors := func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            next(w, r)
        }
    }

    http.HandleFunc("/api/health", cors(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "status": "ok",
            "time":   time.Now().String(),
        })
    }))

    http.HandleFunc("/api/start", cors(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var req struct {
            Target   string `json:"target"`
            Duration int    `json:"duration"`
            Cookie   string `json:"cookie"`
        }

        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }

        if req.Target == "" || req.Duration <= 0 || req.Duration > 120 {
            http.Error(w, "Invalid parameters", http.StatusBadRequest)
            return
        }

        go func() {
            fmt.Printf("🔥 Attack started on %s for %ds\n", req.Target, req.Duration)
            time.Sleep(time.Duration(req.Duration) * time.Second)
            fmt.Printf("✅ Attack finished on %s\n", req.Target)
        }()

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "status":  "started",
            "message": "Attack started",
            "id":      time.Now().Unix(),
        })
    }))

    fmt.Printf("🚀 Backend running on port %s\n", port)
    http.ListenAndServe(":"+port, nil)
}
