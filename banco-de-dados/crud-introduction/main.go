package main

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "docker"
	dbname   = "go-mba"
)

type Product struct {
	ID    string
	Name  string
	Price float64
}

type ProductDTO struct {
	name  string
	price float64
}

func NewProduct(product ProductDTO) *Product {
	return &Product{
		ID:    uuid.New().String(),
		Name:  product.name,
		Price: product.price,
	}
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	data := ProductDTO{
		name:  "Notebook",
		price: 1899.99,
	}

	product := NewProduct(data)

	err = insert(db, product)

	if err != nil {
		panic(err)
	}

	product.Price = 1499.99

	err = save(db, product)

	if err != nil {
		panic(err)
	}

	p, err := findById(db, product.ID)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Product: %v, possui o preço de %.2f", p.Name, p.Price)

	products, err := findMany(db)

	if err != nil {
		panic(err)
	}

	for _, p := range products {
		fmt.Printf("Product: %v, possui o preço de %.2f\n", p.Name, p.Price)
	}

	err = remove(db, product.ID)

	if err != nil {
		panic(err)
	}
}

func findById(db *sql.DB, id string) (*Product, error) {
	stmt, err := db.Prepare("select id, name, price from products where id = ?")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var product Product

	err = stmt.QueryRow(id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
	)

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func findMany(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("select id, name, price from products")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []Product

	for rows.Next() {
		var product Product

		err = rows.Scan(&product.ID, &product.Name, &product.Price)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func insert(db *sql.DB, product *Product) error {
	stmt, err := db.Prepare("insert into products (id, name, price) values (?, ?, ?)")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(product.ID, product.Name, product.Price)

	if err != nil {
		return err
	}

	return nil
}

func save(db *sql.DB, product *Product) error {
	stmt, err := db.Prepare("update products set name = ?, price = ? where id = ?")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.Price, product.ID)

	if err != nil {
		return err
	}

	return nil
}

func remove(db *sql.DB, id string) error {
	stmt, err := db.Prepare("delete from products where id = ?")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		return err
	}

	return nil
}
