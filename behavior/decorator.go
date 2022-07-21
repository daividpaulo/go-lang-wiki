package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type MyEntity struct {
	Id   string
	Name string
}

type serviceFunc func(MyEntity) MyEntity

func loggerDecorator(fun serviceFunc, logger *log.Logger) serviceFunc {
	return func(entity MyEntity) MyEntity {
		fn := func(entity MyEntity) (result MyEntity) {
			defer func(t time.Time) {
				logger.Printf("took=%v, n=%v, result=%v", time.Since(t), entity.Id, result)
			}(time.Now())

			return fun(entity)
		}
		return fn(entity)
	}
}

func cacheDecorator(fun serviceFunc, cache *sync.Map) serviceFunc {
	return func(entity MyEntity) MyEntity {
		fn := func(entity MyEntity) MyEntity {
			key := fmt.Sprintf("n=%d", entity.Id)
			val, ok := cache.Load(key)
			if ok {
				return val.(MyEntity)
			}
			result := fun(entity)
			cache.Store(key, result)
			return result
		}
		return fn(entity)
	}
}

func saveEntity(n MyEntity) MyEntity {
	ch := make(chan MyEntity)

	go func(ch chan MyEntity, k MyEntity) {
		ch <- k
	}(ch, n)

	result := <-ch

	return result
}

func updateEntity(n MyEntity) MyEntity {
	ch := make(chan MyEntity)

	go func(ch chan MyEntity, k MyEntity) {
		k.Name = "New name configured"
		ch <- k
	}(ch, n)

	result := <-ch

	return result
}

func main() {
	f := cacheDecorator(saveEntity, &sync.Map{})
	g := loggerDecorator(f, log.New(os.Stdout, "Save entity ", 1))

	g(MyEntity{Id: "123", Name: "Teste Minha Entidade"})
	g(MyEntity{Id: "123", Name: "Teste Minha Entidade"})
	g(MyEntity{Id: "123", Name: "Teste Minha Entidade"})

	f = cacheDecorator(updateEntity, &sync.Map{})
	g = loggerDecorator(f, log.New(os.Stdout, "Update Entity ", 1))

	g(MyEntity{Id: "123", Name: "Teste Minha Entidade"})
	g(MyEntity{Id: "123", Name: "Teste Minha Entidade"})
}
