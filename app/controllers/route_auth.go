package controllers

import (
	"fmt"
	"log"
	"net/http"
	"todo_app/app/models"
)

func signup(w http.ResponseWriter, r *http.Request) {
	//getメソッドの時の処理とpostメソッドの時の処理を分ける
	if r.Method == "GET" {
		_, err := session(w, r)
		if err != nil {
			generateHTML(w, nil, "layout", "public_navbar", "signup")
		} else {
			http.Redirect(w, r, "/todos", 302)
		}
	} else if r.Method == "POST" {
		//inputタグの入力フォームの値を取得する
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		user := models.User{
			//inputタグのname属性で指定した値を取得
			Name:     r.PostFormValue("name"),
			Email:    r.PostFormValue("email"),
			PassWord: r.PostFormValue("password"),
		}
		//CreateUserでDBに登録
		if err := user.CreateUser(); err != nil {
			log.Println(err)
		}
		//redirectでpost後はトップページへ遷移するように指示
		http.Redirect(w, r, "/", 302)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		generateHTML(w, nil, "layout", "public_navbar", "login")
	} else {
		http.Redirect(w, r, "/todos", 302)
	}
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	//入力されたメールアドレスからユーザ情報を取得
	user, err := models.GetUserByEmail(r.PostFormValue("email"))

	//エラーの場合リダイレクトする
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login", 302)
	}

	//入力したパスワードがユーザ情報と一致した場合は、セッションを作成する
	if user.PassWord == models.Encrypt(r.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			log.Println(err)
		}

		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.UUID,
			HttpOnly: true,
		}

		//クッキーにcookieポインタをセット
		http.SetCookie(w, &cookie)

		//ログイン成功後にリダイレクトするページを指定
		http.Redirect(w, r, "/", 302)
	} else {
		http.Redirect(w, r, "/login", 302)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err != nil {
		log.Println(err)
	}
	if err != http.ErrNoCookie {
		session := models.Session{UUID: cookie.Value}
		session.DeleteSessionByUUID()
	}
	http.Redirect(w, r, "/login", 302)
}
