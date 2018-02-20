package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "strconv"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "log"
    "html/template"
)

type CrawlerData struct {
  Url   string
  Count string
}

func main() {
    http.HandleFunc("/", homePage)
    http.HandleFunc("/doCrawler", doCrawler)
    http.ListenAndServe(":8080", nil)
}

//                                                                                           
func homePage(w http.ResponseWriter, r *http.Request){
    db, err := getDB()
        checkErr(err)
        
    rows, listSize, err := getRows(db)
        checkErr(err)
    
    db.Close()
    defer rows.Close()
        
    fmt.Println("listSize", listSize)
    fmt.Println("rows", rows)
    
    // make a slice
    crawlerList := make([]CrawlerData, listSize, 2*listSize)
    index := 0

    // print out some useful info server side - load the actual crawlerList
    for rows.Next() {
            var uid int
            var url string
            var count string
            err = rows.Scan(&uid, &url, &count)
            checkErr(err)
            
            fmt.Print("url = ", url +" ")
            fmt.Println("count = ", count)
            
            entry := CrawlerData{
                Url:   url,
                Count: count,
            }
            
            crawlerList[index] = entry
            index++
        }
    
    // go templates are cool. you will like them
    t, err := template.ParseFiles("templates/crawler.html") 
    if err != nil { 
  	  log.Print("template parsing error: ", err) 
  	}
  
    err = t.Execute(w, crawlerList)
    if err != nil { 
  	  log.Print("template executing error: ", err) 
  	}
}

//                                                                                           
func doCrawler(w http.ResponseWriter, r *http.Request){
  
  //to do - get the db handle first. if that doesnt work? rest of code doesnt matter. handle that first
  
  // some view stuff we will need, regrdless of errors
  message := "database error, unable to get DB handle. is database running and configured properly?"
  var crawlerList []CrawlerData
  
  //                                                                            
  db, err := getDB()
  if err != nil { 
    // we have a DB error, generate view info --- no more processing needed
    crawlerList = generateDbErrorData(message)
  } else {
      url, count, err := handleFormAndUrl(r)
      if err != nil{
        // we have a form/url error
        fmt.Println(message, err)
        crawlerList = generateDbErrorData("we have a form/url error")
      } else {
          // all is well, try to do insert 
          message, err = doInsert(db, url, count)
          if err != nil {
            fmt.Println(message, err)
            crawlerList = generateDbErrorData(message)
          } else {
              // finally, get data from crawler table to populate view, build view
              rows, listSize, err := getRows(db)
              if err != nil {
                fmt.Println(message, err)
                message = "error on func getRows"
                crawlerList = generateDbErrorData(message)
              } else {
                  // make a slice
                  crawlerList = make([]CrawlerData, listSize, 2*listSize)
                  index := 0
              
                  for rows.Next() {
                    var uid int
                    var url string
                    var count string
                    err = rows.Scan(&uid, &url, &count)
                    checkErr(err)
                                
                    entry := CrawlerData{
                        Url:   url,
                        Count: count,
                    }
                    
                    crawlerList[index] = entry
                    index++
                  }
                } // end retrieve rows worked
              rows.Close()
            } // end insert worked
        } // end form - url stuff succesful
        
      db.Close()
      
    } // end db work    
    
  
  // finally, update view
  t, err := template.ParseFiles("templates/crawler.html")
  err = t.Execute(w, crawlerList)
}

//                                                                                
func handleFormAndUrl(r *http.Request) (string, string, error){
  // hit the website, get the 'body', get length of the text
  message := "no worries"
  
  r.ParseForm()
  
  url := r.Form.Get("url")
  fmt.Println("url:", url)
  
  resp, err := http.Get(url)
  if err != nil {
    message = "error on http.Get"
    fmt.Println(message, err) 
  } else {
      defer resp.Body.Close()
      body, err := ioutil.ReadAll(resp.Body)
      if err != nil {
        message = "error on ioutil.ReadAll"
        fmt.Println(message, err)
      } else {
          count := len(string(body))
          countAsString := strconv.Itoa(count)
          fmt.Println("char count = ", countAsString)
          message = countAsString
        } 
    }
  
  return url, message, err
}

//                                                      
func doInsert(db *sql.DB, url string, countAsString string) (string, error){
  stmt, err := db.Prepare("INSERT crawler SET url=?,count=?")
  errMessage := "coding error --- you should never see this message :). method = doInsert"
  
  if err != nil {
    errMessage = "error on db.Prepare"
    fmt.Println(errMessage, err) 
    } else {
        res, err := stmt.Exec(url, countAsString)
        _ = res // dont plan on using this var so... go requires me to make that clear
        if err != nil {
          errMessage = "error on stmt.Exec"
          fmt.Println(errMessage, err)
        }
      }
  
  return errMessage, err
}

//                                                                                   
func generateDbErrorData(message string) ([]CrawlerData){
  crawlerList := make([]CrawlerData, 1, 1)
  
  entry := CrawlerData{
            Url:   message,
            Count: "count",
  }
  
  crawlerList[0] = entry
  
  return crawlerList
}

//                                                                                   
func getRows(db *sql.DB) (*sql.Rows, int, error){
  rows, err := db.Query("SELECT * FROM crawler")
  
  row := db.QueryRow("SELECT count(*) as numRows FROM crawler")
  
  numRows:=0
  row.Scan(&numRows)
  
  fmt.Println("numRows = ", numRows)
  
  return rows, numRows, err
}

//                                                                              
func getDB() (*sql.DB, error){
  //con, err := sql.Open("mysql", store.user+":"+store.password+"@/"+store.database)
  db, err := sql.Open("mysql", "root:admin@/netapp?charset=utf8")
  
  return db, err
}


//                                                                                           
func checkErr(err error) {
        if err != nil {
            fmt.Println("error = ", err)
            panic(err)
        }
}