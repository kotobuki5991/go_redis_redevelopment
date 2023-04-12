package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
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

		redisRequest := getCmdAndArg(input)
		writeResponse(conn, redisRequest)
	}
}

type RedisRequest struct {
	command string
	args []string
}

// RESP arrayのフォーマットからコマンドを解析する
// *2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n
// *の次の数値がRESP arrayの要素数
// 1つ目$の次の数値が1つ目の値の文字数（ここでは$4なのでECHOの4文字）
// 2つ目の$の次の数値が2つ目の値の文字数（ここでは$3なのでheyの3文字）
func getCmdAndArg(input []byte) RedisRequest{
	inputRespAry := strings.Split(string(input), "\r\n")

	respAryLength, err := strconv.Atoi(strings.ReplaceAll(inputRespAry[0], "*", ""))
	if err != nil {
		fmt.Println("RESP Array format invalid", err.Error())
	}
	fmt.Println(respAryLength)

	command := ""
	args := make([]string, 0)
	for i := 1; i < respAryLength; i++ {
		if strings.Index(inputRespAry[0], "$") != -1 {continue}

		if i == 2 {
			// コマンドを取得
			command = inputRespAry[i]
			continue
		}

		// コマンドの引数を追加
		args = append(args, inputRespAry[i])
	}
	fmt.Println(command)
	return RedisRequest{
		command,
		args,
	}
}



func writeResponse(conn net.Conn, redisRequest RedisRequest) {
	fmt.Println("writeResponse")
	cmd := redisRequest.command
	args := redisRequest.args
	// コネクションにデータを書き込む
	switch cmd {
	case "echo":
		conn.Write([]byte(args[0]))
	case "ping":
		conn.Write([]byte("+PONG\r\n"))
	default:
		fmt.Println("command invalid. your command is ", cmd)
	}
}
