package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/zmashaly/cp-leaderboard/internal/database"
	"github.com/zmashaly/cp-leaderboard/internal/models"
)

func main() {
	handle := flag.String("handle", "", "Codeforces Handle (مطلوب)")
	// غيرنا النوع لـ Float64 عشان يقبل 0.5
	add := flag.Float64("add", 0.0, "قيمة البونص بالنقاط (مثلاً 0.5, 1.5, 2)")

	flag.Parse()

	if *handle == "" || *add == 0 {
		log.Fatal("❌ لازم تكتب الـ handle وتحدد قيمة البونص، مثال: -handle=User -add=0.5")
	}

	db, _ := database.Connect(database.Config{})

	var user models.User
	err := db.Where("handle = ?", *handle).First(&user).Error

	if err != nil {
		log.Fatalf("❌ المتسابق %s مش موجود في الداتابيز!", *handle)
	}

	// تزويد أو خصم البونص من النقاط الصريحة
	user.ManualBonus += *add
	db.Save(&user)

	fmt.Printf("🎁 تم إضافة %.1f نقطة بونص للمتسابق %s بنجاح! (إجمالي البونص الحالي بتاعه: %.1f نقطة)\n", *add, *handle, user.ManualBonus)
}
