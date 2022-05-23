package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Customer struct {
	ID        string
	FirstName string
	LastName  string
	Username  string
	Phones    []Phone
	Adresses  []Adress
	Product   []Products
	Email     string
	Gender    string
	Birthdate time.Time
	Password  string // should be hashed and validate password should be 8 symbols
	Status    string
	Numbers   []int64
}

type Phone struct {
	Id      string
	Numbers []int64
	Code    string
}

type Adress struct {
	ID         string
	Country    string
	City       string
	District   string
	Postalcode int64
}

type Products struct {
	ID          string
	Name        string
	Types       []Type
	Cost        int64
	OrderNumber int64
	Amount      int64
	Currency    string
	Rating      int64
}

type Type struct {
	ID   string
	Name string
}

func insert() {

	connStr := "user=postgres password=1234 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Unable to connect!")
		panic(err)
	}
	defer db.Close()

	user := Customer{}

	tx, err := db.Begin()
	if err != nil {
		signal := fmt.Sprintf("%v\n", err)
		fmt.Println(signal)
	}

	user = Customer{
		FirstName: `Sultankhodja`,
		LastName:  `Jorabekov`,
		Username:  `s_jurabekoff`,
		Phones: []Phone{
			{
				Numbers: []int64{998418295, 996074503},
				Code:    `+998`,
			},
		},
		Adresses: []Adress{
			{
				Country:    `Uzbekistan`,
				City:       `Toshkent`,
				District:   `Chirchiq`,
				Postalcode: 12345,
			},
		},
		Product: []Products{Products{
			Name: `Olma`,
			Types: []Type{
				{
					Name: `Meva`,
				},
			},
			Cost:        20500,
			OrderNumber: 12,
			Amount:      20,
			Currency:    `Sum`,
			Rating:      9,
		},
		},
		Email:     `sjurabekov1@gmail.com`,
		Gender:    `Male`,
		Birthdate: time.Date(1995, time.November, 28, 0, 0, 0, 0, time.Local),
		Password:  `123123`,
		Status:    `active`,
		Numbers:   []int64{1233, 456465},
	}

	customerId, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
	}

	queryCustomerInsert := (`INSERT INTO customers (id, first_name, last_name, username, email, gender, birthday, password, status, numbers) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`)
	err = db.QueryRow(queryCustomerInsert, customerId, user.FirstName, user.LastName, user.Username, user.Email, user.Gender, user.Birthdate, user.Password, user.Status, pq.Array(user.Numbers)).Scan(&customerId)
	if err != nil {
		fmt.Println("Error while inserting customer data!")
		tx.Rollback()
		panic(err)
	}

	for _, prod := range user.Product {
		productId, err := uuid.NewRandom()
		if err != nil {
			fmt.Println(err)
		}

		queryProductInsert := (`INSERT INTO products (id, customer_id, name, cost, ordernumber, amount, currency, rating) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`)
		err = db.QueryRow(queryProductInsert, productId, customerId, prod.Name, prod.Cost, prod.OrderNumber, prod.Amount, prod.Currency, prod.Rating).Scan(&productId)
		if err != nil {
			fmt.Println("Error while inserting products data!")
			tx.Rollback()
			panic(err)
		}

		for _, typ := range prod.Types {
			typeId, err := uuid.NewRandom()
			if err != nil {
				fmt.Println(err)
			}

			queryTypeInsert := (`INSERT INTO types (id, product_id, name) VALUES($1, $2, $3) RETURNING id`)
			err = db.QueryRow(queryTypeInsert, typeId, productId, typ.Name).Scan(&typeId)
			if err != nil {
				fmt.Println("Error while inserting type data!")
				tx.Rollback()
				panic(err)
			}
		}
	}

	for _, ad := range user.Adresses {
		adressId, err := uuid.NewRandom()
		if err != nil {
		}

		queryAdressInsert := (`INSERT INTO adress (id, customer_id, country, city, district, postalcodes) VALUES($1, $2, $3, $4, $5, $6) RETURNING id`)
		err = db.QueryRow(queryAdressInsert, adressId, customerId, ad.Country, ad.City, ad.District, ad.Postalcode).Scan(&adressId)
		fmt.Println(err)
		if err != nil {
			fmt.Println("Error while inserting adresses data!")
			tx.Rollback()
			panic(err)
		}
	}

	for _, ph := range user.Phones {

		phoneId, err := uuid.NewRandom()
		if err != nil {
			fmt.Println(err)
		}

		queryPhoneInsert := (`INSERT INTO phones (id, customer_id, numbers, code) VALUES($1, $2, $3, $4) RETURNING id`)
		err = db.QueryRow(queryPhoneInsert, phoneId, customerId, pq.Array(ph.Numbers), ph.Code).Scan(&phoneId)
		if err != nil {
			fmt.Println("Error while inserting phones data!")
			tx.Rollback()
			panic(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

func get(customerId string) (Customer, error) {

	var (
		cust       Customer
		err        error
		adres      []Adress
		prod       []Products
		type_vales []Type
		phone      []Phone
	)
	connStr := "user=postgres password=1234 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Unable to connect!")
		panic(err)
	}
	defer db.Close()

	customerQuery := `Select id, first_name, last_name, username, email, gender, birthday, password, status, numbers from customers where id = $1`
	row := db.QueryRow(customerQuery, customerId)

	if err != nil {
		fmt.Println("Error while selecting customer data!")
		panic(err)
	}
	err = row.Scan(
		&cust.ID,
		&cust.FirstName,
		&cust.LastName,
		&cust.Username,
		&cust.Email,
		&cust.Gender,
		&cust.Birthdate,
		&cust.Password,
		&cust.Status,
		pq.Array(&cust.Numbers),
	)

	if err != nil {
		fmt.Println("Error while selecting customer data!")
		panic(err)
	}

	adressQuery := `select id, country, city, district, postalcodes from adress where customer_id = $1`
	rows, err := db.Query(adressQuery, customerId)

	for rows.Next() {
		var adr Adress
		err := rows.Scan(
			&adr.ID,
			&adr.Country,
			&adr.City,
			&adr.District,
			&adr.Postalcode,
		)
		if err != nil {
			panic(err)
			fmt.Println("Error while selecting adress data!")
		}
		adres = append(adres, adr)
	}

	phoneQuery := `select id, numbers, code from phones where customer_id = $1`
	rows, err = db.Query(phoneQuery, customerId)

	for rows.Next() {
		var ph Phone
		err := rows.Scan(
			&ph.Id,
			pq.Array(&ph.Numbers),
			&ph.Code,
		)
		if err != nil {
			panic(err)
			fmt.Println("Error while selecting phone data!")
		}
		phone = append(phone, ph)
	}

	productQuery := `select id, name, cost, ordernumber, amount, currency, rating from products where customer_id = $1`
	rows, err = db.Query(productQuery, customerId)

	for rows.Next() {
		var pd Products
		err := rows.Scan(
			&pd.ID,
			&pd.Name,
			&pd.Cost,
			&pd.OrderNumber,
			&pd.Amount,
			&pd.Currency,
			&pd.Rating,
		)
		if err != nil {
			fmt.Println("Error while selecting product data!")
			panic(err)
		}

		typeQuery := `select id, name from types`
		rows, err := db.Query(typeQuery)
		for rows.Next() {
			var tp Type
			err := rows.Scan(
				&tp.ID,
				&tp.Name,
			)
			if err != nil {
				fmt.Println("Error while selecting type data!")
				panic(err)
			}
			type_vales = append(type_vales, tp)
		}
		pd.Types = type_vales
		prod = append(prod, pd)
	}
	defer rows.Close()
	cust.Adresses = adres
	cust.Phones = phone
	cust.Product = prod
	return cust, err
}
func main() {
	customer, err := get("d6a0e8de-ee5e-47ed-b304-d11c1ccd8c2d")
	if err != nil {
		fmt.Println(err)
	}
	info := customer
	info1, _ := json.MarshalIndent(info, ",", "  ")

	fmt.Println(string(info1))
	//insert()
}
