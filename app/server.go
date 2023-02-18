package main

import (
	"fmt"
	"io"
	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// リッスンの開始
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		// プログラムを終了させる（0は成功、0以外はエラーを示す）
		os.Exit(1)
	}

	receiveTCPConnection(l)
}

func receiveTCPConnection(l net.Listener){
	// net.listenで得たリスナーへの接続を待機し、返す。
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	// 接続を閉じる。receiveTCPConnection関数の終了時に実行される。
	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		// コネクションからデータを読み取る
		_, err := conn.Read(buf)
		if err == io.EOF {
			fmt.Println("Connection closed")
			break
		}
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		redisHandler(conn)
	}
}

func redisHandler(conn net.Conn) {
	fmt.Println("call redisHandler")
	// コネクションにデータを書き込む
	conn.Write([]byte("+PONG\r\n"))
}
