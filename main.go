package main

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"
    "github.com/KenmyZhang/image/app"
)

var (
    imagPath = "http://imgsrc.baidu.com/imgad/pic/item/267f9e2f07082838b5168c32b299a9014c08f1f9.jpg"
    destPath = "./test.jpg"
    width = 300
    height = 300
    option = 1
)


func main() {
    resp, _ := http.Get(imagPath)

    var body []byte
    body, _ = ioutil.ReadAll(resp.Body)
   
    var err error
    var data *bytes.Buffer
    if data, err = app.SetScaleImage(body, width, height, option); err != nil {
        fmt.Println(err.Error())
    }

    if data != nil {
        if err := app.SaveImage(data.Bytes(), destPath); err != nil {
            fmt.Println(err.Error())
        }
    } else {
        fmt.Println("data is nil")
    }

    return
}