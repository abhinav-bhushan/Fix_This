package handlers

import (
	"context"
	"encoding/json"
	"golang-project/database"
	"golang-project/models"
	"net/http"
	"time"

	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Laptophandler struct {
	DB *database.Laptop
}

func (l *Laptophandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		glog.Errorln("Method not Implemented")
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Method not implemented"))
	}

	session, err := store.Get(r, "user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve our session values
	username := session.Values["username"]
	password := session.Values["password"]
	if username == nil || password == nil {
		glog.Errorln("Unauthorized access")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Try Login to get Access"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	documents, err := l.DB.GetAll(ctx)
	if err != nil {
		glog.Errorln(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Something went wrong.Please contact admin"))
		return
	}
	w.Header().Set("Content-type", "application/json")
	newJson, err := json.MarshalIndent(documents, "", "")
	if err != nil {
		glog.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong.Please contact admin"))
		return
	}
	w.Write(newJson)

}

func (l *Laptophandler) AddLaptopProblem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		glog.Errorln("Method not Implemented")
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Method not implemented"))
		return
	}

	session, err := store.Get(r, "user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve our session values
	username := session.Values["username"]
	password := session.Values["password"]
	if username == nil || password == nil {
		glog.Errorln("Unauthorized access")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Try Login to get Access"))
		return
	}

	//Encode data to the json
	userproblem := new(models.UserProblem)
	err = json.NewDecoder(r.Body).Decode(userproblem)
	userproblem.UID = session.Values["id"].(string)
	userproblem.Type = "laptop"

	if err != nil {
		glog.Errorln(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Data.Please contact admin"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Check if brand and model is in database
	err=l.DB.Find(ctx,userproblem.Brand,userproblem.Model)
	if err != nil {
		glog.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Brand or Model not in our database"))
		return
	}

	//Insert the response to Problems database
	result, err := l.DB.Insert(ctx, userproblem)
	if err != nil {
		glog.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Problem not in our database"))
		return
	}
	w.Write([]byte(result.(primitive.ObjectID).String()))
}
