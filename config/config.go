package config

import (
	"log"
	"todo_app/utils"

	"gopkg.in/go-ini/ini.v1"
)

type ConfigList struct {
	Port      string
	SQLDriver string
	DbName    string
	LogFile   string
	Static    string //静的ファイルがあるディレクトリを指定する
}

// 外部から呼び出されるように大文字で変数宣言しておく
var Config ConfigList

// main関数より前にconfig情報を読み込む
func init() {
	LoadConfig()
	utils.LoggingSettings(Config.LogFile)
}

// configファイルから構造体へデータを入れる
func LoadConfig() {

	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln(err)
	}
	Config = ConfigList{
		Port:      cfg.Section("web").Key("port").MustString("8080"),
		SQLDriver: cfg.Section("db").Key("driver").MustString("postgres"),
		DbName:    cfg.Section("name").Key("name").MustString("webapp.sql"),
		LogFile:   cfg.Section("web").Key("logfile").MustString("webapp.log"),
		Static:    cfg.Section("web").Key("static").String(),
	}

}
