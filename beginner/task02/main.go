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

func GetPost(pid int, c chan <- string)  {
	var post Post
	resp, err := http.Get(fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d",pid))
	defer resp.Body.Close()
	if err!=nil {
		log.Printf("Error occured: %v\n",err)
		c <- err.Error()
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		log.Printf("Error occured: %v\n",err)
		c <- err.Error()
		return
	}
	json.Unmarshal(body,&post)
	fmt.Println(post.String())
	c <- post.String()
}

func main()  {
	c := make(chan string)
	for i:=1;i<=10;i++ {
		go GetPost(i,c)
	}
	for i:=1;i<10;i++ {
		fmt.Println(<-c)
	}
}