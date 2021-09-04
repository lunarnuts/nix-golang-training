package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type Post struct {
	UID int `json:"userId"`
	ID  int `json:"id"`
	Title string `json:"title"`
	Body string `json:"body"`
}

type Comment struct {
	PostID int `json:"postId"`
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Body string `json:"body"`
}


func (p Post) String() string {
	return fmt.Sprintf("" +
		"{\n" +
		"	\"userId\": %d,\n" +
		"	\"id\": %d,\n" +
		"	\"title\": \"%s\",\n" +
		"	\"body\": \"%s\"\n" +
		"}",p.UID,p.ID,p.Title,p.Body)
}

func GetPosts(uuid int, wg *sync.WaitGroup) (posts []Post)  {
	defer wg.Done()
	resp, err := http.Get(fmt.Sprintf("https://jsonplaceholder.typicode.com/posts?userId=%d",uuid))
	if err!=nil {
		log.Printf("Error occured: %v\n",err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		log.Printf("Error occured: %v\n",err)
		return
	}
	err = json.Unmarshal(body,&posts)
	if err != nil {
		log.Printf("Error occured: %v\n",err)
		return
	}
	for _, post:= range posts {
		wg.Add(1)
		go InsertPostsToDB(post,wg)
	}
	//wg.Wait()
	return posts
}
func InsertPostsToDB(post Post,wg *sync.WaitGroup) {
	defer wg.Done()
	db, err := sql.Open("mysql", "lunarnuts:password@tcp(localhost:3306)/test")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	result, err := db.Exec("Insert into posts(user_id,id,title,body) values (?, ?, ?, ?) ON duplicate key Update id =?",post.UID,post.ID,post.Title,post.Body,post.ID)
	if err!=nil {
		log.Printf("Error occured in posts: %v",err)
		return
	}
	wg.Add(1)
	GetComments(post.ID, wg)
	if r,_:=result.RowsAffected();r == 0 {
		log.Printf("Query affected zero rows: post")
		return
	}
}
func GetComments(pid int, wg *sync.WaitGroup)  {

	defer wg.Done()
	var comments []Comment
	resp, err := http.Get(fmt.Sprintf("https://jsonplaceholder.typicode.com/comments?postId=%d",pid))
	if err!=nil {
		log.Printf("Error occured: %v\n",err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		log.Printf("Error occured: %v\n",err)
		return
	}

	err = json.Unmarshal(body,&comments)
	if err!=nil {
		log.Printf("Error occured: %v\n",err)
		return
	}
	for _, comment := range comments {
		wg.Add(1)
		go InsertCommentsToDB(comment, wg)
	}

}

func InsertCommentsToDB(comment Comment, wg *sync.WaitGroup) {
	defer wg.Done()
	db, err := sql.Open("mysql", "lunarnuts:password@tcp(localhost:3306)/test")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	result, err := db.Exec("Insert into comments(post_id,id,name,email,body) values (?, ?, ?, ?,?) ON" +
		" duplicate key Update id=?",comment.PostID,comment.ID,comment.Name,comment.Email,comment.Body,comment.ID)
	if err!=nil {
		log.Panicf("Error occured: %v",err)
		return
	}
	if r,_:=result.RowsAffected();r == 0 {
		log.Printf("Query affected zero rows: comment")
		return
	}
}



func main()  {
	var wg sync.WaitGroup
	wg.Add(1)
	_ = GetPosts(7, &wg)
	wg.Wait()
}