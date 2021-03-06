//package crud is the package that contains all the functions to interact with the database
package crud

import (
	"database/sql"
	"encoding/json"
	structures "inprinte/backend/structures"
	utils "inprinte/backend/utils"
	"strconv"

	"net/http"
)

//Get returns all the informations for the boutique page
func Get(w http.ResponseWriter, r *http.Request) {
	//set the cors headers
	utils.SetCorsHeaders(&w)

	//connect to the database
	db := utils.DbConnect()

	//get all the informations
	newProducts := getNewProducts(db)
	mostSales := getMostSales(db)
	allProducts := getAllProducts(r, db)
	categories := getCategories(db)

	//add the values to the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(structures.Boutique{
		NewProducts: newProducts,
		MostSales:   mostSales,
		AllProducts: allProducts,
		Categories:  categories,
	})
}

//getNewProducts returns the three last products added
func getNewProducts(db *sql.DB) []structures.NewProduct {
	//global vars
	var newProducts []structures.NewProduct

	//execute the sql query and check errors
	rows, err := db.Query("SELECT product.id, name, price FROM product WHERE product.pending_validation = false AND product.is_alive = true ORDER BY product.id DESC LIMIT 3")
	utils.CheckErr(err)

	//parse the query
	for rows.Next() {
		//global vars
		var name string
		var price float64
		var id int
		//retrieve the values and check errors
		err = rows.Scan(&id, &name, &price)
		utils.CheckErr(err)
		product_picture := getProductPicture(db, strconv.Itoa(id))

		//add the values to the response
		newProducts = append(newProducts, structures.NewProduct{
			Id_product: id,
			Name:       name,
			Price:      price,
			Picture:    product_picture[0],
		})
	}
	//close the rows

	//create the json response
	return newProducts
}

func getMostSales(db *sql.DB) []structures.MostSales {
	//global vars
	var mostSales []structures.MostSales
	var product_picture []string

	//execute the sql query and check errors
	rows, err := db.Query("SELECT COUNT(command_line.id_product) AS nbrOrder, command_line.id_product, product.name, product.price FROM command_line INNER JOIN product ON command_line.id_product = product.id WHERE pending_validation = false AND product.is_alive = true GROUP BY command_line.id_product ORDER BY nbrOrder DESC LIMIT 3")
	utils.CheckErr(err)

	//parse the query
	for rows.Next() {
		//global vars
		var name string
		var id, nbrOrder int
		var price float64

		//retrieve the values and check errors
		err = rows.Scan(&nbrOrder, &id, &name, &price)
		utils.CheckErr(err)

		product_picture = getProductPicture(db, strconv.Itoa(id))
		//add the values to the response
		mostSales = append(mostSales, structures.MostSales{
			Id_product: id,
			Name:       name,
			Price:      price,
			Picture:    product_picture[0],
		})
	}

	//create the json response
	return mostSales
}

//getAllProducts returns all the products
func getAllProducts(r *http.Request, db *sql.DB) []structures.BoutiqueProduct {
	//global vars
	var allProducts []structures.BoutiqueProduct
	var product_picture []string

	//execute the sql query and check errors
	rows, err := db.Query("SELECT product.id, name, price, description FROM product WHERE product.pending_validation = false AND product.is_alive = true")
	utils.CheckErr(err)

	//parse the query
	for rows.Next() {
		//global vars
		var name, description string
		var price float64
		var id int

		//retrieve the values and check errors
		err = rows.Scan(&id, &name, &price, &description)
		utils.CheckErr(err)

		product_picture = getProductPicture(db, strconv.Itoa(id))

		//add the values to the response
		allProducts = append(allProducts, structures.BoutiqueProduct{
			Id_product:  id,
			Name:        name,
			Price:       price,
			Picture:     product_picture[0],
			Description: description,
		})
	}
	//close the rows

	//create the json response
	return allProducts
}

//getCategories returns all the categories
func getCategories(db *sql.DB) []string {
	//global vars
	var allCategories []string

	//execute the sql query and check errors
	rows, err := db.Query("SELECT name FROM `category`")
	utils.CheckErr(err)

	//parse the query
	for rows.Next() {
		//global vars
		var name string

		//retrieve the values and check errors
		err = rows.Scan(&name)
		utils.CheckErr(err)

		//add the values to the response
		allCategories = append(allCategories, name)
	}
	//close the rows

	//create the json response
	return allCategories
}

//getProductPicture returns all the pictures for a product
func getProductPicture(db *sql.DB, id_product string) []string {
	//global vars
	var product_picture []string
	var url string
	//create the sql query
	sqlQuery := ("SELECT picture.url FROM picture INNER JOIN product_picture ON product_picture.id_picture = picture.id INNER JOIN product ON product.id = product_picture.id_product WHERE product.id = " + id_product + ";")

	rows, err := db.Query(sqlQuery)
	utils.CheckErr(err)

	for rows.Next() {
		//retrieve the values and check errors
		err = rows.Scan(&url)
		utils.CheckErr(err)

		//add the values to the response
		product_picture = append(product_picture, url)
	}
	//close the rows

	//create the json response
	return product_picture
}
