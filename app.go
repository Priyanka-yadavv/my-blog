package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) error {

	connectionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", user, password, dbname)
	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	checkError(err)
	a.Router = mux.NewRouter().StrictSlash(true)

	a.handleRoutes()
	return nil
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// coverting string to int
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendResponse(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p := product{ID: id}
	if err := p.getProduct(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			sendResponse(w, http.StatusNotFound, "Product not found")
		default:
			sendResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	sendResponse(w, http.StatusOK, p)
}

//	respondWithJSON(w, http.StatusOK, p)
//}

/*
Okay, so before we start with our route handler,
let us create a method that would send back the response,
along with other information such as status code

It would recieve the ResponseWriter,
the status code,
and the payload (i.e the data to be sent back)

as the parameters
*/
func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {

	// Let us use the Marshal method in json to convert our payload into json
	response, _ := json.Marshal(payload)

	// let use set the Header as json format
	w.Header().Set("Content-Type", "application/json")

	// let us also add the status code sent to the Header
	w.WriteHeader(statusCode)

	// and finally writing and sending our response
	w.Write(response)
}

/*
This handler would have the same function signature
- It will have a ResponseWriter
- and a pointer to Request
- Let us say we receive all our data from another function called getProducts
- we would be passing our DB pointer to it.
- now if there is an error, we would send a 500 Response to it, we can use the http.StatusInternalServerError to get the code
- and we would be passing error into a map here, as it will be parsed into json later on

- for a successful from the DB, we call the sendResponse method with 200 status code and our data
*/
func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts(a.DB)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	sendResponse(w, http.StatusOK, products)
}

func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var p product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		sendResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}
	err = p.createProduct(a.DB)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusCreated, p)
}

func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendResponse(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		sendResponse(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	p.ID = id

	if err := p.updateProduct(a.DB); err != nil {
		sendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, "Payload modified")
}

func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendResponse(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	p := product{ID: id}
	if err := p.deleteProduct(a.DB); err != nil {
		sendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) handleRoutes() {
	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/product/{id}", a.updateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id}", a.deleteProduct).Methods("DELETE")
}
