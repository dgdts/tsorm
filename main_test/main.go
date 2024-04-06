package main

import (
	"fmt"
	"tsorm"

	_ "github.com/mattn/go-sqlite3"
)

// Item represents an item entity in the database.
type Item struct {
	Name string // Name represents the name of the item.
	Id   uint32 // Id represents the unique identifier of the item.
}

func main() {
	// Open a new engine instance to interact with the database.
	engine, _ := tsorm.NewEngine("sqlite3", "ts.db")
	defer engine.Close() // Close the engine connection when main function exits.

	// Create a new session to perform database operations.
	s := engine.NewSession()
	s.Model(&Item{})                                             // Set the model for the session to operate on the Item entity.
	s.DropTable()                                                // Drop the table if it exists.
	s.CreateTable()                                              // Create the table based on the Item schema.
	s.Insert(&Item{Name: "a1", Id: 1}, &Item{Name: "a2", Id: 2}) // Insert sample data into the table.

	// Retrieve items from the database and store them in the items slice.
	var items []Item
	s.Find(&items)

	// Print the retrieved items.
	for _, item := range items {
		fmt.Println(item.Id, item.Name)
	}
}
