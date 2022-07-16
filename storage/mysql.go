package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/BorisIosifov/mend-home-assignment/object"
)

var mySQLMu sync.Mutex

type MySQL struct {
	db *sqlx.DB
}

func PrepareMySQL() (mySQL MySQL, err error) {
	db, err := sql.Open("mysql", "root:@tcp(mysql:3306)/mend_home_assignment")

	if err != nil {
		return MySQL{}, err
	}

	// defer db.Close()

	dbSqlx := sqlx.NewDb(db, "mysql")

	mySQL = MySQL{
		db: dbSqlx,
	}
	log.Print("MySQL storage is ready")

	return mySQL, nil
}

func (mySQL MySQL) GetList(objectType string) (objects []object.Object, err error) {
	if err != nil {
		return objects, err
	}
	selectQuery := fmt.Sprintf("select * from %s order by id", objectType)

	mySQLMu.Lock()
	defer mySQLMu.Unlock()

	switch objectType {
	case "books":
		var books []object.Book
		err = mySQL.db.Select(&books, selectQuery)
		for _, book := range books {
			obj := book
			objects = append(objects, &obj)
		}
	case "cars":
		var cars []object.Car
		err = mySQL.db.Select(&cars, selectQuery)
		for _, car := range cars {
			obj := car
			objects = append(objects, &obj)
		}
	default:
		return nil, fmt.Errorf("Unexpected objectType %s", objectType)
	}

	return objects, err
}

func (mySQL MySQL) Get(objectType string, ID int) (result object.Object, isNotFound bool, err error) {
	var objects []object.Object

	selectQuery := fmt.Sprintf("select * from %s where id = ?", objectType)

	mySQLMu.Lock()
	defer mySQLMu.Unlock()

	switch objectType {
	case "books":
		var books []object.Book
		err = mySQL.db.Select(&books, selectQuery, ID)
		for _, obj := range books {
			objects = append(objects, &obj)
		}
	case "cars":
		var cars []object.Car
		err = mySQL.db.Select(&cars, selectQuery, ID)
		for _, obj := range cars {
			objects = append(objects, &obj)
		}
	default:
		return nil, false, fmt.Errorf("Unexpected objectType %s", objectType)
	}

	if len(objects) == 0 {
		return nil, true, fmt.Errorf("Object %s with id %d not found", objectType, ID)
	}

	return objects[0], false, err
}

func (mySQL MySQL) Post(objectType string, obj object.Object) (result object.Object, err error) {
	fieldsList, err := getFieldsList(obj)
	if err != nil {
		return nil, err
	}

	fields := strings.Join(fieldsList, ", ")
	values := ":" + strings.Join(fieldsList, ", :")

	mySQLMu.Lock()
	defer mySQLMu.Unlock()

	tx := mySQL.db.MustBegin()
	// insert into books (author, title) values (:author, :title)
	query := fmt.Sprintf("insert into %s (%s) values (%s)", objectType, fields, values)
	queryResult, err := tx.NamedExec(query, obj)
	if err != nil {
		log.Printf("error in NamedExec: %s", err)
		return obj, err
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("error in Commit: %s", err)
		return obj, err
	}

	ID, err := queryResult.LastInsertId()
	if err != nil {
		return obj, err
	}
	obj.SetID(int(ID))
	return obj, err
}

func (mySQL MySQL) Put(objectType string, ID int, obj object.Object) (result object.Object, isNotFound bool, err error) {
	var set string
	fieldsList, err := getFieldsList(obj)
	if err != nil {
		return obj, false, err
	}

	for index, field := range fieldsList {
		if index > 0 {
			set += ", "
		}
		set += fmt.Sprintf("%s = :%s", field, field)
	}
	obj.SetID(ID)

	mySQLMu.Lock()
	defer mySQLMu.Unlock()

	tx := mySQL.db.MustBegin()
	// update books set author = :author, title = :title where id = :id
	query := fmt.Sprintf("update %s set %s where id = :id", objectType, set)
	queryResult, err := tx.NamedExec(query, obj)
	if err != nil {
		return obj, false, err
	}
	err = tx.Commit()
	if err != nil {
		return obj, false, err
	}

	rowsAffected, err := queryResult.RowsAffected()
	if err != nil {
		return obj, false, err
	}

	if rowsAffected == 0 {
		return obj, true, fmt.Errorf("Object %s with id %d not found", objectType, ID)
	}

	return obj, false, err
}

func (mySQL MySQL) Delete(objectType string, ID int) (isNotFound bool, err error) {
	mySQLMu.Lock()
	defer mySQLMu.Unlock()

	tx := mySQL.db.MustBegin()
	query := fmt.Sprintf("delete from %s where id = ?", objectType)
	queryResult := tx.MustExec(query, ID)
	err = tx.Commit()
	if err != nil {
		return false, err
	}

	rowsAffected, err := queryResult.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAffected == 0 {
		return true, fmt.Errorf("Object %s with id %d not found", objectType, ID)
	}

	return false, err
}

func getFieldsList(obj object.Object) (fieldsList []string, err error) {
	var objectMap map[string]interface{}
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return fieldsList, err
	}
	err = json.Unmarshal(jsonData, &objectMap)
	if err != nil {
		return fieldsList, err
	}

	fieldsList = make([]string, 0, len(objectMap))
	for field, _ := range objectMap {
		if field == "ID" {
			continue
		}
		fieldsList = append(fieldsList, strings.ToLower(field))
	}

	return fieldsList, nil
}
