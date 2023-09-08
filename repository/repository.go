package repository

import (
	"Todo-App/models"
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type Repository interface {
	CreateTask(t models.Task) (int, error)
	GetTask(id string) (*models.Task, error)
	CreateUser(user models.Users) (int, error)

	GetUser(id string) (*models.Users, error)

	UpdateTask(t models.Task) error
	CompleteTask(t models.Task) error
	TempDeleteTask(t models.Task) error
	DeleteTask(id string) (int, error)
	DeleteUserTask(id string) (int, error)
	ListTasks(userId string, page, perPage int) ([]models.Task, error)
	GetUserByName(username string) (models.Users, error)
	RegisterUser(user models.Users) (int, error)
	LoginUser(user models.Users) (string, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateTask(t models.Task) (int, error) {
	query := "INSERT INTO tasks (name, description, due_date, status, created_at, user_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	err := r.db.QueryRow(query, t.Name, t.Desc, t.DueDate, t.Status, t.CreatedAT, t.UserId).Scan(&t.ID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to create task")
	}

	return http.StatusCreated, nil
}

func (r *repository) GetTask(id string) (*models.Task, error) {
	query := "SELECT id, name, description, due_date, status, created_at FROM tasks WHERE id = $1"
	row := r.db.QueryRow(query, id)

	t := new(models.Task)
	err := row.Scan(&t.ID, &t.Name, &t.Desc, &t.DueDate, &t.Status, &t.CreatedAT)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}
	return t, nil
}
func (r *repository) CreateUser(u models.Users) (int, error) {
	query := "INSERT INTO users (username) VALUES ($1) RETURNING id"
	err := r.db.QueryRow(query, u.Username).Scan(&u.Id)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to create task")
	}

	return http.StatusCreated, nil
}

func (r *repository) GetUser(id string) (*models.Users, error) {
	query := "SELECT id, username FROM Users WHERE id = $1"
	row := r.db.QueryRow(query, id)

	u := new(models.Users)
	err := row.Scan(&u.Id, &u.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}
	return u, nil
}

func (r *repository) UpdateTask(t models.Task) error {
	query := "UPDATE tasks SET name = $1, description = $2, due_date = $3, status = $4, updated_at= $5 WHERE id = $6"
	_, err := r.db.Exec(query, t.Name, t.Desc, t.DueDate, t.Status, t.UpdatedAt, t.ID)
	if err != nil {
		return fmt.Errorf("failed to update task")
	}
	return nil
}
func (r *repository) CompleteTask(t models.Task) error {
	query := "UPDATE tasks SET status = $1 WHERE id = $2"
	_, err := r.db.Exec(query, t.Status, t.ID)
	if err != nil {
		return fmt.Errorf("failed to complete task")
	}
	return nil
}
func (r *repository) TempDeleteTask(t models.Task) error {
	query := "UPDATE tasks SET deleted_at = $1 WHERE id = $2"
	_, err := r.db.Exec(query, t.DeletedAt, t.ID)
	if err != nil {
		return fmt.Errorf("failed to delete task")
	}
	return nil
}
func (r *repository) DeleteTask(id string) (int, error) {
	query := "DELETE FROM tasks WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to delete task")

	}
	return http.StatusOK, nil
}
func (r *repository) DeleteUserTask(id string) (int, error) {
	query := "DELETE FROM Users WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to delete task")
	}
	return http.StatusOK, nil
}
func (r *repository) ListTasks(userId string, page, perPage int) ([]models.Task, error) {
	offset := (page - 1) * perPage
	query := "SELECT t.id, t.name, t.description, t.due_date, t.status, t.created_at ,t.user_id,u.id,u.username FROM tasks as t left join Users as u on u.id=t.user_id WHERE t.user_id = $1 LIMIT $2  OFFSET $3"
	row, err := r.db.Query(query, userId, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get task")
	}
	defer row.Close()

	var tasks []models.Task
	for row.Next() {

		t := models.Task{}
		err := row.Scan(&t.ID, &t.Name, &t.Desc, &t.DueDate, &t.Status, &t.CreatedAT, &t.UserId, &t.User.Id, &t.User.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task")
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
func (r *repository) GetUserByName(username string) (models.Users, error) {
	query := "SELECT id, username, password FROM Users WHERE username = $1"
	row := r.db.QueryRow(query, username)

	var u models.Users
	err := row.Scan(&u.Id, &u.Username, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, nil
		}
		return u, err
	}

	return u, nil
}
func (r *repository) RegisterUser(user models.Users) (int, error) {
	exuser, err := r.GetUserByName(user.Username)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to create user")
	}
	if exuser.Username != "" {
		return http.StatusConflict, fmt.Errorf("username already exists")
	}
	query := "INSERT INTO Users (username, password) VALUES ($1 , $2) RETURNING id"
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to create user")
	}
	//query := "INSERT INTO Users (username, password) VALUES ($1 , $2) RETURNING id"
	err = r.db.QueryRow(query, user.Username, string(pass)).Scan(&user.Id)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to create user")
	}

	user.Password = ""
	return http.StatusCreated, nil
}
func (r *repository) LoginUser(user models.Users) (string, error) {
	exuser, err := r.GetUserByName(user.Username)
	if err != nil {
		return "", fmt.Errorf("failed to login")
	}
	if exuser.Username == "" {
		return "", fmt.Errorf("user not exist")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(exuser.Password), []byte(user.Password)); err != nil {

		return "", fmt.Errorf("invalid data")
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = exuser.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	secretKey := "todoapp"

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to create token")
	}
	return tokenString, nil
}

//LoginUser(user users) (string, error)
//CreateTask(c echo.Context) error
//CreateUserTask(c echo.Context) error
//GetTask(c echo.Context) error
//GetUserTask(c echo.Context) error
//GetUserByName(username string) (models.Users, error)
//UpdateTask(c echo.Context) error
//CompleteTask(c echo.Context) error
//TempDeleteTask(c echo.Context) error
//DeleteTask(c echo.Context) error
//DeleteUserTask(c echo.Context) error
//ListTasks(c echo.Context) error
//RegisterUser(c echo.Context) error
//LoginUser(c echo.Context) error

//func (r *Reapository) CreateTask(c echo.Context) error {
//	var t models.Task
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
//
////func NewRepository(db *sql.DB) repository {
////	return &repository{db}
////}
//
//func (r *Reapository) CreateUserTask(c echo.Context) error {
//	//TODO implement me
//	var u models.Users
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
//func (r *Reapository) GetTask(c echo.Context) error {
//	//TODO implement me
//	id := c.Param("id")
//	query := "SELECT id, name, description, due_date, status, created_at FROM tasks WHERE id = $1"
//	row := r.db.QueryRow(query, id)
//
//	t := new(models.Task)
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
//func (r *Reapository) GetUserTask(c echo.Context) error {
//	//TODO implement me
//	id := c.Param("id")
//	query := "SELECT id, username FROM Users WHERE id = $1"
//	row := r.db.QueryRow(query, id)
//
//	u := new(models.Users)
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
//func (r *Reapository) GetUserByName(username string) (models.Users, error) {
//
//	//username := c.Param("username")
//	query := "SELECT id, username, password FROM Users WHERE username = $1"
//	row := r.db.QueryRow(query, username)
//
//	var u models.Users
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
//func (r *Reapository) UpdateTask(c echo.Context) error {
//	//TODO implement me
//	id := c.Param("id")
//	var t models.Task
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
//func (r *Reapository) CompleteTask(c echo.Context) error {
//	//TODO implement me
//	id := c.Param("id")
//	var t models.Task
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
//
//func (r *Reapository) TempDeleteTask(c echo.Context) error {
//	//TODO implement me
//	id := c.Param("id")
//	var t models.Task
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
//
//func (r *Reapository) DeleteTask(c echo.Context) error {
//	//TODO implement me
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
//
//func (r *Reapository) DeleteUserTask(c echo.Context) error {
//	//TODO implement me
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
//func (r *Reapository) ListTasks(c echo.Context) error {
//	//TODO implement me
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
//	var tasks []models.Task
//	for row.Next() {
//
//		t := models.Task{}
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
//func (r *Reapository) RegisterUser(c echo.Context) error {
//	//TODO implement me
//	var user models.Users
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
//}
//
//func (r *Reapository) LoginUser(c echo.Context) error {
//	//TODO implement me
//	var user models.Users
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
