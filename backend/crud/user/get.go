//package crud is the package that contains all the functions to interact with the database
package CRUD

import (
	"database/sql"
	"encoding/json"
	structures "inprinte/backend/structures"
	utils "inprinte/backend/utils"
	"net/http"

	"github.com/gorilla/mux"
)

//Get returns all the needed data to display the user page
func Get(w http.ResponseWriter, r *http.Request) {
	//set cors headers
	utils.SetCorsHeaders(&w)

	//retrieve url parameters
	vars := mux.Vars(r)
	id_user := vars["id_user"]

	//get the db connection
	db := utils.DbConnect()

	//get informations related to the user
	userData := getUserData(w, db, id_user)
	userFavorite := getUserFavorite(db, id_user)
	userCommandHistory := getCommandHistory(db, id_user)

	//create the json response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(structures.User{
		UserData:         userData,
		FavoriteProducts: userFavorite,
		CommandHistory:   userCommandHistory,
	})
}

//getUserData retrieve the user credentials
func getUserData(w http.ResponseWriter, db *sql.DB, id_user string) structures.UserData {
	//global vars
	var id int
	var firstname, lastname, email, phone string
	var street, city, state, country, zipCode string

	//execute the sql query and check errors
	row := db.QueryRow(`SELECT user.id, first_name, last_name, email, phone, street, city, state, country, zip_code FROM user INNER JOIN address ON user.id = address.id WHERE user.id = ?`, id_user)
	err := row.Scan(&id, &firstname, &lastname, &email, &phone, &street, &city, &state, &country, &zipCode)
	if err == sql.ErrNoRows {
		w.WriteHeader(404)
	} else {
		utils.CheckErr(err)
	}

	//create the json response
	return structures.UserData{
		Id:        id,
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Phone:     phone,
		Address: structures.Address{
			Street:  street,
			City:    city,
			State:   state,
			Country: country,
			ZipCode: zipCode,
		},
	}
}

//getUserFavorite retrieve the user favorite products
func getUserFavorite(db *sql.DB, id_user string) []structures.UserFavoriteProducts {
	//global vars
	var productFavorite []structures.UserFavoriteProducts

	//execute the sql query and check errors
	rows, err := db.Query("SELECT DISTINCT product.id AS id_product, name, price, picture.url FROM product INNER JOIN favorite ON favorite.id_product = product.id INNER JOIN user ON user.id = favorite.id_user INNER JOIN product_picture ON product_picture.id_product = product.id INNER JOIN picture ON picture.id = product_picture.id_picture AND pending_validation = false AND product.is_alive = true WHERE user.id = ?; ", id_user)
	utils.CheckErr(err)

	//parse the query
	for rows.Next() {
		//global vars
		var name, picture string
		var id int
		var price float64

		//retrieve the values and check errors
		err = rows.Scan(&id, &name, &price, &picture)
		utils.CheckErr(err)

		//add the values to the response
		productFavorite = append(productFavorite, structures.UserFavoriteProducts{
			Id:      id,
			Name:    name,
			Price:   price,
			Picture: picture,
		})
	}
	//close the rows

	//create the json response
	return productFavorite
}

//getCommandHistory retrieve the user command history
func getCommandHistory(db *sql.DB, id_user string) []structures.UserCommandHistory {
	//global vars
	var productHistory []structures.UserCommandHistory

	//execute the sql query and check errors
	rows, err := db.Query("SELECT product.id, name, price, picture.url, quantity, command_line.state, command.id FROM product INNER JOIN user ON product.id = user.id INNER JOIN command_line ON product.id = command_line.id INNER JOIN command ON command.id = command_line.id_command INNER JOIN picture ON command_line.id = picture.id WHERE user.id = ?", id_user)
	utils.CheckErr(err)

	//parse the query
	for rows.Next() {
		//global vars
		var name, picture, state string
		var id, quantity, commandNumber int
		var price float64

		//retrieve the values and check errors
		err = rows.Scan(&id, &name, &price, &picture, &quantity, &state, &commandNumber)
		utils.CheckErr(err)

		//add the values to the response
		productHistory = append(productHistory, structures.UserCommandHistory{
			Id:            id,
			Name:          name,
			Price:         price,
			Picture:       picture,
			Quantity:      quantity,
			Status:        state,
			CommandNumber: commandNumber,
		})
	}
	//close the rows

	//create the json response
	return productHistory
}
