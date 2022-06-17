package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string
	Age  int
}

type UserRepository interface {
	GetUser(id int) (User, error)
	PostUser(user User) (User, error)
}

type UserService struct {
	UserRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		UserRepository: userRepository,
	}
}

func (u *UserService) GetUser(id int) (User, error) {
	return u.UserRepository.GetUser(id)
}

func (u *UserService) PostUser(user User) (User, error) {
	user.Name = strings.ToUpper(user.Name)
	dbUser, err := u.UserRepository.PostUser(user)
	if err != nil {
		return User{}, err
	}
	return dbUser, nil
}

type SQLLiteUserRepository struct {
	DB *sql.DB
}

func NewSQLLiteUserRepository(url string) *SQLLiteUserRepository {
	db, err := sql.Open("sqlite3", url)
	if err != nil {
		panic(err)
	}
	return &SQLLiteUserRepository{db}
}

func (u *SQLLiteUserRepository) GetUser(id int) (User, error) {
	query := "SELECT name, age FROM users WHERE id = ?"
	row := u.DB.QueryRow(query, id)
	var name string
	var age int
	err := row.Scan(&name, &age)
	if err != nil {
		return User{}, err
	}
	return User{name, age}, nil
}

func (u *SQLLiteUserRepository) PostUser(user User) (User, error) {
	query := "INSERT INTO users (name, age) VALUES (?, ?)"
	_, err := u.DB.Exec(query, user.Name, user.Age)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func CreateTable(db *sql.DB) {
	query := "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT, age INTEGER)"
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func main() {
	dbRepository := NewSQLLiteUserRepository("users.db")
	UserService := NewUserService(dbRepository)
	user, err := UserService.PostUser(User{"Daniel", 30})
	if err != nil {
		panic(err)
	}
	user, err = UserService.GetUser(3)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
}
