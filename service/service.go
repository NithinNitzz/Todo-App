package service

import (
	"Todo-App/models"
	"Todo-App/repository"
)

type Service interface {
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
type service struct {
	Repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{Repo: repo}
}

func (s *service) CreateTask(t models.Task) (int, error) {

	return s.Repo.CreateTask(t)
}

func (s *service) GetTask(id string) (*models.Task, error) {
	return s.Repo.GetTask(id)
}

func (s *service) CreateUser(user models.Users) (int, error) {
	return s.Repo.CreateUser(user)
}

func (s *service) GetUser(id string) (*models.Users, error) {
	return s.Repo.GetUser(id)
}

func (s *service) UpdateTask(t models.Task) error {
	return s.Repo.UpdateTask(t)
}

func (s *service) CompleteTask(t models.Task) error {
	return s.Repo.CompleteTask(t)
}

func (s *service) TempDeleteTask(t models.Task) error {
	return s.Repo.TempDeleteTask(t)
}

func (s *service) DeleteTask(id string) (int, error) {
	return s.Repo.DeleteTask(id)
}

func (s *service) DeleteUserTask(id string) (int, error) {
	return s.Repo.DeleteUserTask(id)
}

func (s *service) ListTasks(userId string, page, perPage int) ([]models.Task, error) {
	return s.Repo.ListTasks(userId, page, perPage)
}

func (s *service) GetUserByName(username string) (models.Users, error) {
	return s.Repo.GetUserByName(username)
}

func (s *service) RegisterUser(user models.Users) (int, error) {
	return s.Repo.RegisterUser(user)
}

func (s *service) LoginUser(user models.Users) (string, error) {
	return s.Repo.LoginUser(user)
}
