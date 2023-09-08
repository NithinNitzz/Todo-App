// package main
//
// import (
//
//	"Todo-App/api"
//	"Todo-App/database"
//	"Todo-App/repository"
//	"Todo-App/service"
//	"database/sql"
//	"fmt"
//	"net/http"
//	"net/http/httptest"
//	"os"
//	"testing"
//
//	"github.com/DATA-DOG/go-sqlmock"
//	"github.com/labstack/echo/v4"
//	"github.com/stretchr/testify/assert"
//
// )
//
//	func setupTestDatabase(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
//		// Open a new database connection and create a SQL mock
//		db, mock, err := sqlmock.New()
//		if err != nil {
//			t.Fatalf("Error creating mock database: %v", err)
//		}
//
//		// Initialize the database schema (create tables, etc.)
//		if err := database.InitializeSchema(db); err != nil {
//			t.Fatalf("Error initializing schema: %v", err)
//		}
//
//		return db, mock
//	}
//
//	func teardownTestDatabase(db *sql.DB) {
//		// Close the database connection
//		db.Close()
//	}
//
//	func TestCreateTaskAPI(t *testing.T) {
//		// Set up the test database and get the mock
//		db, mock := setupTestDatabase(t)
//		defer teardownTestDatabase(db)
//
//		// Create a new Echo instance
//		e := echo.New()
//
//		// Initialize the mock database and repository
//		database.InitForTesting(db)
//		repo := repository.NewRepository(db)
//		srv := service.NewService(repo)
//		api := api.NewApi(srv)
//		api.Routes(e)
//
//		// Define a test task
//		testTask := `{"name": "Test Task", "description": "Test Description", "due_date": "2023-09-01", "status": "pending", "user_id": 1}`
//
//		// Define the expected SQL query for task creation
//		mock.ExpectQuery(`INSERT INTO tasks`).WithArgs(
//			"Test Task",
//			"Test Description",
//			"2023-09-01",
//			"pending",
//		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
//
//		// Create a request for creating a task
//		req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
//		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//		rec := httptest.NewRecorder()
//		c := e.NewContext(req, rec)
//		c.SetPath("/tasks")
//
//		// Simulate a request to create a task
//		req2 := httptest.NewRequest(http.MethodPost, "/tasks", nil)
//		req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//		req2.Body = httptest.DefaultRequest{}.Body
//		rec2 := httptest.NewRecorder()
//		c2 := e.NewContext(req2, rec2)
//
//		if assert.NoError(t, api.CreateTask(c2)) {
//			assert.Equal(t, http.StatusCreated, rec2.Code)
//			assert.JSONEq(t, testTask, rec2.Body.String())
//		}
//
//		// Assert that the mock database expectations were met
//		if err := mock.ExpectationsWereMet(); err != nil {
//			t.Errorf("Database expectations were not met: %s", err)
//		}
//	}
//
//	func TestMain(m *testing.M) {
//		// Set up and tear down the test database
//		db, mock := setupTestDatabase()
//		defer teardownTestDatabase(db)
//		defer mock.ExpectationsWereMet()
//
//		fmt.Println("Setting up test database...")
//		// No need to call database.InitTestDB or database.DropTestDB here
//		// Initialize the database schema and any other setup needed here
//
//		// Run the tests
//		code := m.Run()
//
//		fmt.Println("Tearing down test database...")
//		// Perform any additional teardown if necessary
//
//		// Exit with the status code from the tests
//		os.Exit(code)
//	}
package main

import (
	"Todo-App/api"
	"Todo-App/database"
	"Todo-App/models"
	"Todo-App/repository"
	"Todo-App/service"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//var ErrTaskNotFound = errors.New("task not found")

func TestCreateTask(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()
	connStr := "postgresql://nithin:office@localhost/Todo_App?sslmode=disable"
	db, err := database.Init(connStr)
	if err != nil {
		panic(err)
	}

	//create repository
	repo := repository.NewRepository(db)

	//create service
	srv := service.NewService(repo)

	//create Api
	api := api.NewApi(srv)
	api.Routes(e)

	// Create a request body with JSON data
	requestBody := `{
        "name": "Task Name",
        "description": "Task Description",
        "due_date": "2023-09-30T00:00:00Z",
        "status": "Incomplete",
        "created_at": "2023-08-30T00:00:00Z",
        "user_id": 1
    }`

	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call your CreateTask handler
	if err := api.CreateTask(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check the response status code (expected 201 for successful creation)
	if rec.Code != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, rec.Code)
	}

	// Parse the response body into a Task struct or validate it as needed
	var createdTask models.Task
	if err := json.Unmarshal(rec.Body.Bytes(), &createdTask); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// You can now perform assertions on the createdTask or any other needed validations
	if createdTask.Name != "Task Name" {
		t.Errorf("expected task name to be 'Task Name', got '%s'", createdTask.Name)
	}
	if createdTask.Status != "Incomplete" {
		t.Errorf("expected task status to be 'Incomplete', got '%s'", createdTask.Status)
	}

	//check response headers
	//contentType := rec.Header().Get(echo.HeaderContentType)
	//expectedContentType := echo.MIMEApplicationJSON
	//if !strings.HasPrefix(contentType, expectedContentType) {
	//	t.Errorf("expected content type to start with '%s',got '%s'", expectedContentType, contentType)
	//}
}
func TestUserTask(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()
	connStr := "postgresql://nithin:office@localhost/Todo_App?sslmode=disable"
	db, err := database.Init(connStr)
	if err != nil {
		panic(err)
	}

	//create repository
	repo := repository.NewRepository(db)

	//create service
	srv := service.NewService(repo)

	//create Api
	api := api.NewApi(srv)
	api.Routes(e)

	// Create a request body with JSON data
	requestBody := `{
        "username": "Task Name",
        "password":  "Task Password"
        
    }`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call your CreateTask handler
	if err := api.CreateUser(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check the response status code (expected 201 for successful creation)
	if rec.Code != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, rec.Code)
	}

	// Parse the response body into a Task struct or validate it as needed
	var createdTask models.Users
	if err := json.Unmarshal(rec.Body.Bytes(), &createdTask); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// You can now perform assertions on the createdTask or any other needed validations
	if createdTask.Username != "Task Name" {
		t.Errorf("expected task name to be 'Task Name', got '%s'", createdTask.Username)
	}
	//if createdTask.Status != "Incomplete" {
	//	t.Errorf("expected task status to be 'Incomplete', got '%s'", createdTask.Status)
	//}
}

//func TestGetTask(t *testing.T) {
//	// Create a new Echo instance
//	e := echo.New()
//	connStr := "postgresql://nithin:office@localhost/Todo_App?sslmode=disable"
//	db, err := database.Init(connStr)
//	if err != nil {
//		panic(err)
//	}
//
//	//create repository
//	repo := repository.NewRepository(db)
//
//	//create service
//	srv := service.NewService(repo)
//
//	//create Api
//	api := api.NewApi(srv)
//	api.Routes(e)
//
//	// Define a sample task ID for testing
//	taskID := "1"
//
//	// Create a GET request to retrieve a task
//	req := httptest.NewRequest(http.MethodGet, "/tasks/:id"+taskID, nil)
//	rec := httptest.NewRecorder()
//	c := e.NewContext(req, rec)
//
//	// Call your GetTask handler
//	if err := api.GetTask(c); err != nil {
//		if rec.Code == http.StatusNotFound {
//			expecteedError := `{"error : "task not found"}`
//			if rec.Body.String() != expecteedError {
//				t.Errorf("expected error response:%s,got:%s", expecteedError, rec.Body.String())
//			}
//			//t.Logf("task not found: %s", err)
//			return
//		}
//		//if err == ErrTaskNotFound {
//		//	rec.WriteHeader(http.StatusNotFound)
//		//	rec.Write([]byte(`{"error":"task not found"`))
//		//	return
//		//}
//		t.Fatalf("expected no error, got %v", err)
//	}
//
//	// Check the response status code (expected 200 for successful retrieval)
//	if rec.Code != http.StatusOK {if task.Name != "Sample Task" {
//		t.Errorf("expected task name to be 'Sample Task', got '%s'", task.Name)
//	}
//		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
//	}
//
//	// Parse the response body into a Task struct or validate it as needed
//	var retrievedTask models.Task
//	if err := json.Unmarshal(rec.Body.Bytes(), &retrievedTask); err != nil {
//		t.Fatalf("failed to unmarshal response: %v", err)
//	}
//
//	// Assert specific fields in the retrieved task
//	if retrievedTask.ID != 1 {
//		t.Errorf("expected task ID to be 1, got %d", retrievedTask.ID)
//	}

//if retrievedTask.Name != "Sample Task" {
//	t.Errorf("expected task name to be 'Sample Task', got '%s'", retrievedTask.Name)
//}
// Add more field-specific assertions as needed

// Optionally, you can check response headers, content type, etc.
//}

//func TestGetTask(t *testing.T) {
//	// Create a new Echo instance and set up your API routes
//	e := echo.New()
//	connStr := "postgresql://nithin:office@localhost/Todo_App?sslmode=disable"
//	db, err := database.Init(connStr)
//	if err != nil {
//		panic(err)
//	}
//
//	//create repository
//	repo := repository.NewRepository(db)
//
//	//create service
//	srv := service.NewService(repo)
//	api := api.NewApi(srv)
//	api.Routes(e)
//
//	// Create a request to retrieve a task with a specific ID (replace 1 with the desired task ID)
//	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
//	rec := httptest.NewRecorder()
//	c := e.NewContext(req, rec)
//
//	// Call the GetTask handler
//	if assert.NoError(t, api.GetTask(c)) {
//		// Check the response status code (replace 200 with the expected status code)
//		assert.Equal(t, http.StatusOK, rec.Code)
//
//		// You can also check the response body for the expected task properties (e.g., ID and Name)
//		// Replace the expected values with the ones you expect to receive
//		expectedID := 1
//		//expectedName := "Sample Task"
//
//		// Parse the response body into a structure (e.g., models.Task)
//		// Replace models.Task with the actual structure you use in your code
//		var task models.Task
//		err := json.Unmarshal(rec.Body.Bytes(), &task)
//		if assert.NoError(t, err) {
//			// Check the task's properties
//			assert.Equal(t, expectedID, task.ID)
//			//assert.Equal(t, expectedName, task.Name)
//		}
//	}
//}

//func TestGetTask(t *testing.T) {
//	// Create a new Echo instance
//	//e := echo.New()
//	//connStr := "postgresql://nithin:office@localhost/Todo_App?sslmode=disable"
//	//db, err := database.Init(connStr)
//	//if err != nil {
//	//	panic(err)
//	//}
//	//
//	////create repository
//	//repo := repository.NewRepository(db)
//	//
//	////create service
//	//srv := service.NewService(repo)
//	//
//	////create Api
//	//api := api.NewApi(srv)
//	//api.Routes(e)
//
//	e := echo.New()
//	req := httptest.NewRequest(http.MethodGet, "/tasks/:id", nil)
//	rec := httptest.NewRecorder()
//	c := e.NewContext(req, rec)
//	c.SetPath("/tasks/:id")
//	c.SetParamNames("id")
//	c.SetParamValues("1")
//	//h := &handler{mockDB}
//
//	var createdTask models.Task
//	if err := json.Unmarshal(rec.Body.Bytes(), &createdTask); err != nil {
//		t.Fatalf("failed to unmarshal response: %v", err)
//	}
//	if createdTask.Name != "Task Name" {
//		t.Errorf("expected task name to be 'Task Name', got '%s'", createdTask.Name)
//	}
//	if createdTask.Status != "Incomplete" {
//		t.Errorf("expected task status to be 'Incomplete', got '%s'", createdTask.Status)
//	}
//}

//	func TestGetTask(t *testing.T) {
//		// Create a new Echo instance
//		e := echo.New()
//		connStr := "postgresql://nithin:office@localhost/Todo_App?sslmode=disable"
//		db, err := database.Init(connStr)
//		if err != nil {
//			panic(err)
//		}
//
//		//create repository
//		repo := repository.NewRepository(db)
//
//		//create service
//		srv := service.NewService(repo)
//
//		//create Api
//		api := api.NewApi(srv)
//		api.Routes(e)
//		//requestBody := `{
//		//    "name": "Task Name",
//		//    "description": "Task Description",
//		//    "due_date": "2023-09-30T00:00:00Z",
//		//    "status": "Incomplete",
//		//    "created_at": "2023-08-30T00:00:00Z",
//		//    "user_id": 1
//		//}`
//
//		// Create a request to retrieve a task with a specific ID (replace 1 with the desired task ID)
//		req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
//		rec := httptest.NewRecorder()
//		c := e.NewContext(req, rec)
//		c.SetPath("/tasks/:id")
//		c.SetParamNames("id")
//		c.SetParamValues("1")
//
//		// Call your GetTask handler
//		if err := api.GetTask(c); err != nil {
//			t.Fatalf("expected no error, got %v", err)
//		}
//
//		//// Check the response status code (expected 200 for successful retrieval)
//		if rec.Code != http.StatusOK {
//			t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
//		}
//		//
//		//// Parse the response body into a Task struct or validate it as needed
//		//var task models.Task
//		//if err := json.Unmarshal(rec.Body.Bytes(), &task); err != nil {
//		//	t.Fatalf("failed to unmarshal response: %v", err)
//		//}
//		//
//		//// Optionally, perform assertions on the retrieved task or other response aspects
//		//if task.ID != 1 {
//		//	t.Errorf("expected task ID to be 1, got %d", task.ID)
//		//}
//		//if task.Name != "Sample Task" {
//		//	t.Errorf("expected task name to be 'Sample Task', got '%s'", task.Name)
//		//}
//
//		// Optionally, you can also check response headers, content type, and more as shown in the previous example for CreateTask.
//	}
func TestGetTask(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()
	connStr := "postgresql://nithin:office@localhost/Todo_App?sslmode=disable"
	db, err := database.Init(connStr)
	if err != nil {
		panic(err)
	}

	//create repository
	repo := repository.NewRepository(db)

	//create service
	srv := service.NewService(repo)

	//create Api
	api := api.NewApi(srv)
	api.Routes(e)

	// Create a request to retrieve a task with a specific ID
	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tasks/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	if err := api.GetTask(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}
	var tasks models.Task
	if err := json.Unmarshal(rec.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if tasks.ID != 1 {
		t.Errorf("expected task ID to be 1, got %d", tasks.ID)
	}

}
func TestGetUser(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()
	connStr := "postgresql://nithin:office@localhost/Todo_App?sslmode=disable"
	db, err := database.Init(connStr)
	if err != nil {
		panic(err)
	}

	//create repository
	repo := repository.NewRepository(db)

	//create service
	srv := service.NewService(repo)

	//create Api
	api := api.NewApi(srv)
	api.Routes(e)

	// Create a request to retrieve a task with a specific ID
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	if err := api.GetUser(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}
	var users models.Users
	if err := json.Unmarshal(rec.Body.Bytes(), &users); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if users.Id != 1 {
		t.Errorf("expected task ID to be 1, got %d", users.Id)
	}

}
func TestUpdateTask(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()
	connStr := "postgresql://nithin:office@localhost/Todo_App?sslmode=disable"
	db, err := database.Init(connStr)
	if err != nil {
		panic(err)
	}

	//create repository
	repo := repository.NewRepository(db)

	//create service
	srv := service.NewService(repo)

	//create Api
	api := api.NewApi(srv)
	api.Routes(e)

	// Create a request body with JSON data
	requestBody := `{
        "name": "updated Task Name",
        "description": "updated Task Description",
        "due_date": "2023-09-30T00:00:00Z",
        "status": "completed",
        "updated_at": "2023-08-30T00:00:00Z",
        
    }`

	req := httptest.NewRequest(http.MethodPut, "/tasks/1", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tasks/:id")
	c.SetParamNames("name", "description", "due_date", "status", "updated_at")
	c.SetParamValues("snowman", "snow", "2023-08-15T15:35:10.535606+05:30", "IN_PROGRESS", "2024-08-16T15:35:10.535606+05:30")

	// Call your CreateTask handler
	if err := api.UpdateTask(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check the response status code (expected 200 for successful creation)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}
	//

	//var createdTask models.Task
	//if err := json.Unmarshal(rec.Body.Bytes(), &createdTask); err != nil {
	//	t.Fatalf("failed to unmarshal response: %v", err)
	//}
	//
	//// You can now perform assertions on the createdTask or any other needed validations
	//if createdTask.Name != "updated Task Name" {
	//	t.Errorf("expected task name to be 'Task Name', got '%s'", createdTask.Name)
	//}
	//if createdTask.Status != "completed" {
	//	t.Errorf("expected task status to be 'Incomplete', got '%s'", createdTask.Status)
	//}
}
