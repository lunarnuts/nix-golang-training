package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type Post struct {
	UID   int    `json:"userId"`
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (p Post) String() string {
	return fmt.Sprintf(""+
		"{\n"+
		"    \"userId\": %d,\n"+
		"    \"id\": %d,\n"+
		"    \"title\": \"%s\",\n"+
		"    \"body\": \"%s\"\n"+
		"  }", p.UID, p.ID, p.Title, p.Body)
}

func GetPost(pid int, wg *sync.WaitGroup) {
	defer wg.Done()
	var post Post
	resp, err := http.Get(fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", pid))
	if err != nil {
		log.Printf("Error occured: %v\n", err)
		return
	}
	//defer resp.Body.Close()
	if err != nil {
		log.Printf("Error occured: %v\n", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error occured: %v\n", err)
		return
	}
	err = json.Unmarshal(body, &post)
	if err !=nil {
		log.Printf("Error occured: %v",err)
		return
	}
	fmt.Println(post.String())
	SaveToFile(post)
	return
}

func SaveToFile(post Post) {
	id := post.ID
	message := []byte(post.String())
	err := ioutil.WriteFile(fmt.Sprintf("task03/storage/posts/%d.txt", id), message, 0666)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}
	return
}

func main() {
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go GetPost(i, &wg)
	}
	wg.Wait()
}
