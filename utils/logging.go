package utils

import (
	"io"
	"log"
	"os"
)

func LoggingSettings(logFile string) {
	//logFileを読み込む
	//logFileという読み書き可能なファイルを開く
	logfile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}

	//io.MultiWriterで引数にログの書き込み先を指定
	multiLogFile := io.MultiWriter(os.Stdout, logfile)

	//ログのフォーマットを指定
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	//ログの出力先を指定
	log.SetOutput(multiLogFile)
}

/*
	log.Print("ログ情報")
	のようにlogパッケージの情報出力先をmain関数実行前にする
	LogginSettingsはlogfailと
	標準出力に指定のフォーマットで出力するための処理
*/
