package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type KeyVal struct {
	key string
	value string
	expiredDayTime *time.Time //初期値nil
}

type SetCommand struct {}

var (
	setInstance Command
	keyVals []KeyVal
)

func NewSetInstance() Command {
	if setInstance == nil {
		setInstance = &SetCommand{}
	}
	return setInstance
}

func (this *SetCommand)Exec(conn net.Conn, args []string) []byte {

	keyVal := KeyVal{key: args[0], value: args[1]}

	for _, v := range args {
		fmt.Println(v)
	}

	if len(args) == 5 {
		op := args[2]
		fmt.Println(op)
		opVal, err := strconv.Atoi(args[3])
		if err != nil {
			fmt.Println("PXオプションはmillisecondを数値で入力してください")
		}

		if op == "px" {
			expiredTime := time.Now().Add(time.Duration(opVal) * time.Millisecond)
			keyVal.expiredDayTime = &expiredTime
		}
	}

	keyVals = append(keyVals, keyVal)
	return []byte(fmt.Sprint("$", 2, CRLF, "OK", CRLF))
}
