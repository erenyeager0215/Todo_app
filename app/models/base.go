package models

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"
	"os"
	"todo_app/config"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

/*
	tebleの作成
*/

var Db *sql.DB

var err error

func init() {
	//herokuの環境変数の値を取り出すことができる
	//DATABASE＿URLはherokuのpostgreSQLのURLを表す
	url := os.Getenv("DATABASE_URL")
	//urlで取得したリソースをコネクションとして取得する
	connection, _ := pq.ParseURL(url)
	connection += "sslmode=require"
	Db, err = sql.Open(config.Config.SQLDriver, connection)
	if err != nil {
		log.Fatalln(err)
	}
}

// UUIDの生成
func createUUID() (uuidobj uuid.UUID) {
	uuidobj, _ = uuid.NewUUID()
	return
}

// passwordの保存はハッシュ値にする
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}
