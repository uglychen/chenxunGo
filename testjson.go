package main

import (
    "encoding/json"
    "fmt"
)

type registerData struct {
    OpenId   string
    NickName string
    HeadUrl  string
    Gender   string
    Age      string
}

func main() {

    var m registerData
    err := json.Unmarshal(b, &m)

}
