package routes

import (
	"../database"
	"../type_structure"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

var Articles = []type_structure.Article{
	type_structure.Article{Title: "Hello", Desciption: "Article Description", ContentArtcile: "Article Content", Id: "1"},
	type_structure.Article{Title: "Hello 2", Desciption: "Article Description", ContentArtcile: "Article Content", Id: "2"},
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Homepage!")
	fmt.Println("Endpoint Hit: homePage")

}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")

	json.NewEncoder(w).Encode(Articles)
}

func returnAllArticlesFromDB(w http.ResponseWriter, r *http.Request) {
	db := database.OpenConnection()
	rows, err := db.Query("SELECT * FROM Article")
	if err != nil {
		log.Fatal(err)
	}
	var articles []type_structure.Article
	for rows.Next() {
		var temp_article type_structure.Article
		rows.Scan(&temp_article.Id, &temp_article.Title, &temp_article.Desciption, &temp_article.ContentArtcile)
		articles = append(articles, temp_article)
	}
	if len(articles) != 0 {
		artilcesBytes, _ := json.MarshalIndent(articles, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(artilcesBytes)
	} else {
		//msg,_ :=json.MarshalIndent("No data", "", "\t")
		//w.Header().Set("Content-Type", "application/json")
		//w.Write(msg)
		fmt.Fprintln(w, "No Data exist!!!")
	}
	defer rows.Close()
	defer db.Close()
}

func returnSignleArticle(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["id"]
	for _, article := range Articles {
		if article.Id == key {
			json.NewEncoder(w).Encode(article)
		}

	}

}

func returnSignleArticleFromDB(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	db := database.OpenConnection()
	rows, err := db.Query("SELECT * FROM Article WHERE id=$1", id)
	if err != nil {
		log.Fatal(err)
	}

	var articles []type_structure.Article
	for rows.Next() {
		var temp_article type_structure.Article
		rows.Scan(&temp_article.Id, &temp_article.Title, &temp_article.Desciption, &temp_article.ContentArtcile)
		articles = append(articles, temp_article)
	}
	if len(articles) != 0 {
		artilcesBytes, _ := json.MarshalIndent(articles, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(artilcesBytes)
	} else {
		//msg,_ :=json.MarshalIndent("No data", "", "\t")
		//w.Header().Set("Content-Type", "application/json")
		//w.Write(msg)
		fmt.Fprintln(w, "No Data exist!!!")
	}
	defer rows.Close()
	defer db.Close()
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article type_structure.Article
	json.Unmarshal(reqBody, &article)
	Articles = append(Articles, article)
	fmt.Println(article)
	fmt.Fprintf(w, "%+v", string(reqBody))
}

func createNewArticleInDB(w http.ResponseWriter, r *http.Request) {
	db := database.OpenConnection()
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article type_structure.Article
	json.Unmarshal(reqBody, &article)
	sql_insert := "INSERT INTO Article (Id, Title, Desciption, ContentArtcile) VALUES ($1, $2, $3, $4)"
	_, err := db.Exec(sql_insert, article.Id, article.Title, article.Desciption, article.ContentArtcile)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	defer db.Close()

}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	for index, article := range Articles {
		if article.Id == id {
			Articles = append(Articles[:index], Articles[index+1:]...)
		}
	}

}

func deleteArticleFromDB(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	db := database.OpenConnection()
	_, err := db.Exec("DELETE FROM Article WHERE Id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	defer db.Close()
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

func updateArticleFromDB(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article type_structure.Article
	json.Unmarshal(reqBody, &article)
	db := database.OpenConnection()
	_, err := db.Exec("UPDATE Article SET Title = $1 , Desciption = $2 , ContentArtcile = $3 WHERE Id = $4",
		article.Title, article.Desciption, article.Desciption, article.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func HandleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAllArticles)
	myRouter.HandleFunc("/all_db", returnAllArticlesFromDB)
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/creat_db_article", createNewArticleInDB).Methods("POST")
	myRouter.HandleFunc("/article", updateArticle).Methods("PUT")
	myRouter.HandleFunc("/article_db", updateArticleFromDB).Methods("PUT")
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/article_db/{id}", deleteArticleFromDB).Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", returnSignleArticle)
	myRouter.HandleFunc("/article_db/{id}", returnSignleArticleFromDB)
	log.Fatal(http.ListenAndServe(":1313", myRouter))

}
