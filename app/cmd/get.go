package cmd

import (
	"fmt"
	"time"
	"github.com/kotobuki5991/go_redis_redevelopment/app/consts"
)

type GetCommand struct {}

var getInstance Command

func NewGetInstance() Command {
	if getInstance == nil {
		getInstance = &GetCommand{}
	}
	return getInstance
}

func (cmd *GetCommand)Exec(args []string) []byte {
	searchKey := args[0]
	resp := cmd.findValueByKey(searchKey)
	if (resp == nil){
		return []byte("$-1\r\n")
	}
	return []byte(fmt.Sprint("$", len(*resp), consts.CRLF, *resp, consts.CRLF))
}

func (cmd *GetCommand)findValueByKey(key string) *string {
	for _, elem := range keyVals {
			if elem.key == key {
				now := time.Now()
				if elem.expiredDayTime != nil && now.After(*elem.expiredDayTime) {
					fmt.Println("Expierd")
					return nil // null bulk stringを返す
				}
				fmt.Println(elem)
				return &elem.value
			}
	}
	return nil // keyが見つからなかった場合
}
