 package main

 import (
  "encoding/json"
  "fmt"
  "net/http"
  "os"
 )


 type Logs struct {
 Name string `json:"name"`
 Lname string `json:"lname"`
}


func write_data(w http.ResponseWriter, r *http.Request){
 fmt.Println("You got me")
 if r.Method != http.MethodPost{
  return
 }
 r.ParseForm()

 data,_ :=  os.ReadFile("logs.json")
 var l []Logs
 json.Unmarshal(data, &l)



 file,_ := os.Create("logs.json")
 defer file.Close()
 //long version of it
 //encoder := json.NewEncoder(file)
 //encoder.Encode(task)
 logged_names := Logs{

   Name:  r.FormValue("name"),
   Lname: r.FormValue("lname"),

 }
 l = append(l, logged_names)
 encoder := json.NewEncoder(file)
 encoder.SetIndent("", " ")
 encoder.Encode(l)
 fmt.Println("Data is saved successfully!")
 read_data()
}

func read_data(){
 file,_ := os.Open("logs.json")
 defer file.Close()

 var t []Logs

 decoder := json.NewDecoder(file)
 decoder.Decode(&t)

 fmt.Println("")
 fmt.Println("Now we're reading trough this data!")
 fmt.Println(t)


}


func home(w http.ResponseWriter, r *http.Request){
     http.ServeFile(w,r, "index.html")
}

 func main () {

port := os.Getenv("PORT")
if port == "" {
 port = "8080"
}
  http.HandleFunc("/", home)
  http.HandleFunc("/write", write_data)

http.ListenAndServe(":"+port, nil)
}