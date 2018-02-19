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
    //con, err := sql.Open("mysql", store.user+":"+store.password+"@/"+store.database)
    //db, err := sql.Open("mysql", "root:admin@/netapp?charset=utf8")
    db, err := getDB()
        checkErr(err)
        
    rows, err, listSize := getRows(db)
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
    
    // hit the website, get the 'body', get length of the text
    r.ParseForm()
    
    url := r.Form.Get("url")
    fmt.Println("url:", url)
    
    message := "something went wrong."
    
    resp, err := http.Get(url)
    if err != nil {
      message = "error on http.Get:"
      fmt.Println(message, err) 
    } else {
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
          message = "error on ioutil.ReadAll:"
          fmt.Println(message, err)
        } else {
            count := len(string(body))
            countAsString := strconv.Itoa(count)
            fmt.Println("char count = ", countAsString)
            message = countAsString
          } 
      }
          
    //  update database with new'crawl' data
    db, err := getDB()
        
    stmt, err := db.Prepare("INSERT crawler SET url=?,count=?")
        checkErr(err)
        
    res, err := stmt.Exec(url, message)
        checkErr(err)
        
    _ = res // dont plan on using this var so... go requires me to make that clear
        
    // finally, get data from crawler table to populate view, build view
    rows, err, listSize := getRows(db)
        checkErr(err)
        
    db.Close()
    defer rows.Close()
    
    // make a slice
    crawlerList := make([]CrawlerData, listSize, 2*listSize)
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
    
    // update view
    t, err := template.ParseFiles("templates/crawler.html")
    err = t.Execute(w, crawlerList)
    
    /*
    fmt.Fprintf(w, "\n")
    fmt.Fprintf(w, string(body))
    fmt.Fprintf(w, "\n")
    fmt.Fprintf(w, "\n")
    */
}

//                                                                                   
func getRows(db *sql.DB) (*sql.Rows, error, int){
  rows, err := db.Query("SELECT * FROM crawler")
  
  row := db.QueryRow("SELECT count(*) as numRows FROM crawler")
  
  numRows:=0
  row.Scan(&numRows)
  
  fmt.Println("numRows = ", numRows)
  
  return rows, err, numRows
}

//                                                                              
func getDB() (*sql.DB, error){
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