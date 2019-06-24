/***********************************************************************
# File Name: hello.go
# Author: TiDao
# mail: songzhaofeng411@gmail.com
# Created Time: 2019-06-24 16:35:29
*********************************************************************/
package main
import(
    "io"
    "log"
    "net/http"
)

func helloHandler(w http.ResponseWriter,r *http.Request){
    io.WriteString(w,"hello,world")
    io.WriteString(w,"test")
}

func main(){
    go http.HandleFunc("/hello",helloHandler)
    err := http.ListenAndServe(":8000",nil)
    if err != nil{
        log.Fatal("listenAndServer:",err.Error())
    }
}
