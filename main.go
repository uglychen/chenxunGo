package main

import (
    "flag"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
)

var dataChain = make(chan map[string]interface{}, 5)
var HandlerMap = make(map[string]HandlersFunc)

type BaseJsonData struct {
    Code    int         `json:"code"`
    Data    interface{} `json:"data"`
    Message string      `json:"message"`
}

func NewBaseJsonData() *BaseJsonData {
    return &BaseJsonData{}
}

type HandlersFunc func(w http.ResponseWriter, req *http.Request)

func InitHandlerMap(m map[string]HandlersFunc) {

    HandlerMap["/loginAndOut"] = LoginAndOut
    HandlerMap["/signIn"] = userSignIn
    HandlerMap["/updateUserData"] = updateUserData
    HandlerMap["/shareApp"] = shareApp
    HandlerMap["/readAward"] = readAward

    log.Println("Listen HandlerMap's length:", len(HandlerMap))

}

func SetListenHandle() {

    for key, value := range HandlerMap {
        http.HandleFunc(key, value)
    }

}

func main() {

    logFile, err := os.OpenFile("./ppserver.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
    if err != nil {
        fmt.Println("open file error=", err.Error())
        os.Exit(-1)
    }
    defer logFile.Close()

    writers := []io.Writer{
        logFile,
        os.Stdout,
    }

    fileAndStdoutWriter := io.MultiWriter(writers...)

    log.SetOutput(fileAndStdoutWriter)
    log.SetFlags(log.Lshortfile | log.LstdFlags)

    log.Println("===start the ppserver===.")

    //初始化redis连接池
    flag.Parse()
    pool = newPool(*redisServer, *redisPassword)

    conn := pool.Get()
    defer conn.Close()

    //开启备份模式redis
    //go saveRedis(dataChain)

    InitHandlerMap(HandlerMap)
    SetListenHandle()

    if err := http.ListenAndServe(":9091", nil); err != nil {
        panic(err)
    }
}
