package rlock

import (
	"context"
	"fmt"
	"log"
	"time"

	retry "github.com/avast/retry-go"
	v8 "github.com/go-redis/redis/v8"
)

type RLock struct {
	Key      string
	OwnerID  string
	Client   *v8.Client
	Ttl      time.Duration
	Deadline time.Duration
}

func (r *RLock) Lock() error {
	c, cancel := context.WithTimeout(context.Background(), r.Deadline)
	defer cancel()
	return retry.Do(
		r.LockOnce,
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			return 1 * time.Second
		}),
		retry.Attempts(100),
		retry.Context(c),
	)
}

// func lockOnce
func (r *RLock) LockOnce() error {

	args := []interface{}{
		r.OwnerID,
		ToMilliseconds(r.Ttl),
	}
	luaRS, err := r.Client.Eval(context.Background(), luaMutexLock, []string{r.Key}, args...).Result()
	if err != nil {
		return err
	}
	// parse res
	rs, ok := luaRS.([]interface{})
	if !ok {
		return fmt.Errorf("failed to parse luaRS")
	}
	isSuccess, ok := rs[0].(int64)
	if !ok {
		return fmt.Errorf("failed to parse luaRS")
	}
	if isSuccess == 1 {
		log.Printf("key: %s ownerid: %s get lock success\n", r.Key, r.OwnerID)
		return nil
	}
	log.Printf("key: %s ownerid: %s get lock failed\n", r.Key, r.OwnerID)
	return fmt.Errorf("failed to lock")
}

func (r *RLock) Unlock() error {
	luaRS, err := r.Client.Eval(context.Background(), luaKvDelIfExists, []string{r.Key}, r.OwnerID).Result()
	if err != nil {
		return err
	}
	isSuccess, ok := luaRS.(int64)
	if !ok {
		return fmt.Errorf("failed to parse luaRS")
	}
	if isSuccess == 1 {
		log.Printf("key: %s ownerid: %s unlock\n", r.Key, r.OwnerID)
		return nil
	}
	return fmt.Errorf("failed to unlock")
}

func ToMilliseconds(dur time.Duration) uint64 {
	return uint64(dur / time.Millisecond)
}
