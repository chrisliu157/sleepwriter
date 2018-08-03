package store

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Store struct {
	Pool *redis.Pool
}

func NewStore() (*Store, error) {
	store := new(Store)
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	store.Pool = pool
	return store, nil
}

// Get value from key
func (store Store) Get(key string) ([]byte, error) {
	conn := store.Pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("Store - Error - Get - %v : %v", key, err)
	}
	return data, err
}

// Set key value pair
func (store Store) Set(key string, value []byte) error {
	conn := store.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return fmt.Errorf("Store - Error - Set - Key: %v  Value: %v - %v", key, string(value), err)
	}
	return nil
}

// Delete key
func (store Store) Delete(key string) error {
	conn := store.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	if err != nil {
		return fmt.Errorf("Store - Error - Delete - Key: %v  - %v", key, err)
	}
	return nil
}

// Creates a transaction which will add key to priority queue with score
// then adds attributes to key value
// TODO transition key value store to hmset, unnecssary serialization right now
func (store Store) PQueueAdd(queue string, score int, key string, value []byte) error {
	conn := store.Pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	zaddErr := conn.Send("ZADD", queue, score, []byte(key))
	if zaddErr != nil {
		return fmt.Errorf("Store - Error - PQueueAdd - Queue: %v Score: %v : %v", queue, score, zaddErr)
	}
	setErr := conn.Send("SET", key, value)
	if setErr != nil {
		return fmt.Errorf("Store - Error - PQueueAdd - Queue: %v Score: %v : %v", queue, score, setErr)
	}

	_, execErr := conn.Do("EXEC")
	return execErr
}

// Transaction pops lowest scored id from priority queue
// then returns associated attributes
func (store Store) PQueuePop(queue string) ([]byte, error) {
	key := ""

	conn := store.Pool.Get()
	defer conn.Close()

	// loop through transaction until successful
	for {
		members, err := redis.Strings(conn.Do("ZRANGE", queue, 0, 0))
		if err != nil {
			return []byte(""), err
		}

		if len(members) != 1 {
			return []byte(""), redis.ErrNil
		}

		conn.Send("MULTI")
		conn.Send("ZREM", queue, members[0])
		queued, execErr := conn.Do("EXEC")
		if execErr != nil {
			return []byte(""), err
		}

		if queued != nil {
			key = members[0]
			break
		}
	}

	// Get attributes from key value
	result, getErr := store.Get(key)
	if getErr != nil {
		return result, getErr
	}

	return result, nil
}

// func main() {
// 	pool, err := NewStore()
// 	fmt.Println(err)
// 	popErr := pool.PQueueAdd("testqueue", 1, "myjobid", []byte("mystringofstuff"))
//
// 	fmt.Println(popErr)
//
// 	value, popErr := pool.PQueuePop("testqueue")
// 	fmt.Println(popErr)
// 	fmt.Println(string(value))
//
// }
