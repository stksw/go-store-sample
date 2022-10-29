package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"strconv"
	"strings"
	"context"
	"time"
	"fmt"
	"sort"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func Products(c *fiber.Ctx) error {
	var products []models.Product
	database.DB.Find(&products)

	return c.JSON(products)
}

func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return err
	}

	database.DB.Create(&product)
	return c.JSON(product)
}

func GetProduct(c *fiber.Ctx) error {
	var product models.Product
	id, _ := strconv.Atoi(c.Params("id"))
	product.Id = uint(id)

	database.DB.Find(&product)
	return c.JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	product := models.Product{}
	product.Id = uint(id)

	if err := c.BodyParser(&product); err != nil {
		return err
	}

	database.DB.Model(&product).Updates(&product)
	go deleteCache("products_frontend")
	return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	product := models.Product{}
	product.Id = uint(id)

	database.DB.Delete(&product)
	go deleteCache("products_frontend")
	return nil
}

func deleteCache(key string) {
	time.Sleep(3 * time.Second)
	database.Cache.Del(context.Background(), key)
}

func ProductsFrontend(c *fiber.Ctx) error {
	var products []models.Product
	var ctx = context.Background()

	result, err := database.Cache.Get(ctx, "products_frontend").Result()
	if err != nil {
		// cacheがない場合、DBから取得
		fmt.Println(err.Error())
		database.DB.Find(&products)

		// redisに入れるため、jsonに変換
		bytes, err := json.Marshal(products)
		if err != nil {
			panic(err)
		}

		if errKey := database.Cache.Set(ctx, "products_frontend", bytes, 30*time.Minute).Err(); errKey != nil {
			panic(errKey)
		}
	} else {
		// cacheがある場合、jsonを構造体に戻す
		json.Unmarshal([]byte(result), &products)
	}	

	return c.JSON(products)
}

func SearchProducts(c *fiber.Ctx) error {
	var products []models.Product // DBから取得
	var searchResult []models.Product // 検索結果
	var data []models.Product // 

	database.DB.Find(&products)
	
	if s := c.Query("s"); s != "" {
		lower := strings.ToLower(s)
		for _, product := range products {
			if strings.Contains(strings.ToLower(product.Title), lower) || 
					strings.Contains(strings.ToLower(product.Description), lower) {
				searchResult = append(searchResult, product)
			}
		}
	} else {
		searchResult = products
	}

	if sortParams := c.Query("sort"); sortParams != "" {
		sortLower := strings.ToLower(sortParams)
		if sortLower == "asc" {
			sort.Slice(searchResult, func(i, j int) bool {
				return searchResult[i].Price < searchResult[j].Price
			})
		} else if sortLower == "desc" {
			sort.Slice(searchResult, func(i, j int) bool {
				return searchResult[i].Price > searchResult[j].Price
			})
		}
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("perPage", "8"))

	total := len(searchResult)
	lastPage := map[bool]int{ true: total/perPage, false: total/perPage + 1 }[total % perPage == 0]

	if total <= page * perPage && total >= (page-1) * perPage {
		data = searchResult[(page-1)*perPage : total]
	} else if total >= page * perPage {
		data = searchResult[(page-1) * perPage : page * perPage]
	} else {
		data = []models.Product{}
	}

	return c.JSON(fiber.Map {
		"data": data,
		"total": total,
		"page": page,
		"last_page": lastPage,
	})
}