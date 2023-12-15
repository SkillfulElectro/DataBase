package main

import (
	"fmt"
  "net/http"
  "net"
  "time"
  "io" 
)

var bot1 string
var bot2 string
var Chat_ID string
var data_count int = 0
var t_time int = 0

func timer(){
  for{
    time.Sleep(time.Second)
    t_time++
  }
}

func getData_All(){
  
}

func getData_byObject(){
  
}

func chatID_set(response string){

  i:= 26
  for ;response[i] != ',' ; i++{}
  
	Chat_ID = response[26:i]
}

func handleRequest(db_conn net.Conn){
  res , err := io.ReadAll(db_conn)

  if err != nil{
    fmt.Fprintf(db_conn , "Error")
    fmt.Println("Error happened on reciving data :%s" , err)
    return
  }

  data := string(res)
  if data[0] == 's'{
    _ , err = http.Get(fmt.Sprintf(bot2 + "/sendMessage?chat_id=%s&text=%s" , Chat_ID , data[1:] ))
    data_count++
  } 
  if err != nil{
    fmt.Fprintf(db_conn , "Failed to send data : %s" , err)
  }

  if (data[0] == 'g'){
     switch (data[1]){
        case 'a':
        getData_All()
        case '/':
        getData_byObject()
    }
  }
}

func startingDB(){
  fmt.Println("Starting DB servers ...")
  time.Sleep(time.Second*2)
  fmt.Println("checking available ports ... pls wait for a while")

  var db_serv net.Listener
  var err error
  for i := 0;i<20000;i++ {
    db_serv , err = net.Listen("tcp" , fmt.Sprintf("localhost:%d" , i))

    if err != nil {
      continue 
    }else{
      fmt.Println("DB started on port%d" , i)
      break
    }
  }

  if err != nil {
    fmt.Println("all ports on your device are filled so DB failed to start")
    return
  }

  fmt.Println("you can interact with server in a tcp connection\n1. for sending data to DB use this syntax : s/<your data>")
   
  for {
    conn, err := db_serv.Accept()
    if err != nil {
      fmt.Println("Error:", err)
      return
    }
    go handleRequest(conn)
  }
}

func main() {
  go timer()
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
