package main

import (
	"database/sql"
	"fmt"
	"log"
)

/*
now, before I create the methods here,
we know to process data returned as rows in the sql package,
we require a struct,

so let us create a struct here,
we are going to add json tags over here as well.
Now these json tags would be helpful whenever we would encode this data into json format.

Which we are obviously doing in our code.
*/
type product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func (p *product) getProduct(db *sql.DB) error {

	return db.QueryRow(fmt.Sprintf("SELECT name, quantity, price FROM products WHERE id=%v",
		p.ID)).Scan(&p.Name, &p.Quantity, &p.Price)

}

func (p *product) updateProduct(db *sql.DB) error {
	_, err :=
		db.Exec(fmt.Sprintf("UPDATE products SET name='%v', quantity=%v, price=%v WHERE id=%v",
			p.Name, p.Quantity, p.Price, p.ID))

	return err
}

func (p *product) deleteProduct(db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf("DELETE FROM products WHERE id=%v", p.ID))

	return err
}

func (p *product) createProduct(db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO products(name, quantity, price) VALUES('%v', %v, %v)",
		p.Name, p.Quantity, p.Price)
	res, err := db.Exec(query)

	if err != nil {
		return err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("The last inserted row id: %d\n", lastId)
	p.ID = int(lastId)
	return nil
}

/*
Our getProducts method here

- would receive a DB pointer as the input

and would return us a

- slice of product type, which would be holding our row data on successful query fetch
- and also an error
*/
func getProducts(db *sql.DB) ([]product, error) {

	// let us form our select query here
	query := fmt.Sprintf("SELECT id, name,  quantity, price FROM products")

	// and receive the rows
	rows, err := db.Query(query)

	// if there is any error, we are going to send it back
	if err != nil {
		return nil, err
	}

	/*
		if not, let us create an empty slice called products of type product struct.
		we will simply keep on looping our rows and appending the data to our slice
	*/
	products := []product{}

	for rows.Next() {
		var p product
		if err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	// let us return the processed slice now
	return products, nil
}
