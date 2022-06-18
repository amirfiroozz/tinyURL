package middlewares

import (
	"fmt"
	"time"
	"tiny-url/pkg/config"
	"tiny-url/pkg/utils"
)

func Set(duration int, key string, val interface{}) *utils.Error {
	ctx := config.GetCtx()
	rdb := config.GetRdb()
	err := rdb.Set(ctx, key, val, time.Minute*time.Duration(duration)).Err()
	if err != nil {
		return &utils.Error{
			Code:   1,
			Status: 500,
			Msg:    err.Error(),
		}
	}
	fmt.Printf("set to redis memory %v:%v it will be removed in %v minutes\n", key, val, duration)
	return nil
}

func Get(key string) (*string, *utils.Error) {
	ctx := config.GetCtx()
	rdb := config.GetRdb()
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, &utils.Error{
			Code:   1,
			Status: 500,
			Msg:    err.Error(),
		}
	}
	fmt.Printf("get from redis memory %v:%v\n", key, val)
	return &val, nil

}
func Remove(key string) *utils.Error {
	ctx := config.GetCtx()
	rdb := config.GetRdb()
	val, err := rdb.Del(ctx, key).Result()
	if err != nil {
		return &utils.Error{
			Code:   1,
			Status: 500,
			Msg:    err.Error(),
		}
	}
	fmt.Printf("remove url form redis memory key:%v , count:%v\n", key, val)
	return nil

}

func IfExists(key string) bool {
	ctx := config.GetCtx()
	rdb := config.GetRdb()
	val, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("result: ", val)
	return val == 1
}
