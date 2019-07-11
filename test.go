/***********************************************************************
# File Name: test.go
# Author: TiDao
# mail: songzhaofeng411@gmail.com
# Created Time: 2019-07-05 11:14:19
*********************************************************************/
package main
import(
    "fmt"
    "io/ioutil"
//    "path"
)

func main(){
    fileInfoArr,err := ioutil.ReadDir("./view")
    if err!=nil{
        panic(err)
        return
    }
    for _,fileInfo := range fileInfoArr{
//        if ext := path.Ext(fileInfo.Name());ext!=".html"{
//            continue
//        }
        fmt.Println(fileInfo.Name())
    }
}
