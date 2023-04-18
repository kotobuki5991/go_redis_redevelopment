package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"

	"app/cmd"
	"app/consts"
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

		// ゴルーチンプールの設定
		poolSize := 5
		pool := make(chan struct{}, poolSize)
		for i := 0; i < poolSize; i++ {
			pool <- struct{}{}
		}

	// ゴルーチンプールで接続を処理する
	var wg sync.WaitGroup
	for {
		// goroutineプールに空きがなければ待機する
		<-pool
		wg.Add(1)
		go func() {
			createRedisRequestReceiver(l)
			wg.Done()
			// 処理が完了したらgoroutineをプールに戻す
			pool <- struct{}{}
		}()
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
	inputRespAry := strings.Split(string(input), consts.CRLF)
	respAryLength := len(inputRespAry)

	command := ""
	args := make([]string, 0)
	for i := 1; i < respAryLength; i++ {
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
	cmd := redisRequest.command
	args := redisRequest.args
	// コネクションにデータを書き込む
	resp := make([]byte, 0)
	// コマンド名に応じたインスタンスを取得
	cmdInstance := getCmdInstance(cmd)
	resp = append(resp, cmdInstance.Exec(conn, args)...)
	conn.Write(resp)
}

func getCmdInstance(cmdName string) cmd.Command {
	var command cmd.Command
	switch cmdName {
	case "echo":
		command = cmd.NewEchoInstance()
	case "ping":
		command = cmd.NewPingInstance()
	case "set":
		command = cmd.NewSetInstance()
	case "get":
		command = cmd.NewGetInstance()
	default:
		fmt.Println("command invalid. your command is ", cmdName)
	}
	return command
}

// func echoCmdHandler(conn net.Conn, args string) []byte{
// 	return []byte(fmt.Sprint("$", len(args), consts.CRLF, args, consts.CRLF))
// }

// func pingCmdHandler(conn net.Conn) []byte{
// 	return []byte("+PONG\r\n")
// }

// type KeyVal struct {
// 	key string
// 	value string
// 	expiredDayTime *time.Time //初期値nil
// }

// var keyVals []KeyVal

// func setCmdHandler(conn net.Conn, args []string) []byte{
// 	keyVal := KeyVal{key: args[0], value: args[1]}

// 	for _, v := range args {
// 		fmt.Println(v)
// 	}

// 	if len(args) == 5 {
// 		op := args[2]
// 		fmt.Println(op)
// 		opVal, err := strconv.Atoi(args[3])
// 		if err != nil {
// 			fmt.Println("PXオプションはmillisecondを数値で入力してください")
// 		}

// 		if op == "px" {
// 			expiredTime := time.Now().Add(time.Duration(opVal) * time.Millisecond)
// 			keyVal.expiredDayTime = &expiredTime
// 		}
// 	}

// 	keyVals = append(keyVals, keyVal)
// 	return []byte(fmt.Sprint("$", 2, consts.CRLF, "OK", consts.CRLF))
// }

// func getCmdHandler(conn net.Conn, key string) []byte{
// 	resp := findValueByKey(key)
// 	if (resp == nil){
// 		return []byte("$-1\r\n")
// 	}
// 	return []byte(fmt.Sprint("$", len(*resp), consts.CRLF, *resp, consts.CRLF))
// }

// func findValueByKey(key string) *string {
// 	for _, elem := range keyVals {
// 			if elem.key == key {
// 				now := time.Now()
// 				if elem.expiredDayTime != nil && now.After(*elem.expiredDayTime) {
// 					fmt.Println("Expierd")
// 					return nil // null bulk stringを返す
// 				}
// 				fmt.Println(elem)
// 				return &elem.value
// 			}
// 	}
// 	return nil // keyが見つからなかった場合
// }
