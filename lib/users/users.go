package users

import (
	"context"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Users struct {
	conn *pgxpool.Conn
}

type UserResults struct {
	Users []User
}

type User struct {
	ID int16
	Firstname string
	Lastname string
	Email string
	Phone string
	Usertype string
}

func (result *UserResults) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (result *User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewUsers(pool *pgxpool.Pool, ctx context.Context) (Users, error) {
	acquire, err := pool.Acquire(ctx)

	if err != nil {
		return Users{}, err
	}

	return Users{
		conn: acquire,
	}, nil
}

func (u Users) GetUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := u.getUsers()
		render.Render(w, r, users)
	}
}

func (u Users) PostUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
				panic(err)
		}
		var user User
		unmarshalError := json.Unmarshal(body, &user)
		if err != nil || unmarshalError != nil {
			log.Println("Error from reading response body or unmarshaling response body", err, unmarshalError)
			http.Error(w, "Error getting response body", http.StatusUnauthorized)
		}
		userCreated := u.createUser(user)
		render.Render(w, r, userCreated)
	}
}

func (u Users) getUsers() (*UserResults) {
	GET_USERS_SQL := `SELECT * FROM users ORDER BY id ASC`;

	rows, _ := u.conn.Query(context.Background(), GET_USERS_SQL);
	users, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (User, error) {
		var n User
		err := row.Scan(&n.ID, &n.Firstname, &n.Lastname, &n.Email, &n.Phone, &n.Usertype)
		return n, err
	})
	if err != nil {
		log.Infof("CollectRows error: %v", err)
		return nil
	}
	var userResults UserResults
	userResults.Users = users

	return &userResults;
}

func (u Users) createUser(user User) (*User) {
	//TODO: if such email already exist, return bad request

	CREATE_USER_SQL := `INSERT INTO users (firstname, lastname, email, phone, usertype) VALUES ($1, $2, $3, $4, $5)`;
	_, err := u.conn.Exec(context.Background(), CREATE_USER_SQL, user.Firstname, user.Lastname, user.Email, user.Phone, user.Usertype);
	
	if err != nil {
		log.Infof("Error creating user: %v", err)
		return nil
	}

	GET_USER_SQL := `SELECT * FROM users WHERE email = $1`
	row := u.conn.QueryRow(context.Background(), GET_USER_SQL, user.Email);
	if row == nil {
		return nil
	}
	var n User
	_err := row.Scan(&n.ID, &n.Firstname, &n.Lastname, &n.Email, &n.Phone, &n.Usertype)

	if _err != nil {
		log.Infof("CollectOneRow error: %v", err)
		return nil
	}
	return &n;
}