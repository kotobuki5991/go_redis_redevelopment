package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var CRLF = "\r\n"

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

// RESP arrayのフォーマットからコマンドと引数を解析する
// *2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n
// *の次の数値がRESP arrayの要素数
// 1つ目$の次の数値が1つ目の値の文字数（ここでは$4なのでECHOの4文字）
// 2つ目の$の次の数値が2つ目の値の文字数（ここでは$3なのでheyの3文字）
func getCmdAndArg(input []byte) RedisRequest{
	inputRespAry := strings.Split(string(input), CRLF)
	respAryLength := len(inputRespAry)

	command := ""
	args := make([]string, 0)
	for i := 1; i < respAryLength; i++ {
		fmt.Println(i)
		fmt.Println(inputRespAry[i])
		if strings.Index(inputRespAry[i], "$") != -1 {continue}

		if i == 2 {
			// コマンドを取得
			command = inputRespAry[i]
			continue
		}

		// コマンドの引数を追加
		args = append(args, inputRespAry[i])
	}
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
		echoCmdHandler(conn, args[0])
	case "ping":
		pingCmdHandler(conn)
	case "set":
		setCmdHandler(conn, KeyVal{key: args[0], value: args[1]})
	case "get":
		searchKey := args[0]
		getCmdHandler(conn, searchKey)
	default:
		fmt.Println("command invalid. your command is ", cmd)
	}
}

func echoCmdHandler(conn net.Conn, args string){
	conn.Write([]byte(fmt.Sprint("$", len(args), CRLF, args, CRLF)))
}

func pingCmdHandler(conn net.Conn){
	conn.Write([]byte("+PONG\r\n"))
}

type KeyVal struct {
	key string
	value string
}

var keyVals []KeyVal

func setCmdHandler(conn net.Conn, keyVal KeyVal){
	keyVals = append(keyVals, keyVal)
	conn.Write([]byte(fmt.Sprint("$", 2, CRLF, "OK", CRLF)))
}

func getCmdHandler(conn net.Conn, key string){
	resp := findValueByKey(key)
	conn.Write([]byte(fmt.Sprint("$", len(resp), CRLF, resp, CRLF)))
}

func findValueByKey(key string) string {
	for _, elem := range keyVals {
			if elem.key == key {
					return elem.value
			}
	}
	return "OK" // keyが見つからなかった場合
}
