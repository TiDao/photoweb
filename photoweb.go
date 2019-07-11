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
    "html/template"
    "path"
)

const (
    UPLOAD_DIR = "./uploads"
    TEMPLATE_DIR = "./view"
)


//var templates map[string]*template.Template = make(map[string]*template.Template)

var templates = make(map[string]*template.Template)
func init(){
    fileInfoArr,err := ioutil.ReadDir(TEMPLATE_DIR)
    if err != nil{
        panic(err)
        return
    }
    var templateName,templatePath string
    for _,fileInfo := range fileInfoArr{
        templateName = fileInfo.Name()
        if ext := path.Ext(templateName); ext != ".html"{
            continue
        }
        templatePath = TEMPLATE_DIR+"/"+templateName
        t := template.Must(template.ParseFiles(templatePath)) 
        templates[templateName] = t
        log.Println("loading template:",templatePath)
    }
}

func renderHtml(w http.ResponseWriter,templateName string,locals map[string]interface{})  error{
    err := templates[templateName].Execute(w,locals)
    return err
}

func check(err error){
    if err != nil{
        panic(err)
    }
}

func uploadHandler(w http.ResponseWriter,r *http.Request){
    if r.Method == "GET"{
        err := renderHtml(w,"upload.html",nil)
        check(err)
    }
    if r.Method == "POST"{
        f,h,err := r.FormFile("image")
        check(err)
        filename := h.Filename
        defer f.Close()
        t,err := os.Create(UPLOAD_DIR + "/"+filename)
        check(err)
        defer t.Close()
        _,err = io.Copy(t,f)
        check(err)
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
    check(err)
    locals := make(map[string]interface{})
    images := []string{}
    for _,fileInfo := range fileInfoArr{
        images = append(images,fileInfo.Name())
    }
    locals["images"] = images
    err  = renderHtml(w,"list.html",locals)
    check(err)
}


func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter,r *http.Request){
        defer func(){
            if err,ok:= recover().(error); ok{
                http.Error(w,err.Error(),http.StatusInternalServerError)
                log.Println("warning:panic in %v - %v",fn,err)
            }
        }()
        fn(w,r)
    }
}

func main(){
    http.HandleFunc("/",safeHandler(listHandler))
    http.HandleFunc("/view",safeHandler(viewHandler))
    http.HandleFunc("/upload",safeHandler(uploadHandler))
    err := http.ListenAndServe(":8000",nil)
    if err != nil{
        log.Fatal("ListenAndServe: ",err.Error())
    }
}
