package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/zanjs/y-mugg/app/models"
)

// get all products
func AllProducts(c echo.Context) error {
	var (
		products []models.Product
		err      error
	)
	products, err = models.GetProducts()
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}
	return c.JSON(http.StatusOK, products)
}

// get one product
func ShowProduct(c echo.Context) error {
	var (
		product models.Product
		err     error
	)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	product, err = models.GetProductById(id)
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}
	return c.JSON(http.StatusOK, product)
}

// CreateProduct is product
func CreateProduct(c echo.Context) error {

	product := new(models.Product)

	product.Title = c.FormValue("title")
	product.ExternalCode = c.FormValue("external_code")
	sortV := c.FormValue("sort")
	sort, _ := strconv.Atoi(sortV)
	fmt.Println(sort)

	product.Sort = sort

	err := models.CreateProduct(product)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	return c.JSON(http.StatusCreated, product)
}

// UpdateProduct is update product
func UpdateProduct(c echo.Context) error {
	// Parse the content
	product := new(models.Product)

	product.Title = c.FormValue("title")
	product.ExternalCode = c.FormValue("external_code")

	sortV := c.FormValue("sort")
	sort, _ := strconv.Atoi(sortV)

	product.Sort = sort

	// get the param id
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	m, err := models.GetProductById(id)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	// update product data
	err = m.UpdateProduct(product)
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	return c.JSON(http.StatusOK, m)
}

//delete product
func DeleteProduct(c echo.Context) error {
	var err error

	// get the param id
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	m, err := models.GetProductById(id)
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	err = m.DeleteProduct()
	return c.JSON(http.StatusNoContent, err)
}
