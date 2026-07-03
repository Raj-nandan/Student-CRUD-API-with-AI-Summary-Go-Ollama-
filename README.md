# Student CRUD API with AI Summary (Go + Ollama)

This is a simple RESTful API built in Go using the Gorilla Mux router. It allows you to perform basic CRUD operations on a list of students stored in memory. Additionally, it integrates with the Ollama API to generate a summary of a student's profile using the Llama3 language model.

## Features

- Create a new student
- Retrieve all students
- Retrieve a student by ID
- Update a student by ID
- Delete a student by ID
- Generate an AI-based summary of a student using Ollama



## Requirements

- Go 1.24.5
- Ollama installed and running with Llama3 model

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/Raj-nandan/Student-CRUD-API-with-AI-Summary-Go-Ollama-.git
cd Student-CRUD-API-with-AI-Summary-Go-Ollama-

```
### 2. Installing Dependencies

```bash
go mod tidy
```

### 3. Run the Server

```bash
go run main.go
```
server will start at: http://localhost:8080

### Start Ollama and Llama3

```bash
ollama run llama3
```
the model will be running on: http://localhost:11434


### Api EndPoints

/students (POST)
```
json data

{
  "id": 1,
  "name": "Raj",
  "age": 11,
  "email": "raj@gmail.com"
}
```

/students (GET)
```
return all the students data
```
/students/{id} (GET)
```
return students with particular 'ID'
```
/students/{id} {PUT)
```
updates student data with respected ID
```
/students/{id} (DELET)
```
Delete student with respected ID
```

/students/{id}/summary
```
Return summary of student data from LLM model




