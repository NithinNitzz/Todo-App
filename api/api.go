package api

import (
	"Todo-App/jwtutil"
	"Todo-App/models"
	_ "Todo-App/repository"
	"Todo-App/service"
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	_ "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"net/http"
	"strconv"
	"strings"
)

type Api struct {
	//Repo      repository.Repository
	Validator *validator.Validate
	Service   service.Service
}

//	func NewApi(repo repository.Repository) *Api {
//		return &Api{Repo: repo,
//			Validator: validator.New()}
//	}

func NewApi(srv service.Service) *Api {
	return &Api{Service: srv,
		Validator: validator.New()}

}

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")
			fmt.Println("received token:", tokenString)
			if tokenString == "" {
				return c.String(http.StatusUnauthorized, "Missing JWT token")
			}
			tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer"))
			token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte("todoapp"), nil
			})
			if err != nil {
				fmt.Println("error parsing token:", err)
				return c.String(http.StatusUnauthorized, "invalid jwt token")
			}
			claims, ok := token.Claims.(*jwt.StandardClaims)
			if !ok || !token.Valid {
				fmt.Println("invalid token or claims:", claims)
				return c.String(http.StatusUnauthorized, "invalid jwt token")
			}

			//username, err := jwtutil.VerifyToken(tokenString)
			//if err != nil {
			//	return c.String(http.StatusUnauthorized, "Invalid JWT token")
			//}
			//
			//// Set the username in the context for further use
			//c.Set("username", username)

			return next(c)
		}
	}
}

func (a *Api) Routes(e *echo.Echo) {
	//e.Group("", JWTMiddleware())
	//{
	//apiV1 := e.Group("/api/v1")
	e.POST("/tasks", a.CreateTask)
	e.GET("/tasks/:id", a.GetTask, JWTMiddleware())
	e.POST("/users", a.CreateUser)
	e.GET("/users/:id", a.GetUser)
	e.PUT("/tasks/:id", a.UpdateTask)
	e.PATCH("/tasks/:id/complete", a.CompleteTask)
	e.DELETE("/tasks/tempdel/:id", a.TempDeleteTask)
	e.DELETE("/tasks/:id", a.DeleteTask)
	e.DELETE("/users/:id", a.DeleteUserTask)
	e.GET("/users/:user_id/tasks", a.ListTasks)
	//}
	e.POST("/register", a.RegisterUser)
	e.POST("/login", a.LoginUser)
	//e.GET("/swagger/*", echoSwagger.WrapHandler)

}

// CreateTask godoc
// @Summary      create a task
// @Description  create a new task
// @Tags         Task
// @Accept       json
// @Produce      json
// @Param        t   body models.Task true "Task object to be created"
// @Success      200  {object}  models.Task
// @Router       /tasks [post]
func (a *Api) CreateTask(c echo.Context) error {
	var t models.Task
	if err := c.Bind(&t); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request data")
	}

	//if err := a.Validator.Struct(t); err != nil {
	//	return c.String(http.StatusBadRequest, err.Error())
	//}

	id, err := a.Service.CreateTask(t)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to create task")
	}
	t.ID = id

	return c.JSON(http.StatusCreated, t)
}

// GetTask godoc
// @Summary      Get a task
// @Description  Get a task by id
// @Tags         Task
// @Accept       json
// @Produce      json
// @Param        id   path string true "Task ID"
// @Success      200  {object}  models.Task
// @Router       /tasks/{id} [get]
func (a *Api) GetTask(c echo.Context) error {
	//	id := c.Param("id")
	//	taskId, err := strconv.Atoi(id)
	//	if err != nil {
	//		return c.JSON(http.StatusBadRequest, map[string]string{"error": "task not found"})
	//	}
	//	task, err := a.Service.GetTask(strconv.Itoa(taskId))
	//	if err != nil {
	//		//var ErrTaskNotFound = errors.New("task not found")
	//		if err == ErrTaskNotFound {
	//			return c.JSON(http.StatusNotFound, map[string]string{"error": "task not found"})
	//		}
	//		return err
	//		//return c.String(http.StatusInternalServerError, "failed to get task")
	//	}
	//	return c.JSON(http.StatusOK, task)
	//}
	id := c.Param("id")
	task, err := a.Service.GetTask(id)
	if err != nil {
		fmt.Println("error fetching task", err)
		if err == sql.ErrNoRows {
			return c.String(http.StatusNotFound, "task not found")
		}
		return c.String(http.StatusInternalServerError, "failed to get task")
	}
	return c.JSON(http.StatusOK, task)
}

// CreateUser godoc
// @Summary      create a task
// @Description  create a new task
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        u   body models.Users true "Task object to be created"
// @Success      200  {object}  models.Users
// @Router       /users [post]
func (a *Api) CreateUser(c echo.Context) error {
	var u models.Users
	if err := c.Bind(&u); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request data")
	}

	id, err := a.Service.CreateUser(u)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to create task")
	}
	u.Id = id

	return c.JSON(http.StatusCreated, u)
}

// GetUser godoc
// @Summary      Get a User task
// @Description  Get a task by id
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path string true "Users Task ID"
// @Success      200  {object}  models.Users
// @Router       /users/{id} [get]
func (a *Api) GetUser(c echo.Context) error {
	id := c.Param("id")
	user, err := a.Service.GetUser(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to get task")
	}
	return c.JSON(http.StatusOK, user)
}

// UpdateTask godoc
// @Summary      update a task
// @Description  update ex task
// @Tags         Task
// @Id           taskID
// @Accept       json
// @Produce      json
// @Param        id  path string true "Task id"
// @Param        id  body  models.Task true "Task object to be updated"
// @Success      200  {object}  models.Task
// @Router       /tasks/{id} [put]
func (a *Api) UpdateTask(c echo.Context) error {
	id := c.Param("id")
	var t models.Task
	if err := c.Bind(&t); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request data")
	}

	taskID, err := strconv.Atoi(id)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid task id")
	}

	t.ID = taskID

	err = a.Service.UpdateTask(t)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to update task")
	}

	return c.JSON(http.StatusOK, t)

}

// CompleteTask godoc
// @Summary      complete a task
// @Description  complete a task by id
// @Tags         Task
// @Accept       json
// @Produce      json
// @Param        id   path string true "Task ID"
// @Param        id   body models.Task true "Task completed"
// @Success      200  {object}  models.Task
// @Router       /tasks/{id}/complete [patch]
func (a *Api) CompleteTask(c echo.Context) error {
	id := c.Param("id")
	var t models.Task
	if err := c.Bind(&t); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request data")
	}
	taskID, err := strconv.Atoi(id)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid task id")
	}

	t.ID = taskID

	err = a.Service.CompleteTask(t)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to get task")
	}
	return c.JSON(http.StatusOK, t)
}

// TempDeleteTask godoc
// @Summary      tem delete a task
// @Description  tem delete a task by id
// @Tags         Task
// @Accept       json
// @Produce      json
// @Param        id   path string true "Task ID"
// @Param        id   body models.Task true " tem delete task"
// @Success      200  {object}  models.Task
// @Router       /tasks/tempdel/{id} [delete]
func (a *Api) TempDeleteTask(c echo.Context) error {
	id := c.Param("id")
	var t models.Task
	if err := c.Bind(&t); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request data")
	}
	taskID, err := strconv.Atoi(id)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid task id")
	}

	t.ID = taskID
	err = a.Service.TempDeleteTask(t)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to get task")
	}
	return c.JSON(http.StatusOK, t)

}

// DeleteTask godoc
// @Summary      Delete a task
// @Description  Delete a task by id
// @Tags         Task
// @Accept       json
// @Produce      json
// @Param        id   path string true "Task ID"
// @Success      200  {object}  models.Task
// @Router       /tasks/{id} [delete]
func (a *Api) DeleteTask(c echo.Context) error {
	id := c.Param("id")
	status, err := a.Service.DeleteTask(id)
	if err != nil {
		return c.String(status, fmt.Sprintf("Failed to delete task: %s", err.Error()))
	}
	return c.String(http.StatusOK, fmt.Sprintf("Task with ID %s deleted", id))
}
func (a *Api) DeleteUserTask(c echo.Context) error {
	id := c.Param("id")
	status, err := a.Service.DeleteUserTask(id)
	if err != nil {
		return c.String(status, fmt.Sprintf("Failed to delete task: %s", err.Error()))
	}
	return c.String(http.StatusOK, fmt.Sprintf("Task with ID %s deleted", id))
}

// ListTasks godoc
// @Summary      list all
// @Description  list all tasks by id
// @Tags         Task
// @Accept       json
// @Produce      json
// @Param        user_id   path string true "Users_ID"
// @Success      200  {object}  models.Task
// @Router       /users/{user_id}/tasks [get]
func (a *Api) ListTasks(c echo.Context) error {
	userId := c.Param("user_id")
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	perPage, err := strconv.Atoi(c.QueryParam("perPage"))
	if err != nil || perPage <= 0 {
		perPage = 10
	}

	tasks, err := a.Service.ListTasks(userId, page, perPage)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, tasks)
}

// RegisterUser godoc
// @Summary      Register a user
// @Description  add a new user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user   body models.Users true "User Data"
// @Success      200  {object}  models.Users
// @Router       /register [post]
func (a *Api) RegisterUser(c echo.Context) error {
	var user models.Users
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, "invalid user data")
	}
	if err := a.Validator.Struct(user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	status, err := a.Service.RegisterUser(user)
	if err != nil {
		return c.String(status, err.Error())
	}
	tokenString, err := jwtutil.GenerateToken(user.Username)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create token")
	}

	return c.JSON(status, map[string]string{
		"message": "user registered",
		"token":   tokenString,
	})
	//return c.JSON(status, user)
}

// LoginUser godoc
// @Summary      Login a user
// @Description  Login for  user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user   body models.Users true "User Data"
// @Success      200  {object}  models.Users
// @Router       /login [post]
func (a *Api) LoginUser(c echo.Context) error {
	var user models.Users
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, "invalid user data")
	}
	if err := a.Validator.Struct(user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	//tokenString, err := a.Service.LoginUser(user)
	//if err != nil {
	//	return c.String(http.StatusUnauthorized, "invalid data")
	//}
	tokenString, err := jwtutil.GenerateToken(user.Username)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create token")
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "login successful",
		"token":   tokenString,
	})
}
