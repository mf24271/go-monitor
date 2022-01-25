# go-monitor

go monitor

## Usage

```go
package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/mf24271/go-monitor/monitor"
)

func calc(i int32) int32 {
	if i <= 0 {
		return i
	}
	return i + calc(i-1)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	conf := monitor.NewConfig()
	m := monitor.NewMonitor(conf)
	err := m.Start()
	if err != nil {
		panic(err)
	}
	//Simulate computing intensive tasks
	go func() {
		for {
			if rand.Int31n(20) == 0 {
				calc(rand.Int31n(10000000) + 1000)
				calc(rand.Int31n(10000000) + 1000)
				calc(rand.Int31n(10000000) + 1000)
				calc(rand.Int31n(10000000) + 1000)
				calc(rand.Int31n(10000000) + 1000)
			}
			time.Sleep(time.Duration(rand.Int31n(1100)) * time.Millisecond)
		}
	}()

	wg.Wait()
}
```
