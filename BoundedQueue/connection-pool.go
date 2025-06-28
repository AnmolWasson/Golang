package main

import "fmt"

type Connection struct {
	id    string
	query func(q string) string
}

type Pool struct {
	connPoolMaxSize int
	pool            *BoundedBlockingQueue[Connection]
}

func ConnectionPool() *Pool {

	size := 10
	pool := &Pool{
		connPoolMaxSize: size,
		pool:            Init[Connection](size),
	}

	for i := 0; i < size; i++ {
		conn := Connection{
			id: fmt.Sprintf("conn_id_%d", i),
		}

		conn.query = func(q string) string {
			return fmt.Sprintf("Going to use conn %d for query %s", conn.id, q)
		}
		pool.pool.Enqueue(conn)

	}
	return pool
}

func (conn *Pool) Acquire() Connection {
	return conn.pool.Dequeue()
}

func (conn *Pool) Release(c Connection) {
	conn.pool.Enqueue(c)
}
