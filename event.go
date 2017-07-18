package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "io/ioutil"
    "log"
    "net/http"
    //"reflect"
    "strconv"
    "time"
)

type userData struct {
    UserId   int
    OpenId   string
    NickName string
    HeadUrl  string
    Gender   byte
    gold     int
    cash     int
    Age      int
}

func readAward(w http.ResponseWriter, req *http.Request) {

    //阅读奖励 taskId = 5
    retValue := NewBaseJsonData()
    req.ParseForm()

    param_userId, _ := req.Form["userId"]
    param_type, _ := req.Form["type"]
    userId := param_userId[0]
    tab_type := param_type[0]

    db, err := sql.Open("mysql", "reportadmin:123456@tcp(202.120.1.109:3306)/ppserver?charset=utf8")
    if err != nil {
        log.Println(err)
    }

    tx, _ := db.Begin()

    eventMap := getEventMap(userId, tab_type, db)
    taskStatus := eventMap["taskStatus"].(int)
    today := time.Now().Format("2006-01-02 00:00:00")

    if len(eventMap) > 0 && taskStatus != 1 && eventMap["endTime"] != today {
        tx.Exec("update eventTable set taskStatus = ?, endTime = ? where userId=? and taskId=?",
            1, today, userId, tab_type)
        tx.Exec("update userInfo set gold = gold + ? where userId=?", 10, userId)

        retValue.Code = 200
        retValue.Data = 10
        retValue.Message = "success"
        bytes, _ := json.Marshal(retValue)
        fmt.Fprint(w, string(bytes), "\n")

    } else {
        retValue.Code = 300
        retValue.Message = "failed"
        bytes, _ := json.Marshal(retValue)
        fmt.Fprint(w, string(bytes), "\n")
    }

    tx.Commit()

}

func updateUserData(w http.ResponseWriter, req *http.Request) {

    //主要是检查用户年龄age是否更新
    var data userData
    result, _ := ioutil.ReadAll(req.Body)
    req.Body.Close()
    log.Println("updateUserData:", string(result))
    json.Unmarshal([]byte(result), &data)

    userId := strconv.Itoa(data.UserId)
    age := strconv.Itoa(data.Age)

    db, err := sql.Open("mysql", "reportadmin:123456@tcp(202.120.1.109:3306)/ppserver?charset=utf8")
    if err != nil {
        log.Println(err)
    }

    userMap := getUserInfo(userId, db)
    age1 := strconv.Itoa(int(userMap["age"].(uint8)))
    retValue := NewBaseJsonData()

    if age1 != age {
        if uint8(0) == userMap["age"].(uint8) {
            log.Println("update user's age")
            updateAge(userId, db, age)

            retValue.Code = 200
            retValue.Message = "success"
            bytes, _ := json.Marshal(retValue)
            fmt.Fprint(w, string(bytes), "\n")

        } else {
            tx, _ := db.Begin()
            sql := "update userInfo set age = ? where userId = ? "
            tx.Exec(sql, age1)
            tx.Commit()

            retValue.Code = 200
            retValue.Message = "success"
            bytes, _ := json.Marshal(retValue)
            fmt.Fprint(w, string(bytes), "\n")
        }
    } else {
        retValue.Code = 300
        retValue.Message = "failed"
        bytes, _ := json.Marshal(retValue)
        fmt.Fprint(w, string(bytes), "\n")
    }

}

func shareApp(w http.ResponseWriter, req *http.Request) {
    //分享app至微信
    req.ParseForm()
    param_userId, _ := req.Form["userId"]
    userId := param_userId[0]

    db, err := sql.Open("mysql", "reportadmin:123456@tcp(202.120.1.109:3306)/ppserver?charset=utf8")
    if err != nil {
        log.Println(err)
    }

    var taskId string = "4"
    userMap := getUserInfo(userId, db)
    eventMap := getEventMap(userId, taskId, db)
    log.Println(userMap)
    log.Println(eventMap)

    today := time.Now().Format("2006-01-02 00:00:00")
    shareDate := eventMap["endTime"].(string)
    taskStatus := eventMap["taskStatus"].(int)
    retValue := NewBaseJsonData()
    if eventMap["createTime"] == eventMap["endTime"] {
        retValue.Code = 200
        retValue.Data = 0
        retValue.Message = "shared"
        bytes, _ := json.Marshal(retValue)
        fmt.Fprint(w, string(bytes), "\n")
    }

    if today != shareDate && taskStatus != 1 {
        tx, _ := db.Begin()
        tx.Exec("update eventTable set taskStatus = ?, endTime = ? where userId=? and taskId=?",
            1, time.Now(), userId, 4)
        tx.Exec("update userInfo set gold = gold + ? where userId=?", 40, userId)
        tx.Commit()

        retValue.Code = 200
        retValue.Data = 40
        retValue.Message = "success"
        bytes, _ := json.Marshal(retValue)
        fmt.Fprint(w, string(bytes), "\n")
    } else {
        retValue.Code = 200
        retValue.Data = 0
        retValue.Message = "shared"
        bytes, _ := json.Marshal(retValue)
        fmt.Fprint(w, string(bytes), "\n")
    }
}

func loginEvent(userId int, db *sql.DB) {
    //用户微信首次登录事件为一次性事件，给用户加cash = 2，taskId = 1
    sql := `INSERT INTO eventTable(userId, taskId, taskStatus, addGoldCoin, addCash, eventId) 
        values (? , ?, ?, ?, ?, ?)`
    tx, _ := db.Begin()
    tx.Exec(sql, userId, 1, 1, 0, 2, 0)
    tx.Exec("update userInfo set cash = cash + 2 where userId =?", userId)
    tx.Commit()
}

func updateAge(userId string, db *sql.DB, age string) {

    //用户更新资料  taskId = 2
    sql := `INSERT INTO eventTable(userId, taskId, taskStatus, addGoldCoin, addCash, eventId) 
            values (?, ?, ?, ?, ?, ?)`
    tx, _ := db.Begin()
    tx.Exec(sql, userId, 2, 1, 0, 0, 0)
    tx.Exec("update userInfo set age = ?, gold = gold + ? where userId = ? ", age, 400, userId)
    tx.Commit()

}

func createSignEvent(userId int, db *sql.DB) {

    // 创建签到任务 taskId = 3
    sql := `INSERT INTO eventTable(userId, taskId, taskStatus, addGoldCoin, addCash, eventId) 
            values (?, ?, ?, ?, ?, ?)`

    tx, _ := db.Begin()
    tx.Exec(sql, userId, 3, 0, 0, 0, 0)
    tx.Commit()
}

func createShareAppEvent(userId int, db *sql.DB) {

    //创建shareApp分享任务 taskId = 4
    sql := `INSERT INTO eventTable(userId, taskId, taskStatus, addGoldCoin, addCash, eventId) 
            values (?, ?, ?, ?, ?, ?)`

    tx, _ := db.Begin()
    tx.Exec(sql, userId, 4, 0, 0, 0, 0)
    tx.Commit()

}

func createReadAwardEvent(userId int, db *sql.DB) {

    //创建shareApp分享任务 taskId = 5
    today := time.Now().Format("2006-01-02 00:00:00")
    sql := `INSERT INTO eventTable(userId, taskId, createTime,endTime,
        taskStatus, addGoldCoin, addCash, eventId) values (?, ?, ?, ?, ?, ?, ?, ?)`

    tx, _ := db.Begin()
    tx.Exec(sql, userId, 5, today, today, 0, 0, 0, 0)
    tx.Exec(sql, userId, 6, today, today, 0, 0, 0, 0)
    tx.Exec(sql, userId, 7, today, today, 0, 0, 0, 0)
    tx.Exec(sql, userId, 8, today, today, 0, 0, 0, 0)
    tx.Commit()
}

func getEventMap(userId string, taskId string, db *sql.DB) map[string]interface{} {

    log.Println("call getEventMap")

    evevtMap := make(map[string]interface{})
    var taskStatus int
    var createTime string
    var endTime string

    rows, err := db.Query(`SELECT taskStatus, createTime, endTime FROM eventTable where userId = ? and taskId =?`,
        userId, taskId)
    if err != nil {
        log.Println(err)
    }
    defer rows.Close()
    for rows.Next() {
        err := rows.Scan(&taskStatus, &createTime, &endTime)
        if err != nil {
            log.Fatal(err)
        }
        log.Println("get taskStatus:", taskStatus)
    }

    evevtMap["taskStatus"] = taskStatus
    evevtMap["createTime"] = createTime
    evevtMap["endTime"] = endTime

    log.Println("getEventMap:", evevtMap)

    return evevtMap
}

func getUserInfo(userId string, db *sql.DB) map[string]interface{} {

    log.Println("call getUserInfo")

    retMap := make(map[string]interface{})

    var wechatId string
    var nickName string
    var headUrl string
    var gender byte
    var age byte
    var gold int
    var cash int

    sql := "select wechatId, nickName, headUrl, gender, age, gold, cash FROM userInfo where userId = ?"
    rows, err := db.Query(sql, userId)

    if err != nil {
        log.Println(err)
    }
    defer rows.Close()
    if rows.Next() {
        err := rows.Scan(&wechatId, &nickName, &headUrl, &gender, &age, &gold, &cash)
        if err != nil {
            log.Fatal(err)
        }

    } else {

        return retMap
    }

    retMap["userId"] = userId
    retMap["oepnId"] = wechatId
    retMap["nickName"] = nickName
    retMap["headUrl"] = headUrl
    retMap["gender"] = gender
    retMap["age"] = age
    retMap["gold"] = gold
    retMap["cash"] = cash

    log.Println("getUserInfo:retMap", retMap)

    return retMap
}
