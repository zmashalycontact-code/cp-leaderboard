package main

import (
	"fmt"
	"log"

	"github.com/zmashaly/cp-leaderboard/internal/database"
	"github.com/zmashaly/cp-leaderboard/internal/models"
)

func main() {
	db, err := database.Connect(database.Config{})
	if err != nil {
		log.Fatal("❌ فشل الاتصال بقاعدة البيانات")
	}

	var users []models.User
	db.Find(&users)

	for _, u := range users {

		u.BaseSolvedCount = u.TotalSolved
		

		u.ManualBonus = 0
		u.SeasonPoints = 0
		u.CFPoints = 0
		u.AtCoderPoints = 0
		u.HiddenSolved = 0
		u.Activity7D = 0
		u.StruggleCount = 0
		u.PeakWeeklyRating = 0
		

		
		db.Save(&u)
	}

	fmt.Println("✅ تم تصفير جميع المتسابقين وبدء سباق الشهر الجديد بنجاح!")
}
