package main

import (
	"fmt"
	"os"

	"github.com/zmashaly/cp-leaderboard/internal/database"
	"github.com/zmashaly/cp-leaderboard/internal/models"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("⚠️ اكتب اسم المتسابق بعد الأمر. مثال: go run cmd/check/main.go TryHardZ")
		return
	}

	db, _ := database.Connect(database.Config{})
	var u models.User
	if err := db.Where("handle = ?", os.Args[1]).First(&u).Error; err != nil {
		fmt.Println("❌ المتسابق ده مش موجود!")
		return
	}

	fmt.Printf("\n🕵️‍♂️ فحص شامل للوحش: %s 🕵️‍♂️\n", u.Handle)
	fmt.Printf("=====================================\n")
	fmt.Printf("1️⃣ التوتال الحالي في البروفايل: %d\n", u.TotalSolved)
	fmt.Printf("2️⃣ التوتال اللي بدأ بيه السيزون (Base): %d\n", u.BaseSolvedCount)
	fmt.Printf("3️⃣ الفرق بينهم (المفروض يكونوا مسائل جديدة): %d\n", u.TotalSolved-u.BaseSolvedCount)
	fmt.Printf("4️⃣ البونص اليدوي اللي إنت اديتهوله (Admin Tool): %d\n", u.ManualBonus)
	fmt.Printf("5️⃣ إجمالي المسائل المخفية المحسوبة حالياً (الشيتات): %d\n", u.HiddenSolved)
	fmt.Printf("6️⃣ إجمالي النقاط (CF Points): %.2f\n", u.CFPoints)
	fmt.Printf("=====================================\n\n")
}
