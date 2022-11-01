package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"todo_app/app/models"
	"todo_app/config"
)

func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}
	//Mustはテンプレートをあらかじめキャッシュしておき効率的な処理する
	//Must引数にパースファイルを渡すとエラーの場合はパニックになる
	//ParseFilesによって生成されたtemplatesからExecuteTemplateを実行する場合は、表示するファイルを明示的に指定する必要があります。
	//ここで指定したファイルにデータを渡し、テンプレートを表示します。
	templates := template.Must(template.ParseFiles(files...))
	//第一引数にResponseWriterの引数
	//第二引数は実行するテンプレート名(defineで指定したテンプレート名)
	//defineを使ったファイルを読み込む時はExecuteTemplateを使う
	templates.ExecuteTemplate(w, "layout", data)
}

// server.goで作成されるcookieがDBにあるか確認する
func session(w http.ResponseWriter, r *http.Request) (sess models.Session, err error) {
	//server.goで作成されるcookieの構造体のNameを指定
	//これにより、cookieを取得できる
	cookie, err := r.Cookie("_cookie")
	//もしエラーがない場合（成功した場合）
	if err == nil {
		//sessに取得したUUIDをもつSessionの構造体をいれる
		sess = models.Session{UUID: cookie.Value}
		//CheckSessionでセッションがDBにあるかどうかをチェックする
		if ok, _ := sess.CheckSession(); !ok {
			err = fmt.Errorf("Invalid session")
		}
	}
	return sess, err
}

// 「 ^ 」は直後の文字が行の先頭にある場合にマッチします。
// 「 $ 」直前の文字が行の末尾にある場合にマッチします。
var validPath = regexp.MustCompile("^/todos/(edit|update|delete)/([0-9]+)$")

// ハンドラ関数を引数とし、ハンドラ関数を戻り値とする関数
func parseURL(fn func(http.ResponseWriter, *http.Request, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//validPathとリクエストURLのパスでマッチしたものをスライスで返す
		q := validPath.FindStringSubmatch(r.URL.Path)
		if q == nil {
			http.NotFound(w, r)
			return
		}
		//qに格納したインデックス番号2は数字だが、文字列型なのでint型に変換
		qi, err := strconv.Atoi(q[2])
		if err != nil {
			http.NotFound(w, r)
			return
		}
		//prseURLの引数に指定した関数を実行
		fn(w, r, qi)
	}
}

// サーバの立ち上げの処理
func StartMainServer() error {
	//cssとjsファイル読み込む
	//捜索するディレクトリをFileServerへ指定する
	files := http.FileServer(http.Dir(config.Config.Static))
	//Dir()に入れたディレクトリ➡Handleの第1引数のパスを捜索する
	//http.StripPrefixは、第1引数に指定したパスを、
	//http.FileServer()が捜索するURLから取り除きます。
	http.Handle("/static/", http.StripPrefix("/static/", files))
	//HundleFuncでDefaultServerMuxにhandler関数を登録する
	http.HandleFunc("/", top)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/authenticate", authenticate)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/todos", index)
	http.HandleFunc("/todos/new", todoNew)
	//post系のハンドラ関数はformタグのactionに記載したurlと統一させる
	http.HandleFunc("/todos/save", todoSave)
	//urlの末尾をスラッシュにすることで要求されたurlの先頭が登録されたものと一致していればハンドラ関数が渡される
	http.HandleFunc("/todos/edit/", parseURL(todoEdit))
	http.HandleFunc("/todos/update/", parseURL(todoUpdate))
	http.HandleFunc("/todos/delete/", parseURL(todoDelete))

	//heroku対応
	//herokuからportを取得
	port := os.Getenv("PORT")
	return http.ListenAndServe(":"+port, nil)

}
