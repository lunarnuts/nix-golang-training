package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Post struct {
	UID int `json:"userId"`
	ID  int `json:"id"`
	Title string `json:"title"`
	Body string `json:"body"`
}

func (p Post) String() string {
	return fmt.Sprintf("" +
		"{\n" +
		"    \"userId\": %d,\n" +
		"    \"id\": %d,\n" +
		"    \"title\": \"%s\",\n" +
		"    \"body\": \"%s\"\n" +
		"  }",p.UID,p.ID,p.Title,p.Body)
}

func main() {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")
	if err!=nil {
		log.Fatalf("Error occured: %v\n",err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		log.Fatalf("Error occured: %v\n",err)
	}
	var Posts []Post
	json.Unmarshal(body,&Posts)
	for _, post := range Posts {
		fmt.Print(post.String()+",\n")
	}
}
