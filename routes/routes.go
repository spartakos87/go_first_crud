package routes

import (
	"../database"
	"../encodecodepass"
	"../jwthandler"
	"../type_structure"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
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

func createUser(w http.ResponseWriter, r *http.Request) {
	db := database.OpenConnection()
	reqBody, _ := ioutil.ReadAll(r.Body)
	var temp_user type_structure.Users
	json.Unmarshal(reqBody, &temp_user)
	hash_pass, err := encodecodepass.HashPassword(temp_user.Pass)
	if err != nil {
		panic(err)
	}
	sql_insert := "INSERT INTO Users (username, pass) VALUES ($1, $2)"
	_, err = db.Exec(sql_insert, temp_user.Username, hash_pass)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	defer db.Close()

}

func logInUser(w http.ResponseWriter, r *http.Request) {
	db := database.OpenConnection()
	reqBody, _ := ioutil.ReadAll(r.Body)
	var temp_user, db_user type_structure.Users
	json.Unmarshal(reqBody, &temp_user)
	sql_select := "SELECT * FROM Users WHERE username=$1"
	row := db.QueryRow(sql_select, temp_user.Username)
	err := row.Scan(&db_user.Username, &db_user.Pass)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	if encodecodepass.CheckPassWord(temp_user.Pass, db_user.Pass) == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		token, _ := jwthandler.GenerateJWT(temp_user.Username)
		token_str, _ := json.MarshalIndent(token, "", "\t")
		w.Write(token_str)
	}
	defer db.Close()
}

func verifyJWT(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var temp_token type_structure.Jwt
	json.Unmarshal(reqBody, &temp_token)
	_, err := jwthandler.ValidateJWT(temp_token.Token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}
	w.WriteHeader(http.StatusOK)
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

func uploadFile(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	// Create the uploads file if it doent already exists
	err = os.Mkdir("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Create a new file in the uploads directory
	dst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Upload successful")
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, "./uploads/1623414617513232878.png")
}

func HandleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAllArticles)
	myRouter.HandleFunc("/all_db", returnAllArticlesFromDB)
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/creat_db_article", createNewArticleInDB).Methods("POST")
	myRouter.HandleFunc("/sign_up", createUser).Methods("POST")
	myRouter.HandleFunc("/log_in", logInUser).Methods("POST")
	myRouter.HandleFunc("/token", verifyJWT).Methods("POST")
	myRouter.HandleFunc("/uploadfile", uploadFile).Methods("POST")
	myRouter.HandleFunc("/downloadfile/{filename}", downloadFile).Methods("GET")
	myRouter.HandleFunc("/article", updateArticle).Methods("PUT")
	myRouter.HandleFunc("/article_db", updateArticleFromDB).Methods("PUT")
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/article_db/{id}", deleteArticleFromDB).Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", returnSignleArticle)
	myRouter.HandleFunc("/article_db/{id}", returnSignleArticleFromDB)
	log.Fatal(http.ListenAndServe(":1313", myRouter))

}
