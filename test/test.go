// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"my-rest-api/gopkg.in/mgo.v2/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gorilla/mux"
)

// Post - Our struct for all Posts
type Post struct {
	Post_Id          string `json:"Post_Id"`
	Caption          string `json:"Caption"`
	Image_URL        string `json:"Image_URL"`
	Posted_timestamp string `json:"Posted_timestamp"`
	Post_User_Id     string `json:"Post_User_Id"`
}

type User struct {
	User_Id  string `json:"User_Id"`
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

var Posts []Post
var Users []User

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["User_Id"]
	for _, post := range Posts {
		if post.Post_User_Id == key {
			json.NewEncoder(w).Encode(post)
		}
	}
}

func returnSinglePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["Post_Id"]

	for _, post := range Posts {
		if post.Post_Id == key {
			json.NewEncoder(w).Encode(post)
		}
	}
}

func createNewPost(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Post struct
	// append this to our Posts array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var post Post
	json.Unmarshal(reqBody, &post)
	// update our global Posts array to include
	// our new Post
	Posts = append(Posts, post)

	json.NewEncoder(w).Encode(post)
}

func createNewUser(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Post struct
	// append this to our Posts array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user User
	json.Unmarshal(reqBody, &user)
	// update our global Posts array to include
	// our new Post
	Users = append(Users, user)

	json.NewEncoder(w).Encode(user)
}

func returnSingleUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["Email"]
	password := vars["Password"]
	for _, user := range Users {
		if user.Email == email {
			if user.Password == password {
				json.NewEncoder(w).Encode(user)
			}
		}
	}
}
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/posts/users/{User_Id}", returnAllPosts).Methods("GET")
	myRouter.HandleFunc("/posts", createNewPost).Methods("POST")
	myRouter.HandleFunc("/posts/{Post_Id}", returnSinglePost).Methods("GET")
	myRouter.HandleFunc("/users", createNewUser).Methods("POST")
	myRouter.HandleFunc("/users/{User_Id}", returnSingleUser).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://yash:manuYASH1234*@cluster0.sy5pe.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
	defer client.Disconnect(ctx)
	Posts = []Post{
		Post{Post_Id: "1", Caption: "Hello", Image_URL: "abffc", Posted_timestamp: "Post Posted_timestamp", Post_User_Id: "1"},
		Post{Post_Id: "2", Caption: "Hello 2", Image_URL: "adbeif324", Posted_timestamp: "Post Posted_timestamp", Post_User_Id: "1"},
	}
	Users = []User{
		User{User_Id: "1", Name: "Yash", Email: "yashsinha@gmail.com", Password: "qwdcfghd"},
	}

	handleRequests()

}
