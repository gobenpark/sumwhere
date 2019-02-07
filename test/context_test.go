package test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	bgContext := context.Background()
	newContext := context.WithValue(bgContext, "key", "value")

	t.Log(bgContext)
	t.Log(newContext)

	if newContext.Value("key") != "value" {
		t.Fail()
	}

	context.WithCancel(newContext)
	context.WithDeadline(bgContext, time.Now().AddDate(0, 0, 1))
	context.WithTimeout(bgContext, time.Second*3)
}

func TestContextUse(t *testing.T) {

	var currentUser string
	currentUser = "user"
	ctx := context.Background()

	ctx = context.WithValue(ctx, "current_user", currentUser)

	myFunc(t, ctx)
}

func myFunc(t *testing.T, ctx context.Context) error {
	var currentUser string

	if v := ctx.Value("current_user"); v != nil {
		u, ok := v.(string)
		if !ok {
			return errors.New("Not user")
		}

		currentUser = u
	} else {
		return errors.New("no have context value")
	}
	t.Log(currentUser)
	return nil
}

func longFunc() string {
	<-time.After(time.Second * 20)
	return "success"
}

func longFuncWithCtx(ctx context.Context) (string, error) {
	done := make(chan string)

	go func() {
		done <- longFunc()
	}()

	select {
	case result := <-done:
		return result, nil
	case <-ctx.Done():
		return "fail", ctx.Err()
	}
}

func TestMultiGoroutine(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		cancel()
	}()

	var jobCount = 10
	var wg sync.WaitGroup
	for i := 0; i < jobCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			result, err := longFuncWithCtx(ctx)
			if err != nil {
				t.Error(err)
			}
			t.Log(result)
		}()
	}
	wg.Wait()
}

func TestTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	go func() {
		// 고루틴을 종료해야 할 상황이 되면 cancel 함수 실행
		cancel()
	}()

	start := time.Now()
	result, _ := longFuncWithCtx(ctx)
	t.Logf("duration: %v result %s\n", time.Since(start), result)
}
