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

type Post struct {
	Id    int
	Title string
	Body  string
}

type GenericRepository[T any, V any] interface {
	Get(id V) (T, error)
	Post(T) (T, error)
}

type UserService struct {
	UserRepository GenericRepository[User, int]
}

func NewUserService(userRepository GenericRepository[User, int]) *UserService {
	return &UserService{
		UserRepository: userRepository,
	}
}

func (u *UserService) Get(id int) (User, error) {
	user, err := u.UserRepository.Get(id)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (u *UserService) Post(user User) (User, error) {
	user.Name = strings.ToUpper(user.Name)
	dbUser, err := u.UserRepository.Post(user)
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

func (u *SQLLiteUserRepository) Get(id int) (User, error) {
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

func (u *SQLLiteUserRepository) Post(user User) (User, error) {
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

type InMemoryUserRepository struct {
	Users map[int]User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		Users: make(map[int]User),
	}
}

func (u *InMemoryUserRepository) Get(id int) (User, error) {
	user, ok := u.Users[id]
	if !ok {
		return User{}, fmt.Errorf("User not found")
	}
	return user, nil
}

func (u *InMemoryUserRepository) Post(user User) (User, error) {
	u.Users[user.Age] = user
	return user, nil
}

func main() {
	dbRepository := NewSQLLiteUserRepository("users.db")
	UserService := NewUserService(dbRepository)
	user, err := UserService.Post(User{"Daniel", 1})
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
	user, err = UserService.Get(2)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
}
