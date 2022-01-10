package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getproductID(c *gin.Context) {
	db, err := sql.Open("postgres", "postgres://postgres:qwerty123@localhost:5432/api")
	if err != nil {
		fmt.Println("could not connect to database: ", err)
	}

	id := c.Param("id")
	fmt.Println(id)
	res := []Product{}
	fmt.Println(res)
	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		fmt.Println("error")
	}
	for rows.Next() {
		emp := Product{}

		err = rows.Scan(&emp.Id, &emp.Name, &emp.Price, &emp.Tax, &emp.Seller_id)
		if err != nil {
			fmt.Println("scan error", err)
		}
		if id == emp.Id {
			c.IndentedJSON(http.StatusOK, emp)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
}

//-----------------------------------------------seller create product----------------------------------//

func createproduct(c *gin.Context) {
	db, err := sql.Open("postgres", "postgres://postgres:qwerty123@localhost:5432/api")
	if err != nil {
		fmt.Println("could not connect to database: ", err)
	}
	email := c.PostForm("email")
	name := c.PostForm("name")
	price := c.PostForm("price")
	tax := c.PostForm("tax")
	seller_id := c.Param("seller")
	fmt.Println(seller_id)

	emp := Seller{}
	err = db.QueryRow("SELECT * FROM seller WHERE email='"+email+"'").Scan(&emp.Id, &emp.Name, &emp.Email, &emp.Phoneno, &emp.Role)
	fmt.Println(emp.Role)
	switch {
	case emp.Role == "2":

		rows, err := db.Query("INSERT INTO products( name , price, tax,seller_id)	VALUES ('" + name + "', '" + price + "','" + tax + "','" + seller_id + "') ")
		if err != nil {
			fmt.Println("inserted successfully", rows)
		} else {
			fmt.Println("error")
		}
		c.IndentedJSON(http.StatusOK, gin.H{
			"email": email,
			"name":  name,
			"price": price,
			"tax":   tax,
		})

		return
	case err != nil:
		c.IndentedJSON(http.StatusOK, gin.H{
			"Message": "Email id not regestered",
		})
		return
	default:
		c.IndentedJSON(http.StatusOK, "empty")
	}
}

//----------------------------------------------orderproduct-----------------------------------------//

func orderproduct(c *gin.Context) {

	db, err := sql.Open("postgres", "postgres://postgres:qwerty123@localhost:5432/api")
	if err != nil {
		fmt.Println("could not connect to database: ", err)
	}
	email := c.PostForm("email")
	quantity := c.PostForm("quantity")

	fmt.Println(email)
	emp := Buyer{}
	res := Product{}
	err = db.QueryRow("SELECT * from buyer where email='"+email+"'").Scan(&emp.Id, &emp.Name, &emp.Email, &emp.Phoneno, &emp.Role)
	if err != nil {
		fmt.Println("selected", err)
	}

	switch {
	case emp.Role == "1":
		err = db.QueryRow("SELECT * from buyer, products Where email='"+email+"'").Scan(&emp.Id, &emp.Name, &emp.Email, &emp.Phoneno, &emp.Role, &res.Id, &res.Name, &res.Price, &res.Tax, &res.Seller_id)
		if err != nil {
			fmt.Println("inserted successfully")
		} else {
			fmt.Println("error", err)
		}
		quan, err := strconv.ParseFloat(quantity, 64)
		fmt.Println(quan, err)
		total_price := (quan * (res.Price))
		fmt.Println(total_price)
		total_tax := (quan) * res.Tax
		fmt.Println(total_tax)
		total := total_price + total_tax
		fmt.Println(total)
		c.IndentedJSON(http.StatusOK, gin.H{
			"Message":   "You can order your product ",
			"email":     emp.Email,
			"Name":      res.Name,
			"Tax":       res.Tax,
			"Price":     res.Price,
			"total":     total,
			"seller_id": res.Seller_id,
		})
		return
	case err != nil:
		c.IndentedJSON(http.StatusOK, gin.H{
			"Message": "email not regester ",
		})
		return

	default:
		c.IndentedJSON(http.StatusOK, "404")
	}
}
