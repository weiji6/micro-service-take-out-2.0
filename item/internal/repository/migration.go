package repository

import "fmt"

func migration() {
	err := DB.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&Item{
				ItemID:   1,
				ItemName: "book",
				Price:    1.0,
				Stock:    100,
			},
		)
	if err != nil {
		fmt.Println("migration err:", err)
	}

	var count int64
	DB.Model(&Item{}).Count(&count)
	if count == 0 {
		initialItem := &Item{
			ItemID:   1,
			ItemName: "book",
			Price:    1.0,
			Stock:    100,
		}
		if err := DB.Create(initialItem).Error; err != nil {
			fmt.Println("initial item err:", err)
		}
	}
}
