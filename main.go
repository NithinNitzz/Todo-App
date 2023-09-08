package main

import (
	"Todo-App/api"
	"Todo-App/database"
	_ "Todo-App/docs"

	"Todo-App/repository"
	_ "Todo-App/repository"
	"Todo-App/service"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//type Task struct {
//	ID        int       `json:"id,omitempty"`
//	Name      string    `json:"name"`
//	Desc      string    `json:"description"`
//	DueDate   time.Time `json:"due_date"`
//	Status    string    `json:"status"`
//	CreatedAT time.Time `json:"created_at"`
//	UpdatedAt time.Time `json:"updated_at"`
//	DeletedAt time.Time `json:"deleted_at"`
//	UserId    int       `json:"user_id"`
//	User      Users     `json:"user"`
//}
//type Users struct {
//	Id       int    `json:"id"`
//	Username string `json:"username"`
//	Password string `json:"password"`
//}

//type repository interface {
//	CreateTask(c echo.Context) error
//	CreateUser(c echo.Context) error
//	GetTask(c echo.Context) error
//	GetUser(c echo.Context) error
//	GetUserByName(c echo.Context) (*Users, error)
//	UpdateTask(c echo.Context) error
//	CompleteTask(c echo.Context) error
//	TempDeleteTask(c echo.Context) error
//	DeleteTask(c echo.Context) error
//	DeleteUser(c echo.Context) error
//	ListTasks(c echo.Context) error
//	RegisterUser(c echo.Context) error
//	LoginUser(c echo.Context) error
//}

//type repository struct {
//	db *sql.DB
//}
//
//func (r repository) CreateTask(c echo.Context) error {
//	var t Task
//	if err := c.Bind(&t); err != nil {
//		return c.String(http.StatusBadRequest, "Invalid request data")
//	}
//
//	query := "INSERT INTO tasks (name, description, due_date, status, created_at, user_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
//	err := r.db.QueryRow(query, t.Name, t.Desc, t.DueDate, t.Status, t.CreatedAT, t.UserId).Scan(&t.ID)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "Failed to create task")
//	}
//
//	return c.JSON(http.StatusCreated, t)
//}
//func (r repository) CreateUserTask(c echo.Context) error {
//	var u Users
//	if err := c.Bind(&u); err != nil {
//		return c.String(http.StatusBadRequest, "Invalid request data")
//	}
//
//	query := "INSERT INTO Users (username) VALUES ($1) RETURNING id"
//	err := r.db.QueryRow(query, u.Username).Scan(&u.Id)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "Failed to create task")
//	}
//
//	return c.JSON(http.StatusCreated, u)
//}
//
//func (r repository) GetTask(c echo.Context) error {
//	id := c.Param("id")
//	query := "SELECT id, name, description, due_date, status, created_at FROM tasks WHERE id = $1"
//	row := r.db.QueryRow(query, id)
//
//	t := new(Task)
//	err := row.Scan(&t.ID, &t.Name, &t.Desc, &t.DueDate, &t.Status, &t.CreatedAT)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return c.String(http.StatusNotFound, "Task not found")
//		}
//		return c.String(http.StatusInternalServerError, "Failed to get task")
//	}
//
//	return c.JSON(http.StatusOK, t)
//}
//
//func (r repository) GetUserTask(c echo.Context) error {
//	id := c.Param("id")
//	query := "SELECT id, username FROM Users WHERE id = $1"
//	row := r.db.QueryRow(query, id)
//
//	u := new(Users)
//	err := row.Scan(&u.Id, &u.Username)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return c.String(http.StatusNotFound, "Task not found")
//		}
//		return c.String(http.StatusInternalServerError, "Failed to get task")
//	}
//
//	return c.JSON(http.StatusOK, u)
//}
//
//func (r repository) GetUserByName(username string) (Users, error) {
//	//username := c.Param("username")
//	query := "SELECT id, username, password FROM Users WHERE username = $1"
//	row := r.db.QueryRow(query, username)
//
//	var u Users
//	err := row.Scan(&u.Id, &u.Username, &u.Password)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return u, nil
//		}
//		return u, err
//	}
//
//	return u, nil
//}
//
//func (r repository) UpdateTask(c echo.Context) error {
//	id := c.Param("id")
//	var t Task
//	if err := c.Bind(&t); err != nil {
//		return c.String(http.StatusBadRequest, "Invalid request data")
//	}
//
//	query := "UPDATE tasks SET name = $1, description = $2, due_date = $3, status = $4, updated_at= $5 WHERE id = $6"
//	_, err := r.db.Exec(query, t.Name, t.Desc, t.DueDate, t.Status, t.UpdatedAt, id)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "Failed to update task")
//	}
//
//	return c.JSON(http.StatusOK, t)
//}
//
//func (r repository) CompleteTask(c echo.Context) error {
//	id := c.Param("id")
//	var t Task
//	if err := c.Bind(&t); err != nil {
//		return c.String(http.StatusBadRequest, "Invalid request data")
//	}
//	query := "UPDATE tasks SET status = $1 WHERE id = $2"
//	_, err := r.db.Exec(query, t.Status, id)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "Failed to complete task")
//	}
//
//	return c.JSON(http.StatusOK, t)
//}
//func (r repository) TempDeleteTask(c echo.Context) error {
//	id := c.Param("id")
//	var t Task
//	if err := c.Bind(&t); err != nil {
//		return c.String(http.StatusBadRequest, "Invalid request data")
//	}
//	query := "UPDATE tasks SET deleted_at = $1 WHERE id = $2"
//	_, err := r.db.Exec(query, t.DeletedAt, id)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "Failed to complete task")
//	}
//
//	return c.JSON(http.StatusOK, t)
//}
//func (r repository) DeleteTask(c echo.Context) error {
//	id := c.Param("id")
//
//	query := "DELETE FROM tasks WHERE id = $1"
//	_, err := r.db.Exec(query, id)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "Failed to delete task")
//	}
//
//	return c.String(http.StatusOK, fmt.Sprintf("Task with ID %s deleted", id))
//}
//func (r repository) DeleteUserTask(c echo.Context) error {
//	id := c.Param("id")
//
//	query := "DELETE FROM Users WHERE id = $1"
//	_, err := r.db.Exec(query, id)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "Failed to delete task")
//	}
//
//	return c.String(http.StatusOK, fmt.Sprintf("Task with ID %s deleted", id))
//}
//
//func (r repository) ListTasks(c echo.Context) error {
//	userId := c.Param("user_id")
//
//	page, err := strconv.Atoi(c.QueryParam("page"))
//	if err != nil || page <= 0 {
//		page = 1
//	}
//	perPage, err := strconv.Atoi(c.QueryParam("perPage"))
//	if err != nil || perPage <= 0 {
//		perPage = 10
//	}
//	offset := (page - 1) * perPage
//
//	query := "SELECT t.id, t.name, t.description, t.due_date, t.status, t.created_at ,t.user_id,u.id,u.username FROM tasks as t left join Users as u on u.id=t.user_id WHERE t.user_id = $1 LIMIT $2  OFFSET $3"
//	row, err := r.db.Query(query, userId, perPage, offset)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "Failed to get task")
//	}
//	defer row.Close()
//
//	var tasks []Task
//	for row.Next() {
//
//		t := Task{}
//		err := row.Scan(&t.ID, &t.Name, &t.Desc, &t.DueDate, &t.Status, &t.CreatedAT, &t.UserId, &t.User.Id, &t.User.Username)
//		if err != nil {
//			if err == sql.ErrNoRows {
//				return c.String(http.StatusNotFound, "Task not found")
//			}
//			return c.String(http.StatusInternalServerError, "Failed to get task")
//		}
//		tasks = append(tasks, t)
//	}
//	return c.JSON(http.StatusOK, tasks)
//}
//
//func (r repository) Registeruser(c echo.Context) error {
//	var user Users
//	if err := c.Bind(&user); err != nil {
//		return c.String(http.StatusBadRequest, "invalid user data")
//	}
//
//	exuser, err := r.GetUserByName(user.Username)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "failed to create user")
//	}
//	if exuser.Username != "" {
//		return c.String(http.StatusConflict, "username already exists")
//	}
//	//cost := 12
//	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "failed to create Users ")
//	}
//
//	query := "INSERT INTO Users (username, password) VALUES ($1 , $2) RETURNING id"
//	err = r.db.QueryRow(query, user.Username, string(pass)).Scan(&user.Id)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "failed to create user")
//	}
//
//	user.Password = ""
//
//	return c.JSON(http.StatusCreated, user)
//
//}
//
//func (r repository) Loginuser(c echo.Context) error {
//	var user Users
//	if err := c.Bind(&user); err != nil {
//		return c.String(http.StatusBadRequest, "invalid user data")
//	}
//
//	exuser, err := r.GetUserByName(user.Username)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "failed to login")
//	}
//	//if exuser.Username == "" {
//	//	return c.String(http.StatusNotFound, "user not exist")
//	//}
//	if err = bcrypt.CompareHashAndPassword([]byte(exuser.Password), []byte(user.Password)); err != nil {
//
//		return c.String(http.StatusUnauthorized, "invalid data")
//	}
//
//	token := jwt.New(jwt.SigningMethodHS256)
//	claims := token.Claims.(jwt.MapClaims)
//	claims["username"] = exuser.Username
//	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
//
//	secretKey := "todoapp"
//
//	tokenString, err := token.SignedString([]byte(secretKey))
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "failed to create token")
//	}
//
//	return c.JSON(http.StatusOK, map[string]string{
//		"message": "login successful",
//		"token":   tokenString,
//	})
//}

//func NewRepository(db *sql.DB) repository {
//	return &repository{db}
//}

//var db *sql.DB

//func init() {
//
//	connStr := "postgresql://nithin:office@localhost/Todo_App?sslmode=disable"
//	var err error
//	db, err = sql.Open("postgres", connStr)
//	if err != nil {
//		panic(err)
//	}
//}

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
//

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	e := echo.New()

	connStr := "postgresql://nithin:office@localhost/Todo_App?sslmode=disable"
	db, err := database.Init(connStr)
	if err != nil {
		panic(err)
	}

	//repo := &repository{db: db}

	//repo := NewRepository(db)

	//create repository
	repo := repository.NewRepository(db)

	//create service
	srv := service.NewService(repo)

	//create Api
	api := api.NewApi(srv)
	api.Routes(e)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

//e.POST("/tasks", func(c echo.Context) error {
//	var t models.Task
//	if err := c.Bind(&t); err != nil {
//		return c.String(http.StatusBadRequest, "Invalid request data")
//	}
//
//	id, err := repo.CreateTask(t)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "failed to create task")
//	}
//	t.ID = id
//
//	return c.JSON(http.StatusCreated, t)
//
//})
//
//e.GET("/tasks/:id", func(c echo.Context) error {
//	id := c.Param("id")
//	task, err := repo.GetTask(id)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "failed to get task")
//	}
//	return c.JSON(http.StatusOK, task)
//
//})

//e.POST("/tasks", repo.CreateTask)
//e.POST("/Users", repo.CreateUserTask)
//e.GET("/tasks/:id", repo.GetTask)
//e.GET("/Users/:id", repo.GetUserTask)
//
//e.PUT("/tasks/:id", repo.UpdateTask)
//e.PATCH("/tasks/:id/complete", repo.CompleteTask)
//e.DELETE("/tasks/tempdel/:id", repo.TempDeleteTask)
//e.DELETE("/tasks/:id", repo.DeleteTask)
//e.DELETE("/users/:id", repo.DeleteUserTask)
//e.GET("/users/:user_id/tasks", repo.ListTasks)
//e.POST("/register", repo.RegisterUser)
//e.POST("/login", repo.LoginUser)
//e.GET("/swagger/*", echoSwagger.WrapHandler)

//	e.Logger.Fatal(e.Start(":8080"))
//}

//func createTask(c echo.Context) error {
//	var t Task
//	if err := c.Bind(&t); err != nil {
//		return c.String(http.StatusBadRequest, "Invalid request data")
//	}
//
//	query := "INSERT INTO tasks (name, description, due_date, status,created_at,user_id) VALUES ($1, $2, $3, $4,$5,$6 ) RETURNING id"
//	err := db.QueryRow(query, t.Name, t.Desc, t.DueDate, t.Status, t.CreatedAT, t.UserId).Scan(&t.ID)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "Failed to create task")
//	}
//
//	return c.JSON(http.StatusCreated, t)
//}
//func createuserTask(c echo.Context) error {
//	var u Users
//	if err := c.Bind(&u); err != nil {
//		return c.String(http.StatusBadRequest, "Invalid request data")
//	}
//
//	query := "INSERT INTO Users (username) VALUES ($1) RETURNING id"
//	err := db.QueryRow(query, u.Username).Scan(&u.Id)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "Failed to create task")
//	}
//
//	return c.JSON(http.StatusCreated, u)
//}
//
//func getTask(c echo.Context) error {
//	id := c.Param("id")
//	query := "SELECT id, name, description, due_date, status, created_at FROM tasks WHERE id = $1"
//	row := db.QueryRow(query, id)
//
//	t := new(Task)
//	err := row.Scan(&t.ID, &t.Name, &t.Desc, &t.DueDate, &t.Status, &t.CreatedAT)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return c.String(http.StatusNotFound, "Task not found")
//		}
//		return c.String(http.StatusInternalServerError, "Failed to get task")
//	}
//
//	return c.JSON(http.StatusOK, t)
//}
//
//func getuserTask(c echo.Context) error {
//	id := c.Param("id")
//	query := "SELECT id, username FROM Users WHERE id = $1"
//	row := db.QueryRow(query, id)
//
//	u := new(Users)
//	err := row.Scan(&u.Id, &u.Username)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return c.String(http.StatusNotFound, "Task not found")
//		}
//		return c.String(http.StatusInternalServerError, "Failed to get task")
//	}
//
//	return c.JSON(http.StatusOK, u)
//}
//
//func getuserbyname(username string) (Users, error) {
//	//username := c.Param("username")
//	query := "SELECT id, username, password FROM Users WHERE username = $1"
//	row := db.QueryRow(query, username)
//
//	var u Users
//	err := row.Scan(&u.Id, &u.Username, &u.Password)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return u, nil
//		}
//		return u, err
//	}
//
//	return u, nil
//}
//
//func updateTask(c echo.Context) error {
//	id := c.Param("id")
//	var t Task
//	if err := c.Bind(&t); err != nil {
//		return c.String(http.StatusBadRequest, "Invalid request data")
//	}
//
//	query := "UPDATE tasks SET name = $1, description = $2, due_date = $3, status = $4, updated_at= $5 WHERE id = $6"
//	_, err := db.Exec(query, t.Name, t.Desc, t.DueDate, t.Status, t.UpdatedAt, id)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "Failed to update task")
//	}
//
//	return c.JSON(http.StatusOK, t)
//}
//
//func completeTask(c echo.Context) error {
//	id := c.Param("id")
//	var t Task
//	if err := c.Bind(&t); err != nil {
//		return c.String(http.StatusBadRequest, "Invalid request data")
//	}
//	query := "UPDATE tasks SET status = $1 WHERE id = $2"
//	_, err := db.Exec(query, t.Status, id)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "Failed to complete task")
//	}
//
//	return c.JSON(http.StatusOK, t)
//}
//func tempdelTask(c echo.Context) error {
//	id := c.Param("id")
//	var t Task
//	if err := c.Bind(&t); err != nil {
//		return c.String(http.StatusBadRequest, "Invalid request data")
//	}
//	query := "UPDATE tasks SET deleted_at = $1 WHERE id = $2"
//	_, err := db.Exec(query, t.DeletedAt, id)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "Failed to complete task")
//	}
//
//	return c.JSON(http.StatusOK, t)
//}
//func deleteTask(c echo.Context) error {
//	id := c.Param("id")
//
//	query := "DELETE FROM tasks WHERE id = $1"
//	_, err := db.Exec(query, id)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "Failed to delete task")
//	}
//
//	return c.String(http.StatusOK, fmt.Sprintf("Task with ID %s deleted", id))
//}
//func deleteuserTask(c echo.Context) error {
//	id := c.Param("id")
//
//	query := "DELETE FROM Users WHERE id = $1"
//	_, err := db.Exec(query, id)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "Failed to delete task")
//	}
//
//	return c.String(http.StatusOK, fmt.Sprintf("Task with ID %s deleted", id))
//}
//
//func listTasks(c echo.Context) error {
//	userId := c.Param("user_id")
//
//	page, err := strconv.Atoi(c.QueryParam("page"))
//	if err != nil || page <= 0 {
//		page = 1
//	}
//	perPage, err := strconv.Atoi(c.QueryParam("perPage"))
//	if err != nil || perPage <= 0 {
//		perPage = 10
//	}
//	offset := (page - 1) * perPage
//
//	query := "SELECT t.id, t.name, t.description, t.due_date, t.status, t.created_at ,t.user_id,u.id,u.username FROM tasks as t left join Users as u on u.id=t.user_id WHERE t.user_id = $1 LIMIT $2  OFFSET $3"
//	row, err := db.Query(query, userId, perPage, offset)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "Failed to get task")
//	}
//	defer row.Close()
//
//	var tasks []Task
//	for row.Next() {
//
//		t := Task{}
//		err := row.Scan(&t.ID, &t.Name, &t.Desc, &t.DueDate, &t.Status, &t.CreatedAT, &t.UserId, &t.User.Id, &t.User.Username)
//		if err != nil {
//			if err == sql.ErrNoRows {
//				return c.String(http.StatusNotFound, "Task not found")
//			}
//			return c.String(http.StatusInternalServerError, "Failed to get task")
//		}
//		tasks = append(tasks, t)
//	}
//	return c.JSON(http.StatusOK, tasks)
//}
//
//func Registeruser(c echo.Context) error {
//	var user Users
//	if err := c.Bind(&user); err != nil {
//		return c.String(http.StatusBadRequest, "invalid user data")
//	}
//
//	exuser, err := getuserbyname(user.Username)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "failed to create user")
//	}
//	if exuser.Username != "" {
//		return c.String(http.StatusConflict, "username already exists")
//	}
//	//cost := 12
//	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "failed to create Users ")
//	}
//
//	query := "INSERT INTO Users (username, password) VALUES ($1 , $2) RETURNING id"
//	err = db.QueryRow(query, user.Username, string(pass)).Scan(&user.Id)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "failed to create user")
//	}
//
//	user.Password = ""
//
//	return c.JSON(http.StatusCreated, user)
//
//}
//
//func loginuser(c echo.Context) error {
//	var user Users
//	if err := c.Bind(&user); err != nil {
//		return c.String(http.StatusBadRequest, "invalid user data")
//	}
//
//	exuser, err := getuserbyname(user.Username)
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "failed to login")
//	}
//	//if exuser.Username == "" {
//	//	return c.String(http.StatusNotFound, "user not exist")
//	//}
//	if err = bcrypt.CompareHashAndPassword([]byte(exuser.Password), []byte(user.Password)); err != nil {
//
//		return c.String(http.StatusUnauthorized, "invalid data")
//	}
//
//	token := jwt.New(jwt.SigningMethodHS256)
//	claims := token.Claims.(jwt.MapClaims)
//	claims["username"] = exuser.Username
//	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
//
//	secretKey := "todoapp"
//
//	tokenString, err := token.SignedString([]byte(secretKey))
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "failed to create token")
//	}
//
//	return c.JSON(http.StatusOK, map[string]string{
//		"message": "login successful",
//		"token":   tokenString,
//	})
//}
