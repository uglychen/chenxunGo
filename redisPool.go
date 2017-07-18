package main

import (
    "flag"
    //"fmt"
    "github.com/garyburd/redigo/redis"
    "time"
)

//声明一些全局变量
var (
    pool          *redis.Pool
    redisServer   = flag.String("redisServer", "127.0.0.1:6379", "")
    redisPassword = flag.String("redisPassword", "123456", "")
)

//初始化一个pool
func newPool(server, password string) *redis.Pool {
    return &redis.Pool{
        MaxIdle:     3,
        MaxActive:   5,
        IdleTimeout: 240 * time.Second,
        Dial: func() (redis.Conn, error) {
            c, err := redis.Dial("tcp", server)
            if err != nil {
                return nil, err
            }
            if _, err := c.Do("AUTH", password); err != nil {
                c.Close()
                return nil, err
            }
            return c, err
        },
        TestOnBorrow: func(c redis.Conn, t time.Time) error {
            if time.Since(t) < time.Minute {
                return nil
            }
            _, err := c.Do("PING")
            return err
        },
    }
}

/*
func main() {
    flag.Parse()
    pool = newPool(*redisServer, *redisPassword)

    conn := pool.Get()
    defer conn.Close()
    //redis操作
    v, err := redis.String(conn.Do("HGET", "userInfo:chen_123456", "userId"))
    if err != nil {
        fmt.Println(err)
        //return
    }
    fmt.Println(v)
}
*/
