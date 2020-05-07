package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
    "database/sql"
  _ "github.com/go-sql-driver/mysql"	

    "io"

)

type Comment struct {
	Id        int  `json:"id"`
	Name      string  `json:"name"`
	Comment   string  `json:"comment"`
	Parent_id int  `json:"parent_id"`
	Inserted  string  `json:"inserted"`
	Children  []Comment `json:"children"`
}

var db *sql.DB
var err error
var Info    *log.Logger


func returnAllComments(w http.ResponseWriter, r *http.Request) {
    //results, err := db.Query("WITH recursive cte (id, name, parent_id, comment, inserted) as ( select id, name, parent_id, comment, inserted from comments where parent_id = 0 union all select c.id, c.name, c.parent_id, c.comment, c.inserted from comments c inner join cte on c.parent_id = cte.id ) select * from cte")
    results, err := db.Query("SELECT * FROM comments")
    if err != nil {
        panic(err.Error())
    }

    var comments = []Comment{}
    for results.Next() {
        var comment Comment
        // for each row, scan the result into our tag composite object
        err = results.Scan(&comment.Id, &comment.Name, &comment.Comment, &comment.Parent_id, &comment.Inserted)
        if err != nil {
            panic(err.Error()) // proper error handling instead of panic in your app
        }
        
        comments = append(comments, comment)
        
    }

    comments = getChildren(comments, 0, 0)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(comments)
}

func getChildren(results []Comment, parent_id int, level int) []Comment {
    var comments = []Comment{}
    level++;
    for i := 0; i < len(results); i++ {
        var result = results[i];

        //Info.Println("ppp")
        // Info.Println(parent_id)
        if(result.Parent_id == parent_id) {
            //result.level = level;
            
            var children []Comment = getChildren(results, result.Id, level)

            if(len(children) > 0) {
                result.Children = children
            }
            comments = append(comments, result)
        }
    }
    
    return comments;
}



func createNewComment(w http.ResponseWriter, r *http.Request) {
    
    var comment Comment
    json.NewDecoder(r.Body).Decode(&comment)

    result, err := db.Exec(`INSERT INTO comments (name, comment, parent_id, inserted) VALUES (?, ?, ?, NOW())`, comment.Name, comment.Comment, comment.Parent_id)
    if err != nil {
        println(err.Error())
    }
    id, err := result.LastInsertId()
    if err != nil {
        println("Error:", err.Error())
    } else {
        println("LastInsertId:", id)
    }

    results := db.QueryRow("SELECT * FROM comments where id = ?", id)


    var comment2 Comment
    // for each row, scan the result into our tag composite object
    err = results.Scan(&comment2.Id, &comment2.Name, &comment2.Comment, &comment2.Parent_id, &comment2.Inserted)
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(comment2)

}

func handleRequests() {
	db, err = sql.Open("mysql", "edge:fusionha87@tcp(127.0.0.1:3306)/coding_challenge1")
	if err != nil {
        panic(err.Error())
    }
    myRouter := mux.NewRouter().StrictSlash(true)
    
    myRouter.HandleFunc("/comments", returnAllComments).Methods("GET")
    myRouter.HandleFunc("/comments", createNewComment).Methods("POST")

    corsObj:=handlers.AllowedOrigins([]string{"*"})

    log.Fatal(http.ListenAndServe(":8080", handlers.CORS(corsObj)(myRouter)))
}


func Init(infoHandle io.Writer) {



    Info = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

}

func main() {

    db, err = sql.Open("mysql", "edge:fusionha87@tcp(127.0.0.1:3306)/coding_challenge1")
	if err != nil {
        panic(err.Error())
    }
    handleRequests()
}
