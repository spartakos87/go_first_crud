package routes

import (
	"FirstRestGo/type_structure"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

var Articles = []type_structure.Article{
	type_structure.Article{Title: "Hello", Desc: "Article Description", Content: "Article Content", Id: "1"},
	type_structure.Article{Title: "Hello 2", Desc: "Article Description", Content: "Article Content", Id: "2"},
}

func homePage(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintln(w, "Welcome to the Homepage!")
	fmt.Println("Endpoint Hit: homePage")

}

func returnAllArticles(w http.ResponseWriter, r *http.Request){
	fmt.Println("Endpoint Hit: returnAllArticles")

	json.NewEncoder(w).Encode(Articles)
}

func returnSignleArticle(w http.ResponseWriter, r *http.Request)  {

	vars := mux.Vars(r)
	key := vars["id"]
	for _, article := range Articles{
		if article.Id == key{
			json.NewEncoder(w).Encode(article)
		}

	}

}

func createNewArticle(w http.ResponseWriter, r *http.Request)  {
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article type_structure.Article
	json.Unmarshal(reqBody, &article)
	Articles = append(Articles, article)
	fmt.Println(article)
	fmt.Fprintf(w, "%+v", string(reqBody))
}

func deleteArticle(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	id := vars["id"]
	for index, article := range Articles{
		if article.Id == id {
			Articles = append(Articles[:index], Articles[index+1:]...)
		}
	}

}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article type_structure.Article
	json.Unmarshal(reqBody, &article)
	fmt.Println(article)
	for index, temp_article := range Articles {
		if temp_article.Id == article.Id {
			fmt.Println("In if statement")
			Articles[index] = article
		}

	}
}

func HandleRequests(){
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAllArticles)
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/article", updateArticle).Methods("PUT")
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", returnSignleArticle)
	log.Fatal(http.ListenAndServe(":1313", myRouter))

}