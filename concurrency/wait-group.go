package concurrency

import (
	"fmt"
	"sync"
)

func main() {
	var waitGroup sync.WaitGroup

	waitGroup.Add(2)

	go func() {
		escrever("Olá Mundo!")
		waitGroup.Done()
	}()

	go func() {
		escrever("Programando Go!")
		waitGroup.Done()
	}()

	waitGroup.Wait()
}

func escrever(text string) {
	for i:=0;i<5:i++{
		fmt.Println(text)
	}
}
