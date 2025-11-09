package routes

import (
	"fmt"
	"log"
	"meal-choices/db/schema"
	"meal-choices/db/tables"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func HandleRecipeAdd(c echo.Context) error {

	c.Request().ParseForm()
	name := c.Request().Form.Get("name")
	book := c.Request().Form.Get("book")
	pageNum, _ := strconv.Atoi(c.Request().Form.Get("page"))

	if name == "" || book == "" {
		var missingValue []string
		if book == "" {
			missingValue = append(missingValue, "book")
		}
		if name == "" {
			missingValue = append(missingValue, "name")
		}

		c.Render(422, "recipe-error", fmt.Sprintf("Problem creating recipe, missing: %v", strings.Join(missingValue, ", ")))
		return c.Render(422, "recipe-form", &schema.Recipe{Id: -1, Name: name, Book: book, Page: pageNum})
	}

	_, err := tables.AddRecipe(name, book, pageNum)
	if err != nil {
		message := "Problem creating recipe"
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			message = fmt.Sprintf("Recipe: \"%v\" in Book: \"%v\" already exists.", name, book)
		}

		c.Render(422, "recipe-error", message)
		return c.Render(422, "recipe-form", &schema.Recipe{Id: -1, Name: name, Book: book, Page: pageNum})
	}

	c.Render(200, "recipe-result", fmt.Sprintf("Recipe \"%v\" added.", name))
	return c.Render(422, "recipe-form", &schema.Recipe{})
}

func HandleGetAllRecipes(c echo.Context) error {
	recipes, err := tables.GetAllRecipes()
	if err != nil {
		log.Fatal(err)
		return c.Render(500, "recipes-list", nil)
	}

	return c.Render(200, "recipes-list", recipes)
}
