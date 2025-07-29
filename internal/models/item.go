package models

type Item struct {
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	Calories        int     `json:"calories"`
	TasteProfile    string  `json:"taste_profile"`
	PopularityScore float64 `json:"popularity_score"`
}