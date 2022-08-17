package xutils

import (
	"context"
	"time"
)
import (
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

type XRedisClient struct {
	client        *goredislib.Client
	clusterClient *goredislib.ClusterClient
	redisMutex    *redsync.Mutex
	bCluster      bool
}

func NewXRedisClient(addr string, dbIndex int, pwd string, lockName string) (*XRedisClient, error) {
	xClient := &XRedisClient{}
	xClient.client = goredislib.NewClient(&goredislib.Options{Addr: addr, DB: dbIndex, Password: pwd})
	_, err := xClient.client.Ping(context.Background()).Result()
	if nil != err {
		return nil, err
	}
	if len(lockName) > 0 {
		pool := goredis.NewPool(xClient.client) // or, pool := redigo.NewPool(...)

		// Create an instance of redisync to be used to obtain a mutual exclusion
		// lock.
		rs := redsync.New(pool)

		// Obtain a new mutex by using the same name for all instances wanting the
		// same lock.
		xClient.redisMutex = rs.NewMutex(lockName)
	}

	return xClient, nil
}

func NewXRedisClusterClient(addrs []string, pwd string, lockName string) (*XRedisClient, error) {
	xClient := &XRedisClient{}
	xClient.bCluster = true
	xClient.clusterClient = goredislib.NewClusterClient(&goredislib.ClusterOptions{Addrs: addrs, Password: pwd})
	_, err := xClient.client.Ping(context.Background()).Result()
	if nil != err {
		return nil, err
	}
	if len(lockName) > 0 {
		pool := goredis.NewPool(xClient.client) // or, pool := redigo.NewPool(...)

		// Create an instance of redisync to be used to obtain a mutual exclusion
		// lock.
		rs := redsync.New(pool)

		// Obtain a new mutex by using the same name for all instances wanting the
		// same lock.
		xClient.redisMutex = rs.NewMutex(lockName)
	}

	return xClient, nil
}

func (x *XRedisClient) Lock(timeOut time.Duration) bool {
	if nil == x.redisMutex {
		return false
	}

	if timeOut > 0 {
		timer := time.NewTimer(timeOut)
		ticker := time.NewTicker(time.Millisecond)
		defer timer.Stop()
		defer ticker.Stop()
		for {
			select {
			case <-timer.C:
				return false

			case <-ticker.C:
				{
					if err := x.redisMutex.Lock(); nil == err {
						return true
					}
				}
			}
		}
	}
	if err := x.redisMutex.Lock(); nil == err {
		return true
	}
	return false
}

func (x *XRedisClient) UnLock(timeOut time.Duration) bool {
	if nil == x.redisMutex {
		return false
	}

	if timeOut > 0 {
		timer := time.NewTimer(timeOut)
		ticker := time.NewTicker(time.Millisecond)
		defer timer.Stop()
		defer ticker.Stop()
		for {
			select {
			case <-timer.C:
				return false

			case <-ticker.C:
				{
					if ok, err := x.redisMutex.Unlock(); ok || nil == err {
						return true
					}
				}
			}
		}
	}
	if ok, err := x.redisMutex.Unlock(); ok || nil == err {
		return true
	}
	return false
}

func (x *XRedisClient) Redis() goredislib.UniversalClient {
	if x.bCluster {
		return x.clusterClient
	}
	return x.client
}
