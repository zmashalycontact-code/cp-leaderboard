package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/zmashaly/cp-leaderboard/internal/database"
	"github.com/zmashaly/cp-leaderboard/internal/models"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("⚠️ الاستخدام الصحيح:")
		fmt.Println("go run ./cmd/admin/fix_base.go <Handle> <New_Base>")
		return
	}

	handle := os.Args[1]
	newBase, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("❌ الرقم لازم يكون صحيح!")
		return
	}

	db, _ := database.Connect(database.Config{})
	var user models.User
	if err := db.Where("handle = ?", handle).First(&user).Error; err != nil {
		fmt.Println("❌ المتسابق ده مش موجود!")
		return
	}

	oldBase := user.BaseSolvedCount
	user.BaseSolvedCount = newBase
	db.Save(&user)
	
	fmt.Printf("✅ تم تعديل الـ Base للوحش %s بنجاح!\n", handle)
	fmt.Printf("📉 الـ Base القديم: %d\n", oldBase)
	fmt.Printf("📈 الـ Base الجديد: %d\n", newBase)
}
