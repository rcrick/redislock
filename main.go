package main

import (
	"example/rlock"
	"fmt"
	"strconv"
	"sync"
	"time"

	v8 "github.com/go-redis/redis/v8"
)

func main() {

	client := v8.NewClient(&v8.Options{Addr: "127.0.0.1:6379"})
	wg := sync.WaitGroup{}
	
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(ownerid int) {
			lock := rlock.RLock{
				Key:      "1",
				OwnerID:  strconv.Itoa(ownerid),
				Client:   client,
				Ttl:      20 * time.Second,
				Deadline: 100 * time.Second,
			}
			lock.Lock()
			time.Sleep(time.Second)
			err := lock.Unlock()
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
