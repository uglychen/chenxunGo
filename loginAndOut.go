package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"
)

type registerData struct {
    OpenId   string
    NickName string
    HeadUrl  string
    Gender   byte
    Age      int
}

type retJson struct {
    OpenId     string
    NickName   string
    HeadUrl    string
    Gender     byte
    Age        int
    TaskStatus int
}

func LoginAndOut(w http.ResponseWriter, req *http.Request) {

    req.ParseForm()
    param_type, _ := req.Form["type"]
    reqType := param_type[0]

    //用户注册
    if reqType == "register" {
        UserRegister(w, req)
    }

    //用户登录接口
    if reqType == "login" {
        UserLogin(w, req)
    }
}

func UserLogin(w http.ResponseWriter, req *http.Request) {

    retValue := NewBaseJsonData()

    if req.Method == "GET" {

        param_openId, _ := req.Form["userId"]
        userId := param_openId[0]

        db, err := sql.Open("mysql", "reportadmin:123456@tcp(202.120.1.109:3306)/ppserver?charset=utf8")
        if err != nil {
            log.Fatalf("Open database error: %s\n", err)
        }
        defer db.Close()

        retMap := getUserInfo(userId, db)
        //log.Println(retMap)
        //userId := retMap["userId"].(int)
        if len(retMap) > 0 {

            evevtMap := getEventMap(userId, "3", db)
            fmt.Println("evevtMap:", evevtMap)

            retMap["signCount"] = evevtMap["taskStatus"]

            fmt.Println(evevtMap)
            fmt.Print(retMap)

            retValue.Code = 200
            retValue.Data = retMap
            retValue.Message = "success"
            bytes, _ := json.Marshal(retValue)
            fmt.Fprint(w, string(bytes), "\n")

        } else {

            retValue.Code = 200
            retValue.Message = "failed"
            bytes, _ := json.Marshal(retValue)
            fmt.Fprint(w, string(bytes), "\n")

        }

    } else {

        retValue.Code = 200
        retValue.Message = "failed"
        bytes, _ := json.Marshal(retValue)
        fmt.Fprint(w, string(bytes), "\n")
    }
}

func UserRegister(w http.ResponseWriter, req *http.Request) {

    if req.Method == "POST" {
        result, _ := ioutil.ReadAll(req.Body)
        req.Body.Close()
        log.Println("recv body:", string(result))

        var m registerData
        var userId int
        var nickName string
        json.Unmarshal([]byte(result), &m)

        db, err := sql.Open("mysql", "reportadmin:123456@tcp(202.120.1.109:3306)/ppserver?charset=utf8")
        if err != nil {
            log.Fatalf("Open database error: %s\n", err)
        }
        defer db.Close()

        rows1, err := db.Query("SELECT userId, nickName FROM userInfo where wechatId = ?", m.OpenId)
        if err != nil {
            log.Println(err)
        }
        defer rows1.Close()

        flag := rows1.Next()
        log.Println("用户请求的openId是否已经注册过了：", flag)

        if flag {

            //查找userId
            rows, err := db.Query("SELECT userId, nickName FROM userInfo where wechatId = ?", m.OpenId)
            if err != nil {
                log.Println(err)
            }
            defer rows.Close()
            for rows.Next() {
                err := rows.Scan(&userId, &nickName)
                if err != nil {
                    log.Fatal(err)
                }
                log.Println("userId:", userId, "nickName:", nickName)
            }

            retValue := NewBaseJsonData()
            retValue.Code = 300
            retValue.Message = "failed"
            bytes, _ := json.Marshal(retValue)
            fmt.Fprintln(w, string(bytes)) //返回json

        } else {

            dataMap := make(map[string]interface{})
            dataMap["OpenId"] = m.OpenId
            dataMap["NickName"] = m.NickName
            dataMap["HeadUrl"] = m.HeadUrl
            dataMap["Gender"] = m.Gender
            dataMap["Age"] = m.Age
            dataMap["Gold"] = 0
            dataMap["Cash"] = 0
            dataMap["SignCount"] = 0

            sql := `INSERT INTO userInfo(wechatId, nickName, headUrl, gender, age, gold, cash)
                   VALUES (?, ?, ?, ?, ?, ?, ?)`

            tx, _ := db.Begin()
            tx.Exec(sql, dataMap["OpenId"].(string), dataMap["NickName"], dataMap["HeadUrl"],
                dataMap["Gender"], dataMap["Age"], 0, 0)
            tx.Commit()

            //查找userId
            rows2, err := db.Query("SELECT userId, nickName FROM userInfo where wechatId = ?", m.OpenId)
            if err != nil {
                log.Println(err)
            }
            defer rows2.Close()
            for rows2.Next() {
                err := rows2.Scan(&userId, &nickName)
                if err != nil {
                    log.Fatal(err)
                }
                log.Println("userId:", userId, "nickName:", nickName)
            }

            dataMap["UserId"] = userId
            retValue := NewBaseJsonData()
            retValue.Code = 200
            retValue.Message = "success"
            retValue.Data = dataMap
            bytes, _ := json.Marshal(retValue)
            fmt.Fprintln(w, string(bytes)) //返回json

            //创建任务
            loginEvent(userId, db)
            createSignEvent(userId, db)
            createShareAppEvent(userId, db)
            createReadAwardEvent(userId, db)

            //dataChain <- dataMap //缓存到redis中去
        }

    } else {

        retValue := NewBaseJsonData()
        retValue.Code = 400
        retValue.Message = "注册失败"
        bytes, _ := json.Marshal(retValue)
        fmt.Fprintln(w, string(bytes))

    }
}

func saveRedis(c chan map[string]interface{}) {

    for {
        data := <-c
        fmt.Println("recv dataMap from chan:", data)

        conn := pool.Get()
        defer conn.Close()

        key := "userInfo:" + strconv.Itoa((data["UserId"]).(int))
        fmt.Println("key:", key)

        //查看redis缓存中是否已经存在openid相关的账号信息
        v, _ := conn.Do("HSET", key, "wechatId", data["OpenId"])
        //fmt.Println("v的值是：", reflect.TypeOf(v))

        if v != int64(0) {

            conn.Do("HSET", key, "userId", data["UserId"])
            conn.Do("HSET", key, "wechatId", data["OpenId"])
            conn.Do("HSET", key, "nickName", data["NickName"])
            conn.Do("HSET", key, "headlUrl", data["HeadUrl"])
            conn.Do("HSET", key, "gender", data["Gender"])
            conn.Do("HSET", key, "age", data["Age"])
            conn.Do("HSET", key, "gold", data["Gold"])
            conn.Do("HSET", key, "cash", data["Cash"])

        } else {
            fmt.Println("redis缓存中已经存在OpenId:", (data["OpenId"]))
        }

    }
}
