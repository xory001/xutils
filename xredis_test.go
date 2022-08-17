package xutils

import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"
)

func StartRedSync() {
	log.Println("StartRedSync start")
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go redisTest(wg)
	go redisTest2(wg)

	wg.Wait()
	log.Println("StartRedSync end")
}

func redisTest(wg *sync.WaitGroup) {
	client, err := NewXRedisClient("192.168.0.42:6379", 0, "123456", "redis_lock")
	if nil != err {
		return
	}
	for i := 0; i < 100; i++ {
		ctx := context.Background()
		if client.Lock(time.Millisecond * 10) {
			rv, err := client.Redis().HGetAll(ctx, "name").Result()
			if nil != err {
				rv[strconv.Itoa(i)] = strconv.Itoa(i) + "_" + "a"
				client.Redis().HSet(context.Background(), "name", "xory", time.Minute*3)
			}

			client.UnLock(time.Millisecond * 10)
		} else {
			log.Println("redisTest, lock err")
		}
	}
	wg.Done()
}
func redisTest2(wg *sync.WaitGroup) {
	client, err := NewXRedisClient("192.168.0.42:6379", 0, "123456", "redis_lock")
	if nil != err {
		return
	}
	for i := 0; i < 100; i++ {
		if client.Lock(time.Millisecond * 10) {
			client.Redis().Set(context.Background(), "name", "zhu", time.Minute*3)
			client.UnLock(time.Millisecond * 10)
		} else {
			log.Println("redisTest2, lock err")
		}
	}
	wg.Done()
}
