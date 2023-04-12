package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
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

		checkCmd(input)
		writeResponse(conn)
	}
}

// RESP arrayのフォーマットからコマンドを解析する
// *2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n
// *の次の数値がRESP arrayの要素数
// 1つ目$の次の数値が1つ目の値の文字数（ここでは$4なのでECHOの4文字）
// 2つ目の$の次の数値が2つ目の値の文字数（ここでは$3なのでheyの3文字）
func checkCmd(input []byte){
	respAryLength, err := strconv.Atoi(string(input[1:2]))
	fmt.Println(input[1:2])
	if err != nil {
		fmt.Println("RESP Array format invalid", err.Error())
	}
	fmt.Println(respAryLength)
	cmdLength, err := strconv.Atoi(string(input[3:4]))
	fmt.Println(input[3:4])
	if err != nil {
		fmt.Println("RESP Array format invalid", err.Error())
	}
	fmt.Println(cmdLength)
}


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
