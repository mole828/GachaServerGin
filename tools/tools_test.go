package tools

import (
	"context"
	"testing"
	"time"
)

func TestLo(t *testing.T) {
	begin := time.Unix(0, 0)
	now := time.Now()
	t.Log(now.Sub(begin) > 0)
	t.Log(begin.Sub(now) > 0)
}

type person struct {
	Name string
}

func TestString(t *testing.T) {
	var s string
	t.Logf("\"%s\"", s)

	t.Logf("\"%s\"", person{}.Name)
}

func TestGo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		time.Sleep(time.Second * 3)
		cancel()
		time.Sleep(time.Second * 1)
	}()
	ch := make(chan int, 1)
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				t.Log("job timeout return")
				return
			case <-ch:
				t.Log("ch job work")
				time.Sleep(time.Second)
			default:
				t.Log("job still working")
				ch <- 1
				time.Sleep(1 * time.Second)
			}
		}
	}(ctx)
}
