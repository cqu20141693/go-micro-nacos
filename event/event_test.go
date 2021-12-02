package event

import (
	"fmt"
	"testing"
	"time"
)

func TestEvent(t *testing.T) {

	RegisterHook(Start, func() {
		fmt.Println("register start")
	})
	RegisterHook(LocalConfigComplete, func() {
		fmt.Println("register LocalConfigComplete 1")
	})
	go func() {
		RegisterHook(ConfigComplete, func() {
			fmt.Println("register ConfigComplete")
		})
	}()
	go func() {
		TriggerEvent(LocalConfigComplete)
	}()
	time.Sleep(1 * time.Second)
	RegisterHook(LocalConfigComplete, func() {
		fmt.Println("register LocalConfigComplete 2")
	})

	time.Sleep(1 * time.Second)

}
