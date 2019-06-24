/***********************************************************************
# File Name: photoweb.go
# Author: TiDao
# mail: songzhaofeng411@gmail.com
# Created Time: 2019-06-24 17:01:23
*********************************************************************/
package main
import(
    "io"
    "log"
    "net/http"
    "os"
    "io/ioutil"
)

const (
    UPLOAD_DIR = "./uploads"
)

func uploadHandler(w http.ResponseWriter,r *http.Request){
    if r.Method == "GET"{
        io.WriteString(w,"<form method=\"POST\" action=\"/upload\" "+" enctype=\"multipart/form-data\">"+"Choose an image to upload: <input name=\"image\" type=\"file\" />"+"<input type=\"submit\" value=\"Upload\" />"+"</form>")
        return
    }
    if r.Method == "POST"{
        f,h,err := r.FormFile("image")
        if err != nil{
            http.Error(w,err.Error(),http.StatusInternalServerError)
            return
        }
        filename := h.Filename
        defer f.Close()
        t,err := os.Create(UPLOAD_DIR + "/"+filename)
        if err != nil{
            http.Error(w,err.Error(),http.StatusInternalServerError)
            return
        }
        defer t.Close()
        if _,err := io.Copy(t,f); err != nil{
            http.Error(w,err.Error(),http.StatusInternalServerError)
            return
        }
        http.Redirect(w,r,"/view?id="+filename,http.StatusFound)
    }
}

func viewHandler(w http.ResponseWriter,r *http.Request){
    imageId := r.FormValue("id")
    imagePath := UPLOAD_DIR+"/"+imageId
    if exists := isExists(imagePath);!exists{
        error_response := "Not found the image id="+imageId
        io.WriteString(w,error_response)
        return
    }
    w.Header().Set("contect-Type","image")
    http.ServeFile(w,r,imagePath)
}

func isExists(path string) bool{
    _,err := os.Stat(path)
    if err == nil{
        return true
    }
    return os.IsExist(err)
}

func listHandler(w http.ResponseWriter,r *http.Request){
    fileInfoArr,err := ioutil.ReadDir("./uploads")
    if err != nil{
        http.Error(w,err.Error(),http.StatusInternalServerError)
        return
    }
    var listhtml string
    for _,fileInfo := range fileInfoArr{
        imgId := fileInfo.Name()
        strings := "<li><a href=\"/view?id="+string(imgId)+"\">imgid</a></li>"
        listhtml += strings
    }
    io.WriteString(w,"<ol>"+listhtml+"</ol>")
}

func main(){
    http.HandleFunc("/",listHandler)
    http.HandleFunc("/view",viewHandler)
    http.HandleFunc("/upload",uploadHandler)
    err := http.ListenAndServe(":8000",nil)
    if err != nil{
        log.Fatal("ListenAndServe: ",err.Error())
    }
}
