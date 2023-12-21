package main

import (
	"fmt"
  "net/http"
  "net"
  "time"
  "strconv"
  "io" 
)

var bot1 string
var bot2 string
var Chat_ID string
var data_count int = 0
var t_time int = 0

func FindinString(data , text string) ([]int , int){
  
  i:=0
  for ;i<len(data);i++ {
    if(i+len(text) > len(data)){
      return []int{-1},-1
    }
    to_check := (data)[0+i:len(text)+i]
    if (to_check == text){
      break
    }  
  }

  if(i+len(text) > len(data)){
    return []int{-1},-1
  }
  
  return []int{i, i+len(text)},0
}

func getUpdates() string {
  res , err:=http.Get(fmt.Sprintf("%sgetUpdates?offset=%d&limit=1" , bot1 , 1))
  if err != nil {
    return fmt.Sprintf("Failed to get data : %s" , err)
  }
  body , _ := io.ReadAll(res.Body)
  body_s := string(body)

  fmt.Println(body_s)
  
  message_id , bumpy :=FindinString(body_s , "\"message_id\":")
  if(bumpy == -1){
    return "DB is empty !"
  }
  update_id , _:=FindinString(body_s , "\"update_id\":")

 // fmt.Println(message_id, update_id)

  message_id[0] = message_id[1]

  for{
     if(body_s[message_id[0]] == ',')      {
       break
     }
     message_id[0]++
  }

  mes_id , _ := strconv.Atoi(body_s[message_id[1]:message_id[0]])
  
  update_id[0] = update_id[1]
  for{
    if(body_s[update_id[0]] == ',')       {
       break
    }
     update_id[0]++
  }
  upd_id , _ := strconv.Atoi(body_s[update_id[1]:update_id[0]])

  fmt.Println(mes_id , upd_id)
  
    res , _ =http.Get(fmt.Sprintf("%sgetUpdates?offset=%d&limit=1" , bot1 , upd_id))

  upd_id++
  http.Get(fmt.Sprintf("%sgetUpdates?offset=%d&limit=1" , bot1 , upd_id))
    http.Get(fmt.Sprintf("%sdeleteMessage?chat_id=%s&message_id=%d" , bot2 , Chat_ID , mes_id))

  //body , _ = io.ReadAll(res.Body)
  //fmt.Println(string(body) , err)
    
    body , _ = io.ReadAll(res.Body)
    body_s = string(body)
    
    value , _:= FindinString(body_s , "\"text\":\"")
    
    value[0] = len(body_s) - 1
    for ;;value[0]-- {
      if (body_s[value[0]] == '}'){
        break
      }
    }
    value[0] = value[0] - 4
    data := body_s[value[1]:value[0]]

   // fmt.Println(data)
  http.Get(fmt.Sprintf("%ssendMessage?chat_id=%s&text=%s" , bot2 , Chat_ID , data))
  return data
}

func timer(){
  for{
    time.Sleep(time.Second*3)
    t_time++
    getUpdates()
  }
}

func getData_All(db_conn net.Conn){
  num := data_count
  for ;num>0;num--{
    fmt.Fprintf(db_conn , getUpdates())
    fmt.Fprintf(db_conn , "\n")
  }
}

func getData_byObject(db_conn net.Conn , object string){
  num := data_count
  for ;num>0;num--{
    data := getUpdates()
    _ , err := FindinString(data , object)
    if err != -1{
      fmt.Fprintf(db_conn, data)
    }
  }
}

func chatID_set(response string){

  i:= 26
  for ;response[i] != ',' ; i++{}
  
	Chat_ID = response[26:i]
}

func handleRequest(db_conn net.Conn){
  fmt.Fprintf(db_conn , "Connected")
  res , err := io.ReadAll(db_conn)

  if err != nil{
    fmt.Fprintf(db_conn , "Error")
    fmt.Println("Error happened on reciving data :%s" , err)
    return
  }

  data := string(res)
  if data[0] == 's'{
    _ , err = http.Get(fmt.Sprintf(bot2 + "/sendMessage?chat_id=%s&text=%s" , Chat_ID , data[1:] ))
    if err != nil{
      fmt.Fprintf(db_conn , "Failed to send data : %s" , err)
    }else{
      data_count++
      fmt.Fprintf(db_conn , "Success")
    }
  } 

  if (data[0] == 'g'){
     switch (data[1]){
        case 'a':
        getData_All(db_conn)
        case '/':
        next := 1
        for ;;next++{
          if (data[next] == '/'){
            break
          }
        }
        getData_byObject(db_conn , data[1:next])
    }
  }
}

func startingDB(){
  fmt.Println("Starting DB servers ...")
  time.Sleep(time.Second*2)
  fmt.Println("checking available ports ... pls wait for a while")

  var db_serv net.Listener
  var err error
  for i := 1000;i<20000;i++ {
    db_serv , err = net.Listen("tcp" , fmt.Sprintf("localhost:%d" , i))

    if err != nil {
      continue 
    }else{
      fmt.Println("DB started on port : " , i)
      break
    }
  }

  if err != nil {
    fmt.Println("all ports on your device are filled so DB failed to start")
    return
  }

  fmt.Println("you can interact with server in a tcp connection\n1. for sending data to DB use this syntax : s/<your data>\n for getting data from DB use this syntax : g/<object>/ or if you want all of data you have ever sent to DB: ga/\n\n")
  fmt.Println("for making everything safe always encode data and also for being able to use getting by object always give special objects to them or id")
  go timer()
  for {
    conn, err := db_serv.Accept()
    if err != nil {
      fmt.Println("Error:", err)
      continue
    }
    go handleRequest(conn)
  }
}

func main() {
  bale := "https://tapi.bale.ai/"
   telegram:="https://api.telegram.org/"
  
  fmt.Println("Hello welcome to world of ultimate storage databases\n so how does it work ?\n1. creating account in Telegram or Bale\n2. go to @BotFather and generate your bot tokens (we need two of them)\n3. create a channel to store data in it and add that bots there\n4. insert tokens of your bots and username of channel without @ here !\n let's begin our journey ;)\n")

  var platform string
  fmt.Println(">>Bale or Telegram : (B Or T)")
  fmt.Scanln(&platform)

  

  fmt.Println("insert token of your first bot :\n")
  fmt.Scanln(&bot1)

  fmt.Println("insert token of your second bot :\n")
  fmt.Scanln(&bot2)

  
  switch (platform){
    case "B":
      bot1 = bale + "bot" + bot1 + "/"
      bot2 = bale + "bot" + bot2 + "/"
    case "T":
      bot1 = telegram + "bot" + bot1 + "/"
      bot2 = telegram + "bot" + bot2 + "/"
  }

  var username string
  fmt.Println("insert username of your channel :")
  fmt.Scanln(&username)
  
  res , err := http.Get(fmt.Sprintf(bot1 + "getChat?chat_id=@%s" , username))

  if err != nil{
    fmt.Println("Error happened in getting channel information ! : %s" , err)
    return
  }

  body , _ := io.ReadAll(res.Body)
 chatID_set(string(body))

  startingDB()
}
