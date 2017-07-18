package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "log"
    "net/http"
    "time"
)

//import "github.com/garyburd/redigo/redis"

type signJson struct {
    UserId    string
    SignCount int
}

func userSignIn(w http.ResponseWriter, req *http.Request) {

    retValue := NewBaseJsonData()

    req.ParseForm()
    param_id, _ := req.Form["userId"]
    userId := param_id[0]
    taskId := "3"

    db, err := sql.Open("mysql", "reportadmin:123456@tcp(202.120.1.109:3306)/ppserver?charset=utf8")
    if err != nil {
        log.Println(err)
    }

    eventMap := getEventMap(userId, taskId, db)
    userMap := getUserInfo(userId, db)
    fmt.Println(eventMap)
    fmt.Println(userMap)

    old_gold, _ := userMap["gold"].(int)
    sign_date := eventMap["endTime"].(string)
    sign_count := eventMap["taskStatus"].(int)
    log.Println("old_gold:", old_gold, "sign_date:", sign_date, "sign_count:", sign_count)

    s := []int{10, 15, 20, 25, 30, 35, 40}

    t, _ := time.Parse("2006-01-02 15:04:05", sign_date)
    tmp := t.Format("2006-01-02 00:00:00")
    log.Println("[last time sign date:]", tmp)

    date := time.Now()
    d, _ := time.ParseDuration("-24h")
    t2 := date.Add(d).Format("2006-01-02 00:00:00")
    log.Println("[last date:]", t2)

    var addGoldCoin int
    tx, _ := db.Begin()
    if date.Format("2006-01-02 00:00:00") != tmp {

        if tmp == t2 {

            if sign_count < 7 {
                addGoldCoin = s[sign_count]
            } else {
                addGoldCoin = 40
            }

            tx.Exec("update eventTable set `endTime`=?,`taskStatus`=? where `userId`=? and `taskId`=?",
                date, sign_count+1, userId, taskId)
            tx.Exec("update userInfo set `gold` =? where `userId`=?", old_gold+addGoldCoin, userId)
            log.Println("addGoldCoin:", addGoldCoin)

            var retData signJson
            retData.SignCount = sign_count + 1
            retData.UserId = userId

            retValue.Code = 200
            retValue.Data = retData
            retValue.Message = "success"
            bytes, _ := json.Marshal(retValue)
            fmt.Fprintln(w, string(bytes))

        } else {

            addGoldCoin = 10
            tx.Exec("update eventTable set `endTime`=? ,`taskStatus`=? where `userId`=? and `taskId`=?",
                date, 1, userId, taskId)
            tx.Exec("update userInfo set `gold` = ? where `userId`=?", old_gold+addGoldCoin, userId)
            log.Println("addGoldCoin:", addGoldCoin)

            var retData signJson
            retData.SignCount = 1
            retData.UserId = userId

            retValue.Code = 200
            retValue.Data = retData
            retValue.Message = "success"
            bytes, _ := json.Marshal(retValue)
            fmt.Fprintln(w, string(bytes))
        }
    } else {

        retValue.Code = 200
        retValue.Message = "failed"
        bytes, _ := json.Marshal(retValue)
        fmt.Fprintln(w, string(bytes))
    }

    tx.Commit()
}
