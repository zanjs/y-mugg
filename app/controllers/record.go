package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/zanjs/y-mugg/app/models"
	"github.com/zanjs/y-mugg/app/services"
	"github.com/zanjs/y-mugg/app/utils"
)

// AllProductWareroom get all products and wareroom
func AllProductWareroom(c echo.Context) error {
	var (
		data models.ProductWareroom
		err  error
	)
	data.Products, err = models.GetProducts()
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	data.Warerooms, err = models.GetWarerooms()
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	go func() {

		for _, v := range data.Warerooms {
			var wareroom models.Wareroom
			wareroom = v
			numbering := wareroom.Numbering
			fmt.Println(numbering)

			for _, v2 := range data.Products {

				var product models.Product
				product = v2
				fmt.Println(product)

				var qmProduct models.QMProduct
				qmProduct.ItemCode = product.ExternalCode
				qmProduct.WarehouseCode = wareroom.Numbering
				qmProduct.OwnerCode = "bkyy"
				qmProduct.InventoryType = "ZP"

				var qmRequest models.QMRequest
				qmRequest = utils.Parameter("inventory.query", qmProduct)

				var record models.Record
				record.ProductID = product.ID
				record.WareroomID = wareroom.ID

				services.QMHTTPPost(qmRequest, record)

			}

		}

	}()

	return c.JSON(http.StatusOK, data)
}

// AllProductWareroomRecords get all products and wareroom
func AllProductWareroomRecords(c echo.Context) error {
	var (
		data     models.ProductWareroomExcel
		products []models.Product
		err      error
	)
	products, err = models.GetProducts()
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	data.Warerooms, err = models.GetWarerooms()
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	productExcelInt := []models.ProductExcel{}

	productExcel := productExcelInt[0:]

	for _, v := range products {
		var product models.Product
		product = v

		pID := product.ID

		productExcelQuantitysInt := []models.ProductExcelQuantity{}

		productExcelQuantitys := productExcelQuantitysInt[0:]

		for _, v2 := range data.Warerooms {
			fmt.Println(v2)
			wID := v2.ID
			var record models.Record

			record, err = models.GetRecordLast(wID, pID)

			if err != nil {
				fmt.Println("查询最好一个 err: ", err, record)
			}

			var productExcelQuantity models.ProductExcelQuantity

			productExcelQuantity.Quantity = record.Quantity
			productExcelQuantity.Sales = record.Sales

			productExcelQuantitys = append(productExcelQuantitys, productExcelQuantity)
		}

		var pExcle models.ProductExcel
		pExcle.ProductTitle = product.Title
		pExcle.ProductExcelQuantitys = productExcelQuantitys
		productExcel = append(productExcel, pExcle)
	}

	fmt.Println("productExcel:", productExcel)

	data.Products = productExcel
	return c.JSON(http.StatusOK, data)
}

// AllProductWareroomRecordsTime get all products and wareroom
func AllProductWareroomRecordsTime(c echo.Context) error {
	var (
		data        models.ProductWareroomExcel
		products    []models.Product
		queryparams models.QueryParams
		err         error
	)

	qps := c.QueryParams()

	startTimeq := c.QueryParam("start_time")
	endTime := c.QueryParam("end_time")

	fmt.Println(qps)
	fmt.Println(startTimeq)
	fmt.Println(endTime)

	queryparams.StartTime = startTimeq
	queryparams.EndTime = endTime

	products, err = models.GetProducts()
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	data.Warerooms, err = models.GetWarerooms()
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	productExcelInt := []models.ProductExcel{}

	productExcel := productExcelInt[0:]

	for index := 0; index < len(products); index++ {
		var product models.Product
		product = products[index]

		pID := product.ID

		productExcelQuantitysInt := []models.ProductExcelQuantity{}

		productExcelQuantitys := productExcelQuantitysInt[0:]

		for _, v2 := range data.Warerooms {
			fmt.Println(v2)
			wID := v2.ID
			var records []models.Record
			var record models.Record

			var qyps models.QueryParams

			qyps.StartTime = queryparams.StartTime
			qyps.EndTime = queryparams.EndTime
			qyps.ProductID = pID
			qyps.WareroomID = wID

			record, err = models.GetRecordTimeLast(qyps)

			records, err = models.GetRecordWhereTime(qyps)

			if err != nil {
				fmt.Println("查询 where time err: ", err, records)
			}

			sales := 0
			for _, v3 := range records {
				fmt.Println("v3: ", v3)

				sales = sales + v3.Sales
			}

			record.Sales = sales

			var productExcelQuantity models.ProductExcelQuantity

			productExcelQuantity.Quantity = record.Quantity
			productExcelQuantity.Sales = record.Sales
			productExcelQuantity.CreatedAt = record.CreatedAt

			productExcelQuantitys = append(productExcelQuantitys, productExcelQuantity)
		}

		var pExcle models.ProductExcel
		pExcle.ProductTitle = product.Title
		pExcle.ProductExcelQuantitys = productExcelQuantitys
		productExcel = append(productExcel, pExcle)
	}

	fmt.Println("productExcel:", productExcel)

	data.Products = productExcel
	return c.JSON(http.StatusOK, data)
}

// GetRecordWhereTime is
func GetRecordWhereTime(c echo.Context) error {
	var (
		data        []models.Record
		queryparams models.QueryParams
		err         error
	)

	qps := c.QueryParams()

	limitq := c.QueryParam("limit")
	offsetq := c.QueryParam("offset")
	startTimeq := c.QueryParam("start_time")
	endTime := c.QueryParam("end_time")
	pIDq := c.QueryParam("product_id")
	wIDq := c.QueryParam("wareroom_id")

	limit, _ := strconv.Atoi(limitq)
	offset, _ := strconv.Atoi(offsetq)
	pID, _ := strconv.Atoi(pIDq)
	wID, _ := strconv.Atoi(wIDq)
	fmt.Println(qps)
	fmt.Println(limit)
	fmt.Println(offset)

	if limit == 0 {
		limit = 10
	}

	queryparams.Limit = limit
	queryparams.Offset = offset
	queryparams.StartTime = startTimeq
	queryparams.EndTime = endTime
	queryparams.ProductID = pID
	queryparams.WareroomID = wID

	data, err = models.GetRecordWhereTime(queryparams)
	// data, err = models.GetRecordsAll()

	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}
	return c.JSON(http.StatusOK, data)
}

// AllRecordsPage  get all records
func AllRecordsPage(c echo.Context) error {
	var (
		data        models.RecordPage
		queryparams models.QueryParams
		err         error
	)

	qps := c.QueryParams()

	limitq := c.QueryParam("limit")
	offsetq := c.QueryParam("offset")
	startTimeq := c.QueryParam("start_time")
	endTime := c.QueryParam("end_time")

	limit, _ := strconv.Atoi(limitq)
	offset, _ := strconv.Atoi(offsetq)
	fmt.Println(qps)
	fmt.Println(limit)
	fmt.Println(offset)

	if limit == 0 {
		limit = 10
	}

	queryparams.Limit = limit
	queryparams.Offset = offset
	queryparams.StartTime = startTimeq
	queryparams.EndTime = endTime

	data, err = models.GetRecords(queryparams)
	// data, err = models.GetRecordsAll()

	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}
	return c.JSON(http.StatusOK, data)
}

// AllRecords  get all records
func AllRecords(c echo.Context) error {

	var (
		records []models.Record
		err     error
	)
	records, err = models.GetRecordsAll()
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}
	return c.JSON(http.StatusOK, records)
}

// ShowRecord get one record
func ShowRecord(c echo.Context) error {
	var (
		record models.Record
		err    error
	)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	record, err = models.GetRecordById(id)
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}
	return c.JSON(http.StatusOK, record)
}

//DeleteRecord is record
func DeleteRecord(c echo.Context) error {
	var err error

	// get the param id
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	m, err := models.GetRecordById(id)
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	err = m.DeleteRecord()
	return c.JSON(http.StatusNoContent, err)
}

// UpdateRecord is update record
func UpdateRecord(c echo.Context) error {
	// Parse the content
	record := new(models.Record)

	quantity, _ := strconv.Atoi(c.FormValue("quantity"))
	sales, _ := strconv.Atoi(c.FormValue("sales"))

	record.Quantity = quantity
	record.Sales = sales

	// get the param id
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	m, err := models.GetRecordById(id)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	// update record data
	err = m.UpdateRecord(record)
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	return c.JSON(http.StatusOK, m)
}
