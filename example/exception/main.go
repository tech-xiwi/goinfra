package main

import (
	"fmt"
	"math/rand"
	"runtime/debug"
	"time"

	"github.com/tech-xiwi/goinfra/exception"
)

var sRecover = func(params ...interface{}) {
	if r := recover(); r != nil {
		fmt.Printf("[sRecover]:%s\n[params]:%v\n[stack]:\n%s\n", r, params, debug.Stack())
	}
}

var randSeed = rand.New(rand.NewSource(time.Now().Local().Unix()))

func main() {
	job := func(params ...interface{}) {
		fmt.Println("start working ...")
		fmt.Println("doing some work", params)
		time.Sleep(time.Duration(randSeed.Uint64()%4) * time.Second)
		now := time.Now().UTC().Unix()

		if now%2 == 1 {
			panic("panic") // can catch
			// log.Fatal("fault") // not catch , the process has exit(2)

			// ss := make([]int, 10)
			// s := ss[:11] // can catch
			// fmt.Println(s)
		}

		fmt.Println("finish work ...")

	}

	b := exception.Block{
		Try:           job,
		ReTry:         job,
		MaxRetryCount: 3,
		Catch: func(stack []byte, params ...interface{}) {
			fmt.Printf("[catch]:%s\n[stack]:\n%s\n", params[0], stack)
		},
		Finally: func(params ...interface{}) {
			defer sRecover(params...)
			fmt.Println("[finally]")
		},
	}
	b.Do("hard-job", 1, 2, 3, 4, 5)
	fmt.Println("next logic", b.MaxRetryCount, b.RetryCount)
}
