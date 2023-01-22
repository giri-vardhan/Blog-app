package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	//"golang.org/x/tools/go/analysis/passes/nilfunc"
)

var DB = connectDB()

func connectDB() *sql.DB {

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	//connStr := "postgres://:@127.0.0.1:5433/giri?sslmode=disable"

	DB, err := sql.Open("postgres", psqlconn)
	//CheckError(err)
	if err != nil {
		fmt.Println(err.Error())
		log.Println("init error")
		panic(err)
	}
	// fmt.Println("yes")
	// // db, err := sql.Open("postgres", psqlconn)
	// // CheckError(err)

	err = DB.Ping()
	CheckError(err)
	fmt.Println("successfully connected")
	return DB
}

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "giriopen"
	dbname   = "giri"
)

type post struct {
	PostID      int    `json:"postid"`   // primary key
	Title       string `json:"title"`
	UserName    string `json:"username`
	Description string `json:"description"`
	PostTime    string `json:"posttime"`
}

type comment struct {
	CommentID          string `json:"commentid`    // primary key
	CommentPostID      string `json:"commentpostid"`  //foreign key  referencing to primary key PostID in post table
	CommentedUser      string `json: commenteduser`
	CommentTime        string `json:commenttime`
	CommentDescription string `json:commentdescription`
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	fmt.Println("Successfully connected")

}

func getPost(c *gin.Context) {    // fetching posts from database
	
	// var posts []post
	// for rows.Next() {
	// 	var PostID int
	// 	var Title string
	// 	var UserName string
	// 	var Description string
	// 	var PostTime string
	// 	err = rows.Scan(&PostID, &Title, &UserName, &Description, &PostTime)
	// 	if err != nil {
	// 		log.Println("scan in get post")
	// 		//c.IndentedJSON(http.StatusConflict, gin.H{"message": "scan in get post", "error": err})
	// 	}
	// 	posts = append(posts, post{PostID: PostID, Title: Title, UserName: UserName, Description: Description, PostTime: PostTime})
	// }
	// fmt.Println(posts)
	//c.IndentedJSON(http.StatusOK, posts)
	rows, err := DB.Query("select * from post")
	CheckError(err)
	if err != nil {
		log.Println("not able to fetch data from get post")
	}
	posts:=make([]post,0)

	// var PostID string
	// var Title string
	// var UserName string
	// var Description string
	// var PostTime string
	for rows.Next() {
		p := post{}
		err = rows.Scan(&p.PostID, &p.Title, &p.UserName, &p.Description, &p.PostTime)
		if err != nil {
			log.Println("scan in get post")
		}
		posts = append(posts, p)
	}
	fmt.Println(posts)
	 c.IndentedJSON(http.StatusOK, posts)

}

func getCommentByPostID(c *gin.Context) {  //fetching comment by postid
// 	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
// 	db, err := sql.Open("postgres", psqlconn)
// 	CheckError(err)

// 	// get the id from the url
	id := c.Param("id")
	// fmt.Println(id)
	// get data from database
	// return data
	 d, _ := strconv.Atoi(id)
	
	rows, err := DB.Query("SELECT * FROM comment WHERE PostID = $1",d)
	CheckError(err)
	// temp := make([]comment, 0)
	// for rows.Next() {
	// 	c := comment{}
	// 	err = rows.Scan(&c.CommentID, &c.PostID, &c.CommentedUser, &c.CommentTime,&c.CommentDescription)
	// 	CheckError(err)
	// 	temp = append(temp, c)
	// }

	var temp comment
	for rows.Next() {
		var commentid string
		var commentpostid string
		var commenteduser string
		var commenttime  string
		var commentdescription string
		err = rows.Scan(&commentid, &commentpostid, &commenteduser, &commenttime , &commentdescription)
		CheckError(err)
		temp = comment {CommentID: commentid, CommentPostID: commentpostid, CommentedUser: commenteduser, CommentTime: commenttime, CommentDescription: commentdescription}
	}
	c.IndentedJSON(http.StatusOK, temp)
}

func createComment(c *gin.Context) {  // adding comment to a post
	var newComment comment

	if err := c.BindJSON(&newComment); err != nil {
		return
	}

	insert := `insert into "Comment" ("CommentID","PostID","CommentedUser","CommentTime","CommentDescription") values ($1,$2,$3,$4,$5)`
	d := time.Now().Format("2006-01-02 15:04:05")
	newComment.CommentTime = d

	_, err := DB.Exec(insert, newComment.CommentID, newComment.CommentPostID, newComment.CommentedUser, newComment.CommentTime, newComment.CommentDescription)
	CheckError(err)

	c.IndentedJSON(http.StatusCreated, newComment)
}

func CreatePost(c *gin.Context) {   // creating a new post

	var newPost post

	if err := c.BindJSON(&newPost); err != nil {
		return
	}

	insert := `insert into "post" ("PostID","Title","UserName","Description","PostTime") values ($1,$2,$3,$4,$5)`

	// "d := time.Now().Format("2006-01-02 15:04:05")
	// newPost.PostTime = d

	_, err := DB.Exec(insert, newPost.PostID, newPost.Title, newPost.UserName, newPost.Description, newPost.PostTime)
	CheckError(err)

	c.IndentedJSON(http.StatusCreated, newPost)
}

func main() {
	// getPost()
	router := gin.Default()
	router.GET("/posts", getPost)
	router.POST("/posts", CreatePost)
	router.GET("/comment/:id", getCommentByPostID)
	router.POST("/comment", createComment)
	router.Run("localhost:8080")
}

// curl localhost:8080/posts -i -H "Content-Type: application/json" -d @body.json --request "POST"


/*
curl http://localhost:8080/comment \
--include \
--header "Content-Type: application/json" \
--request "POST" \
--data @comment.json
*/

/*
curl http://localhost:8080/posts --include --header "Content-Type: application/json" --data @body.json --request "POST"
*/

/*
curl localhost:8080/comment/2
curl localhost:8080/comment --include --header "Content-Type: application/json" --data @comment.json --request "POST"
*/

