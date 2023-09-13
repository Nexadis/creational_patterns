package main

import (
	"fmt"
	"sync"
	"time"
)

type singleInstance struct{}

var (
	singleton *singleInstance
	once      sync.Once
)

func getSingleton() *singleInstance {
	once.Do(
		func() {
			fmt.Println("Initialize singleton")
			singleton = &singleInstance{}
		},
	)
	return singleton
}

func main() {
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Printf("Singleton: %p\n", getSingleton())
		}(i)
	}
	time.Sleep(100 * time.Millisecond)
}
