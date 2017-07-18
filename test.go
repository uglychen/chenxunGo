package main

import (
    "fmt"
)

/* 这段是抄的 */
func ByteToBinaryString(data byte) (str string) {
    var a byte
    for i := 0; i < 8; i++ {
        a = data
        data <<= 1
        data >>= 1
        switch a {
        case data:
            str += "0"
        default:
            str += "1"
        }
        data <<= 1
    }
    return str
}

/* 自己写了个 */
func ByteToBinaryString2(data byte) (str string) {
    var a byte = 0x80
    for i := 0; i < 8; i++ {
        switch a & data {
        case 0:
            str += "0"
        default:
            str += "1"
        }
        a >>= 1
    }
    return str
}
func main() {
    var t byte = 15
    fmt.Printf("%d = [%s] [%s]\n", t, ByteToBinaryString(t), ByteToBinaryString2(t))
}
