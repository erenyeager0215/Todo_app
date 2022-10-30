package models

import (
	"database/sql"
	"log"
	"time"
)

/*
	Userに関わる処理
*/

type User struct {
	ID       int
	UUID     string
	Name     string
	Email    string
	PassWord string
	CreateAt time.Time
	//todosをUserの構造体についか
	Todos []Todo
}

type Session struct {
	ID        int
	UUID      string
	Email     string
	UserID    int
	CreatedAt time.Time
}

func (u *User) CreateUser() (err error) {
	cmd := `insert into users(
		uuid,
		name,
		email,
		password,
		created_at) values($1,$2,$3,$4,$5)`

	_, err = Db.Exec(cmd,
		createUUID(),
		u.Name,
		u.Email,
		Encrypt(u.PassWord),
		time.Now())

	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func GetUser(id int) (user User, err error) {
	user = User{}
	cmd := `SELECT * FROM users WHERE id = $1`
	row := Db.QueryRow(cmd, id)
	err = row.Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PassWord,
		&user.CreateAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No row")
		} else {
			log.Println(err)
		}
	}
	return user, err
}

func (u *User) UpdateUser() (err error) {
	cmd := `UPDATE users set name = $1, email = $2  WHERE id = $3`
	_, err = Db.Exec(cmd, u.Name, u.Email, u.ID)

	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func (u *User) DeleteUser() (err error) {
	cmd := `DELETE FROM users WHERE id = $1`
	_, err = Db.Exec(cmd, u.ID)
	if err != nil {
		log.Fatalln(err)
	}

	return err
}

// メアドを入力したらそれに紐づくユーザー情報を返す
func GetUserByEmail(email string) (user User, err error) {
	user = User{}
	cmd := `select * from users where email = $1`
	err = Db.QueryRow(cmd, email).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PassWord,
		&user.CreateAt)
	return user, err
}

// sessionの作成
func (u *User) CreateSession() (session Session, err error) {
	session = Session{}
	cmd1 := `insert into sessions (
		uuid,
		email,
		user_id,
		created_at) values($1,$2,$3,$4)`

	_, err = Db.Exec(cmd1, createUUID(), u.Email, u.ID, time.Now())
	if err != nil {
		log.Println(err)
	}

	cmd2 := `select * from sessions where user_id = $1 and email = $2`

	err = Db.QueryRow(cmd2, u.ID, u.Email).Scan(&session.ID,
		&session.UUID,
		&session.Email,
		&session.UserID,
		&session.CreatedAt)

	return session, err
}

// sessionがあるかの判定
func (sess *Session) CheckSession() (valid bool, err error) {
	cmd := `select * from sessions where uuid = $1 `
	err = Db.QueryRow(cmd, sess.UUID).Scan(
		&sess.ID,
		&sess.UUID,
		&sess.Email,
		&sess.UserID,
		&sess.CreatedAt)

	//エラーがあればsessionは存在しない処理にして返す
	if err != nil {
		valid = false
		return
	}

	//もしsessionIDが０でなければ（sessionテーブルへ登録があればIDに数値が付与される）
	if sess.ID != 0 {
		valid = true
	}
	return valid, err
}

func (sess *Session) DeleteSessionByUUID() (err error) {
	cmd := `delete from sessions where uuid = $1`
	_, err = Db.Exec(cmd, sess.UUID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func (sess *Session) GetUserBySession() (user User, err error) {
	user = User{}
	cmd := `select * From users where id = $1`
	err = Db.QueryRow(cmd, sess.UserID).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PassWord,
		&user.CreateAt)

	return user, err
}
