package logic

import (
	"encoding/csv"
	"math/rand"
	"menu-api/internal/models"
	"os"
	"strconv"
	"time"
)

type DailyMenu struct {
	Date       string       `json:"date"`
	MealOptions []MealOption `json:"meal_options"`
}

type MealOption struct {
	Main              models.Item `json:"main"`
	Side              models.Item `json:"side"`
	Drink             models.Item `json:"drink"`
	TotalCalories     int         `json:"total_calories"`
	CombinedPopularity float64    `json:"combined_popularity"`
}

func loadCSVData() ([]models.Item, error) {
	file, err := os.Open("data/menu.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var items []models.Item
	for i, record := range records {
		if i == 0 { // Skip header
			continue
		}

		calories, _ := strconv.Atoi(record[2])
		popularity, _ := strconv.ParseFloat(record[4], 64)

		items = append(items, models.Item{
			Name:            record[0],
			Category:        record[1],
			Calories:        calories,
			TasteProfile:    record[3],
			PopularityScore: popularity,
		})
	}

	return items, nil
}


func GenerateThreeDayMenu(startDate string) ([]DailyMenu, error) {
	return generateMenu(startDate, 3)
}

func GenerateWeeklyMenu(startDate string) ([]DailyMenu, error) {
	return generateMenu(startDate, 7)
}

func generateMenu(startDate string, days int) ([]DailyMenu, error) {
	rand.Seed(time.Now().UnixNano())

	allItems, err := loadCSVData()
	if err != nil {
		return nil, err
	}

	var mains, sides, drinks []models.Item
	for _, item := range allItems {
		switch item.Category {
		case "main":
			mains = append(mains, item)
		case "side":
			sides = append(sides, item)
		case "drink":
			drinks = append(drinks, item)
		}
	}

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}

	var menu []DailyMenu
	usedCombinations := make(map[string]bool)

	for day := 0; day < days; day++ {
		currentDate := start.AddDate(0, 0, day).Format("2006-01-02")
		var options []MealOption

		for len(options) < 3 {
			main := mains[rand.Intn(len(mains))]
			side := sides[rand.Intn(len(sides))]
			drink := drinks[rand.Intn(len(drinks))]

			comboKey := main.Name + "|" + side.Name + "|" + drink.Name
			if usedCombinations[comboKey] {
				continue
			}

			usedCombinations[comboKey] = true
			options = append(options, MealOption{
				Main:              main,
				Side:              side,
				Drink:             drink,
				TotalCalories:     main.Calories + side.Calories + drink.Calories,
				CombinedPopularity: main.PopularityScore + side.PopularityScore + drink.PopularityScore,
			})
		}

		menu = append(menu, DailyMenu{
			Date:       currentDate,
			MealOptions: options,
		})
	}

	return menu, nil
}