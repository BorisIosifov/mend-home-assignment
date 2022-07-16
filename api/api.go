package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/schema"

	"github.com/BorisIosifov/mend-home-assignment/object"
)

type Storage interface {
	GetList(objectType string) (objects []object.Object, err error)
	Get(objectType string, ID int) (result object.Object, isNotFound bool, err error)
	Post(objectType string, obj object.Object) (result object.Object, err error)
	Put(objectType string, ID int, obj object.Object) (result object.Object, isNotFound bool, err error)
	Delete(objectType string, ID int) (isNotFound bool, err error)
}

type API struct {
	Storage Storage
}

func (api API) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var (
		objectType string
		obj        object.Object
		ID         int
		err        error
	)
	w.Header().Set("Content-Type", "text/json")

	pathSlice := strings.Split(req.URL.Path, "/")
	if len(pathSlice) <= 1 {
		w.WriteHeader(http.StatusNotFound)
		api.errorResult(w, fmt.Errorf("Page not found"))
		return
	}

	switch pathSlice[1] {
	case "books":
		objectType = "books"
		obj = &object.Book{}
	case "cars":
		objectType = "cars"
		obj = &object.Car{}
	default:
		w.WriteHeader(http.StatusNotFound)
		api.errorResult(w, fmt.Errorf("Page not found"))
		return
	}

	if len(pathSlice) == 2 || (len(pathSlice) == 3 && pathSlice[2] == "") {
		// /books || /books/
		ID = 0
	} else if len(pathSlice) == 3 || (len(pathSlice) == 4 && pathSlice[3] == "") {
		// /books/123 || /books/123/
		ID, err = strconv.Atoi(pathSlice[2])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			api.errorResult(w, fmt.Errorf("Bad request"))
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		api.errorResult(w, fmt.Errorf("Bad request"))
		return
	}

	switch req.Method {
	case "GET":
		api.get(w, req, objectType, ID)
	case "POST":
		api.post(w, req, objectType, obj)
	case "PUT":
		api.put(w, req, objectType, ID, obj)
	case "DELETE":
		api.delete(w, req, objectType, ID)
	default:
		w.WriteHeader(http.StatusBadRequest)
		api.errorResult(w, fmt.Errorf("Bad request"))
		return
	}
}

func (api API) get(w http.ResponseWriter, req *http.Request, objectType string, ID int) {
	// fmt.Fprintf(w, "GET request\n")
	if ID == 0 {
		objects, err := api.Storage.GetList(objectType)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			api.errorResult(w, fmt.Errorf("Internal server error: %s", err))
			return
		}

		resJSON, err := json.Marshal(objects)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			api.errorResult(w, fmt.Errorf("Internal server error: %s", err))
			return
		}
		fmt.Fprintln(w, string(resJSON))

	} else {
		obj, isNotFound, err := api.Storage.Get(objectType, ID)

		if isNotFound {
			w.WriteHeader(http.StatusNotFound)
			api.errorResult(w, fmt.Errorf("%s", err))
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			api.errorResult(w, fmt.Errorf("Internal server error: %s", err))
			return
		}

		resJSON, err := json.Marshal(obj)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			api.errorResult(w, fmt.Errorf("Internal server error: %s", err))
			return
		}
		fmt.Fprintln(w, string(resJSON))
	}
}

func (api API) post(w http.ResponseWriter, req *http.Request, objectType string, obj object.Object) {
	var err error
	// fmt.Fprintf(w, "POST request\n")
	err = fillObjectFromForm(req, obj)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		api.errorResult(w, fmt.Errorf("Internal server error: %s", err))
		return
	}

	obj, err = api.Storage.Post(objectType, obj)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		api.errorResult(w, fmt.Errorf("Internal server error: %s", err))
		return
	}

	resJSON, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		api.errorResult(w, fmt.Errorf("Internal server error: %s", err))
		return
	}
	fmt.Fprintln(w, string(resJSON))
}

func (api API) put(w http.ResponseWriter, req *http.Request, objectType string, ID int, obj object.Object) {
	var err error
	// fmt.Fprintf(w, "PUT request\n")
	err = fillObjectFromForm(req, obj)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		api.errorResult(w, fmt.Errorf("Internal server error: %s", err))
		return
	}

	obj, isNotFound, err := api.Storage.Put(objectType, ID, obj)

	if isNotFound {
		w.WriteHeader(http.StatusNotFound)
		api.errorResult(w, fmt.Errorf("%s", err))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		api.errorResult(w, fmt.Errorf("Internal server error: %s", err))
		return
	}

	resJSON, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		api.errorResult(w, fmt.Errorf("Internal server error: %s", err))
		return
	}
	fmt.Fprintln(w, string(resJSON))
}

func (api API) delete(w http.ResponseWriter, req *http.Request, objectType string, ID int) {
	var err error
	// fmt.Fprintf(w, "DELETE request\n")
	isNotFound, err := api.Storage.Delete(objectType, ID)

	if isNotFound {
		w.WriteHeader(http.StatusNotFound)
		api.errorResult(w, fmt.Errorf("%s", err))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		api.errorResult(w, fmt.Errorf("Internal server error: %s", err))
		return
	}

	fmt.Fprintln(w, "{}")
}

type errorResult struct {
	Err string `json:"error"`
}

func (api API) errorResult(w http.ResponseWriter, err error) {
	res := errorResult{err.Error()}
	resJSON, _ := json.Marshal(res)
	fmt.Fprintln(w, string(resJSON))
}

func fillObjectFromForm(req *http.Request, obj object.Object) (err error) {
	err = req.ParseForm()
	if err != nil {
		return err
	}

	decoder := schema.NewDecoder()
	err = decoder.Decode(obj, req.PostForm)
	return err
}
