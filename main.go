package main
//importing required packages
import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

// structure of student input
type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

// Building map datastructure for storing student data in memory data storage
var (
	students = make(map[int]Student)
	lock     = sync.RWMutex{}
)

// main function for routing
func main() {
	r := mux.NewRouter()
	

	// Routes
	r.HandleFunc("/students", createStudent).Methods("POST")
	r.HandleFunc("/students", getAllStudents).Methods("GET")
	r.HandleFunc("/students/{id}", getStudentByID).Methods("GET")
	r.HandleFunc("/students/{id}", updateStudentByID).Methods("PUT")
	r.HandleFunc("/students/{id}", deleteStudentByID).Methods("DELETE")
	r.HandleFunc("/students/{id}/summary", getStudentSummary).Methods("GET")

	// Test route for debugging
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Server is running!",
			"status":  "ok",
		})
	}).Methods("GET")

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Available endpoints:")
	fmt.Println("  GET  /test                    - Test endpoint")
	fmt.Println("  POST /students                - Create student")
	fmt.Println("  GET  /students                - Get all students")
	fmt.Println("  GET  /students/{id}           - Get student by ID")
	fmt.Println("  PUT  /students/{id}           - Update student")
	fmt.Println("  DELETE /students/{id}         - Delete student")
	fmt.Println("  GET  /students/{id}/summary   - Get AI summary")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Post request for creating student
func createStudent(w http.ResponseWriter, r *http.Request) {
	var student Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if student.Name == "" || student.Email == "" || student.Age <= 0 {
		http.Error(w, "Missing or invalid fields", http.StatusBadRequest)
		return
	}

	lock.Lock()
	defer lock.Unlock()
	students[student.ID] = student

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

// Get request for getting all students
func getAllStudents(w http.ResponseWriter, r *http.Request) {
	lock.RLock()
	defer lock.RUnlock()

	var result []Student
	for _, s := range students {
		result = append(result, s)
	}
	json.NewEncoder(w).Encode(result)
}

// Get request for getting student by id
func getStudentByID(w http.ResponseWriter, r *http.Request) {
	id := getIDFromRequest(w, r)
	if id == -1 {
		return
	}

	lock.RLock()
	defer lock.RUnlock()
	student, exists := students[id]
	if !exists {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(student)
}

// Put request for updating student by id
func updateStudentByID(w http.ResponseWriter, r *http.Request) {
	id := getIDFromRequest(w, r)
	if id == -1 {
		return
	}

	var student Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	lock.Lock()
	defer lock.Unlock()

	_, exists := students[id]
	if !exists {
		http.NotFound(w, r)
		return
	}

	student.ID = id
	students[id] = student
	json.NewEncoder(w).Encode(student)
}

// Delete request for deleting student by id
func deleteStudentByID(w http.ResponseWriter, r *http.Request) {
	id := getIDFromRequest(w, r)
	if id == -1 {
		return
	}

	lock.Lock()
	defer lock.Unlock()

	_, exists := students[id]
	if !exists {
		http.NotFound(w, r)
		return
	}
	delete(students, id)
	w.WriteHeader(http.StatusNoContent)
}

// Ollama Integration for getting student summary
func getStudentSummary(w http.ResponseWriter, r *http.Request) {
	id := getIDFromRequest(w, r)
	if id == -1 {
		return
	}

	lock.RLock()
	student, exists := students[id]
	lock.RUnlock()

	if !exists {
		http.NotFound(w, r)
		return
	}

	prompt := fmt.Sprintf("Summarize this student profile:\nName: %s\nAge: %d\nEmail: %s",
		student.Name, student.Age, student.Email)

	// Create proper JSON payload
	payload := map[string]interface{}{
		"model":  "llama3",
		"prompt": prompt,
		"stream": false,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to create request payload", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", strings.NewReader(string(jsonPayload)))
	if err != nil {
		http.Error(w, "Failed to call Ollama: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read Ollama response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Failed to parse Ollama response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the response from Ollama
	summary, ok := result["response"]
	if !ok {
		// If response is not found, return the full result for debugging
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  "No response field in Ollama result",
			"result": result,
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"summary": summary,
	})
}

// Utility function for getting id from request for all the routes
func getIDFromRequest(w http.ResponseWriter, r *http.Request) int {
	params := mux.Vars(r)
	idStr, ok := params["id"]
	if !ok {
		http.Error(w, "ID not found", http.StatusBadRequest)
		return -1
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return -1
	}
	return id
}
