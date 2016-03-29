package main

import(
    "net/http"
    "fmt"
    "os"
    "io/ioutil"
    "encoding/json"
    "net/url"
    "strings"	
)

var configData map[string]interface{}
var apiData []interface{}
var contentTypeConfig map[string]interface{}

func handleHttp(w http.ResponseWriter, r *http.Request){
    u, _ := url.Parse(r.RequestURI)
   handleUrl := u.Path
   // 处理是否匹配API接口
  for _,t := range apiData{
	if t.([]interface{})[0].(string) == u.Path{
		handleUrl = t.([]interface{})[1].(string)
		break	
	}
	
  }	
   
  filePath := configData["workspace"].(string) + handleUrl
  fileContent := GetContentFromFile(filePath)
  fmt.Println("请求路径："+u.Path)
   fmt.Println("读取文件："+filePath)	  
	
	if fileContent ==nil{
		w.WriteHeader(404)
		fmt.Println("状态：404")
		return
	}
	
	// 处理Content-Type
	strs := strings.Split(handleUrl, ".")
	str := strs[len(strs) - 1]
	w.Header().Set("Content-Type", "application/octet-stream;charset=utf8")
	 for k,v := range contentTypeConfig{
		if k == "."+str{
			w.Header().Set("Content-Type", v.(string)+";charset=utf8")
			break
		}
	}
	
	w.Write(fileContent)
	fmt.Println("状态：200")
}

func StartServer(){
	
    fmt.Println("服务器启动！")
    fmt.Print("监听域名：" + configData["host"].(string))
    fmt.Print("监听端口：" + configData["port"].(string))	
    http.HandleFunc("/", handleHttp)
    http.ListenAndServe(configData["host"].(string)+":"+configData["port"].(string), nil)	
}


func GetContentFromFile(filePath string) []byte{
	fp, err := os.Open(filePath)
        defer fp.Close()
        if err != nil{
			return nil
		}
        text, _ := ioutil.ReadAll(fp)
        
        return text
}

func HandleConfig(){
        configData = make(map[string]interface{})
        jsonText := GetContentFromFile("config.txt")
        json.Unmarshal([]byte(jsonText), &configData)
		
       // 处理api
        jsonText = GetContentFromFile(configData["api_change_to_file_config"].(string))
        json.Unmarshal([]byte(jsonText), &apiData)
		
        contentTypeConfig = make(map[string]interface{})
        jsonText = GetContentFromFile("contentType.conf")
        json.Unmarshal([]byte(jsonText), &contentTypeConfig)	
				
}

func main(){
    HandleConfig()
    StartServer()
}