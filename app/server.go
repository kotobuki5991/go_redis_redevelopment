package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	// リッスンの開始
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		// プログラムを終了させる（0は成功、0以外はエラーを示す）
		os.Exit(1)
	}

	for {
		go createRedisRequestReceiver(l)
	}
}

func createRedisRequestReceiver(l net.Listener){
	// net.listenで得たリスナーへの接続を待機し、返す。
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	// 接続を閉じる。createRedisRequestReceiver関数の終了時に実行される。
	defer conn.Close()
	redisHandler(conn)
}


func redisHandler(conn net.Conn){
	input := make([]byte, 1024)
	for {
		// コネクションからデータを読み取る
		_, err := conn.Read(input)
		if err == io.EOF {
			fmt.Println("Connection closed")
			break
		}
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
	}
	writeResponse(conn)
}

// func redisHandler(conn net.Conn){
// 	fmt.Println("call redisHandler")

// 	buf := make([]byte, 1024)
// 	for {
// 		// コネクションからデータを読み取る
// 		_, err := conn.Read(buf)
// 		if err == io.EOF {
// 			fmt.Println("Connection closed")
// 			break
// 		}
// 		if err != nil {
// 			fmt.Println("Error reading:", err.Error())
// 		}
// 		writeResponse(conn)
// 	}

// 	fmt.Println(string(buf))
// }

// *2\r\n$4\r\nECHO\r\n$3\r\nhey\r\nを受け取った場合
func isEchoCmd(b []byte){
	artIndex := bytes.Index(b, []byte("\r\n"))
	// "\r\n"が存在しない場合処理終了
	if artIndex == -1 {return}
	// 区切り文字で分轄
	// slice := bytes.Split(b, []byte("\r\n"))
}

func writeResponse(conn net.Conn) {
	fmt.Println("writeResponse")
	// コネクションにデータを書き込む
	conn.Write([]byte("+PONG\r\n"))
}
