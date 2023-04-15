package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
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
		fmt.Println("--------------")
		fmt.Println(inputRespAry[i])
		fmt.Println("--------------")
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
	fmt.Println("command")
	fmt.Println(cmd)
	fmt.Println("command")
	args := redisRequest.args
	// コネクションにデータを書き込む
	resp := make([]byte, 0)
	switch cmd {
	case "echo":
		resp = append(resp, echoCmdHandler(conn, args[0])...)
	case "ping":
		resp = append(resp, pingCmdHandler(conn)...)
	case "set":
		resp = append(resp, setCmdHandler(conn, args)...)
	case "get":
		searchKey := args[0]
		resp = append(resp, getCmdHandler(conn, searchKey)...)
	default:
		fmt.Println("command invalid. your command is ", cmd)
	}
	conn.Write(resp)
	fmt.Println("=======================")
}

func echoCmdHandler(conn net.Conn, args string) []byte{
	return []byte(fmt.Sprint("$", len(args), CRLF, args, CRLF))
}

func pingCmdHandler(conn net.Conn) []byte{
	return []byte("+PONG\r\n")
}

type KeyVal struct {
	key string
	value string
	expiredDayTime *time.Time //初期値nil
}

var keyVals []KeyVal

func setCmdHandler(conn net.Conn, args []string) []byte{
	keyVal := KeyVal{key: args[0], value: args[1]}

	for _, v := range args {
		fmt.Println(v)
	}


	if len(args) == 4 {
		op := args[2]
		opVal, err := strconv.Atoi(args[3])
		if err != nil {
			fmt.Println("PXオプションはmillisecondを数値で入力してください")
			fmt.Println(args[3])
		}

		if op == "px" {
			expiredTime := time.Now().Add(time.Duration(opVal) * time.Millisecond)
			keyVal.expiredDayTime = &expiredTime
		}
	}

	keyVals = append(keyVals, keyVal)
	return []byte(fmt.Sprint("$", 2, CRLF, "OK", CRLF))
}

func getCmdHandler(conn net.Conn, key string) []byte{
	resp := findValueByKey(key)
	fmt.Println("resp")
	fmt.Println(resp)
	fmt.Println("resp")
	return []byte(fmt.Sprint("$", len(resp), CRLF, resp, CRLF))
}

func findValueByKey(key string) string {
	for _, elem := range keyVals {
			if elem.key == key {
				now := time.Now()
				fmt.Println("check isExpierd")
				if elem.expiredDayTime != nil && now.After(*elem.expiredDayTime) {
					fmt.Println("Expierd")
					return "$-1\r\n" // null bulk stringを返す
				}
				fmt.Println("not Expierd")
				fmt.Println(elem)
				return elem.value
			}
	}
	return "" // keyが見つからなかった場合
}
