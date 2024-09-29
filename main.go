package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type shoppingCart struct {
	ID       string  `json:"id"`
	Item     string  `json:"item"`
	Price    float32 `json:"price"`
	Quantity int     `json:"quantity"`
}

var cart = []shoppingCart{
	{ID: "1", Item: "Banana", Price: 4.99, Quantity: 1},
	{ID: "2", Item: "Apple", Price: 10, Quantity: 1},
	{ID: "3", Item: "Cheese", Price: 0.99, Quantity: 1},
}

func getCart(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, cart)
}

func addCart(context *gin.Context) {
	var newCart shoppingCart

	if err := context.BindJSON(&newCart); err != nil {
		return
	}

	cart = append(cart, newCart)

	context.IndentedJSON(http.StatusCreated, newCart)
}

func getItem(context *gin.Context) {
	id := context.Param("id")
	item, err := getItemById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Item not found"})
		return
	}

	context.IndentedJSON(http.StatusOK, item)
}

func updateQuantity(context *gin.Context) {
	id := context.Param("id")
	item, err := getItemById(id)
	var updatedItem struct {
		Quantity int `json:"quantity"`
	}

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Item not found"})
		return
	}

	if err := context.BindJSON(&updatedItem); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	i, err := strconv.Atoi(id)

	item.Quantity = updatedItem.Quantity
	context.JSON(http.StatusOK, cart[i])
	return
}

func deleteItem(context *gin.Context) {
	id := context.Param("id")
	_, err := getItemById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Item not found"})
		return
	}

	i, err := strconv.Atoi(id)

	cart = append(cart[:i-1], cart[i:]...)
	context.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
	return
}

func getItemById(id string) (*shoppingCart, error) {
	for i, t := range cart {
		if t.ID == id {
			return &cart[i], nil
		}
	}

	return nil, errors.New("Item not found")
}

func main() {
	router := gin.Default()
	router.GET("/cart", getCart)
	router.GET("/cart/:id", getItem)
	router.POST("/cart", addCart)
	router.PUT("/cart/:id", updateQuantity)
	router.DELETE("/cart/:id", deleteItem)
	router.Run("localhost:9090")
}
