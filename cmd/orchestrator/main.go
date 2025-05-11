package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sklerakuku/5final/internal/auth"
	calculation "github.com/sklerakuku/5final/internal/calculator"
	"github.com/sklerakuku/5final/internal/config"
	"github.com/sklerakuku/5final/internal/db"
	"github.com/sklerakuku/5final/internal/grpc"
	"github.com/sklerakuku/5final/internal/server"
)

var dbConn *sql.DB
var cfg *config.Config

func main() {
	log.Println("Starting application...")

	cfg = config.Load()

	var err error
	dbConn, err = sql.Open("sqlite3", cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer dbConn.Close()

	if err = dbConn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Database connection established successfully")

	log.Println("Initializing database tables...")
	sqlBytes, err := os.ReadFile("./migrations/001_init.sql")
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}

	if _, err := dbConn.Exec(string(sqlBytes)); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Println("Database initialized successfully")

	auth.SetDB(dbConn)
	auth.SetJWTSecret(cfg.JWTSecret)
	server.SetJWTSecret(cfg.JWTSecret)
	db.SetDB(dbConn)

	r := chi.NewRouter()

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "web"))
	log.Printf("Serving static files from: %s", filepath.Join(workDir, "web"))
	r.Handle("/*", http.FileServer(filesDir))

	log.Println("Setting up API routes...")
	r.Post("/api/v1/register", handleRegister)
	r.Post("/api/v1/login", handleLogin)
	r.With(server.AuthMiddleware).Post("/api/v1/calculate", handleCalculate)
	r.With(server.AuthMiddleware).Get("/api/v1/expressions/{id}", handleGetExpression)
	r.With(server.AuthMiddleware).Get("/api/v1/expressions", handleGetAllExpressions)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving index.html to client %s", r.RemoteAddr)
		http.ServeFile(w, r, filepath.Join(workDir, "web", "index.html"))
	})

	log.Printf("Server started on :%s \nVisit http://localhost:%s/", cfg.ServerPort, cfg.ServerPort)
	http.ListenAndServe(":"+cfg.ServerPort, r)
}

func sendJSONError(w http.ResponseWriter, message string, status int) {
	log.Printf("Sending JSON error (status %d): %s", status, message)
	log.Printf("ERROR: %s (Status: %d)", message, status)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	log.Printf("Register request from %s", r.RemoteAddr)
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if len(req.Login) < 3 || len(req.Password) < 6 {
		sendJSONError(w, "Username must be at least 3 characters and password at least 6", http.StatusBadRequest)
		return
	}

	if err := auth.Register(req.Login, req.Password); err != nil {
		if err.Error() == "username already exists" {
			sendJSONError(w, "Username already exists", http.StatusConflict)
		} else {
			log.Printf("Registration error: %v", err)
			sendJSONError(w, "Registration failed. Please try again.", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("User %s registered successfully", req.Login)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful"})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("Login request from %s", r.RemoteAddr)
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	token, err := auth.Login(req.Login, req.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			sendJSONError(w, "Invalid username or password", http.StatusUnauthorized)
		} else {
			log.Printf("Login error: %v", err)
			sendJSONError(w, "Login failed. Please try again.", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("User %s logged in successfully", req.Login)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func handleCalculate(w http.ResponseWriter, r *http.Request) {
	log.Printf("Calculate request from %s", r.RemoteAddr)
	claims := r.Context().Value("claims").(jwt.MapClaims)

	var userID int
	switch sub := claims["sub"].(type) {
	case float64:
		userID = int(sub)
	case string:
		var id int
		err := dbConn.QueryRow("SELECT id FROM users WHERE login = ?", sub).Scan(&id)
		if err != nil {
			sendJSONError(w, "Invalid user", http.StatusUnauthorized)
			return
		}
		userID = id
	default:
		sendJSONError(w, "Invalid token claims", http.StatusBadRequest)
		return
	}

	client, err := grpc.NewClient(cfg.GRPCAddress, cfg)
	if err != nil {
		log.Printf("Failed to create gRPC client: %v", err)
		sendJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	id, err := db.SaveExpression(userID, req.Expression)
	if err != nil {
		log.Printf("Failed to save expression for user %d: %v", userID, err)
		sendJSONError(w, "Failed to save expression", http.StatusInternalServerError)
		return
	}

	value, err := calculation.Parse(req.Expression)
	if err != nil {
		log.Printf("Failed to parse expression '%s': %v", req.Expression, err)
		db.UpdateExpressionStatus(id, "failed")
		sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	db.UpdateExpressionResult(id, "completed", value)

	log.Printf("User %d calculated: %s = %v", userID, req.Expression, value)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     id,
		"result": value,
	})
}

func handleGetExpression(w http.ResponseWriter, r *http.Request) {
	log.Printf("Get expression request from %s", r.RemoteAddr)
	claims := r.Context().Value("claims").(jwt.MapClaims)
	var userID int
	switch sub := claims["sub"].(type) {
	case float64:
		userID = int(sub)
	case string:
		var id int
		err := dbConn.QueryRow("SELECT id FROM users WHERE login = ?", sub).Scan(&id)
		if err != nil {
			sendJSONError(w, "Invalid user", http.StatusUnauthorized)
			return
		}
		userID = id
	default:
		sendJSONError(w, "Invalid token claims", http.StatusBadRequest)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	expr, err := db.GetExpressionByID(id)
	if err != nil {
		log.Printf("Expression not found: %v", err)
		sendJSONError(w, "Expression not found", http.StatusNotFound)
		return
	}

	if expr.UserID != userID {
		log.Printf("Unauthorized get expression attempt: expr.UserID=%d, userID=%d", expr.UserID, userID)
		sendJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("User %d fetched expression id %d", userID, id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expr)
}

type ExpressionResponse struct {
	ID         int      `json:"id"`
	Expression string   `json:"expression"`
	Status     string   `json:"status"`
	Result     *float64 `json:"result"`
}

func handleGetAllExpressions(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt.MapClaims)
	var userID int
	switch sub := claims["sub"].(type) {
	case float64:
		userID = int(sub)
	case string:
		var id int
		err := dbConn.QueryRow("SELECT id FROM users WHERE login = ?", sub).Scan(&id)
		if err != nil {
			sendJSONError(w, "Invalid user", http.StatusUnauthorized)
			return
		}
		userID = id
	default:
		sendJSONError(w, "Invalid token claims", http.StatusBadRequest)
		return
	}

	expressions, _ := db.GetExpressionsByUser(userID)
	var responses []ExpressionResponse
	for _, expr := range expressions {
		var result *float64
		if expr.Result.Valid {
			result = &expr.Result.Float64
		}
		responses = append(responses, ExpressionResponse{
			ID:         expr.ID,
			Expression: expr.Expr,
			Status:     expr.Status,
			Result:     result,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"expressions": responses,
	})
}
