package dal

import (
	"errors"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	defaultSessionLifetime = 31 * 24 * time.Hour
)

var (
	ErrNotFound = errors.New("Not found")
)

type User struct {
	login string
	conn  redis.Conn
}

func (u *User) Exists() (bool, error) {
	val, err := redis.Bool(u.conn.Do("exists", "user:"+u.login))
	if err != nil {
		return false, err
	}
	return val, nil
}

func (u *User) Delete() error {
	val, err := redis.Bool(u.conn.Do("del", "user:"+u.login))
	if err != nil {
		return err
	}
	if !val {
		return ErrNotFound
	}
	return nil
}

func (u *User) GetPassword() (string, error) {
	pass, err := redis.String(u.conn.Do("hget", "user:"+u.login, "pass"))
	if err != nil {
		return "", err
	}
	return pass, nil
}

func (u *User) SetPassword(pass string) error {
	_, err := u.conn.Do("hset", "user:"+u.login, "pass", pass)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetName() (string, error) {
	name, err := redis.String(u.conn.Do("hget", "user:"+u.login, "name"))
	if err != nil {
		return "", err
	}
	return name, nil
}

func (u *User) SetName(name string) error {
	_, err := u.conn.Do("hset", "user:"+u.login, "name", name)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) CreateSession(id string, ex time.Duration) (*Session, error) {
	//	xid := xid.New().String()
	if ex == 0 {
		ex = defaultSessionLifetime
	}
	_, err := u.conn.Do("multi")
	if err != nil {
		return nil, err
	}
	_, err = u.conn.Do("hmset", "sess:"+id, "user", u.login)
	if err != nil {
		return nil, err
	}
	_, err = u.conn.Do("expire", ex/time.Second)
	if err != nil {
		return nil, err
	}
	_, err = u.conn.Do("exec")
	if err != nil {
		return nil, err
	}
	session := &Session{
		id:   "sess:" + id,
		conn: u.conn,
	}
	return session, nil
}

type Session struct {
	id   string
	conn redis.Conn
}

func (s *Session) ProlongSession(ex time.Duration) error {
	if ex == 0 {
		ex = defaultSessionLifetime
	}
	ok, err := redis.Bool(s.conn.Do("expire", "sess:"+s.id, ex/time.Second))
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("redis cannot set expiration on session" + s.id)
	}
	return nil
}

func (s *Session) Delete() error {
	val, err := redis.Bool(s.conn.Do("del", "sess:"+s.id))
	if err != nil {
		return err
	}
	if !val {
		return ErrNotFound
	}
	return nil
}

func (s *Session) PutString(key, value string) error {
	_, err := s.conn.Do("hset", "sess:"+s.id, key, value)
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) GetString(key string) (string, error) {
	val, err := redis.String(s.conn.Do("hget", "sess:"+s.id, key))
	if err != nil {
		return "", err
	}
	return val, nil
}

func (s *Session) GetUser() (*User, error) {
	userId, err := redis.String(s.conn.Do("hget", "sess:"+s.id, "user"))
	if err != nil {
		return nil, err
	}
	if userId == "" {
		return nil, ErrNotFound
	}
	user := &User{
		login: userId,
		conn:  s.conn,
	}

	return user, nil
}

type UsersStorage struct {
	conn redis.Conn
}

func NewUsersStorage() (*UsersStorage, error) {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	defer conn.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &UsersStorage{conn: conn}, nil

}

func (u *UsersStorage) CreateUser(login, name, pass string) (*User, error) {
	_, err := u.conn.Do("hmset", "user:"+login, "name", name, "pass", pass)
	if err != nil {
		return nil, err
	}
	user := &User{
		login: login,
		conn:  u.conn,
	}
	return user, nil
}

func (u *UsersStorage) FindSessionById(id string) (*Session, error) {
	exists, err := redis.Bool(u.conn.Do("exists", "sess:"+id))
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}
	session := &Session{
		id:   id,
		conn: u.conn,
	}
	return session, nil
}
