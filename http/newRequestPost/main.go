package main

type NewPost struct {
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func CreatePost(url string, newPost NewPost) (Post, error) {

}

func main() {
	url := "https://jsonplaceholder.typicode.com/posts"

}
