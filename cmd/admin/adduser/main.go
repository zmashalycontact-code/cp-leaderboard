package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/zmashaly/cp-leaderboard/internal/database"
	"github.com/zmashaly/cp-leaderboard/internal/models"
	"github.com/zmashaly/cp-leaderboard/internal/repository"
)

func main() {
	handle := flag.String("handle", "", "Codeforces Handle (مطلوب)")
	name := flag.String("name", "", "Display Name (مطلوب عند الإضافة)")
	cheat := flag.Bool("cheat", false, "Mark as cheater (true/false)")
	base := flag.Int("base", 0, "Base solved count")

	flag.Parse()

	if *handle == "" {
		log.Fatal("❌ لازم تكتب الـ handle على الأقل")
	}

	db, _ := database.Connect(database.Config{})
	userRepo := repository.NewUserRepository(db)
	ctx := context.Background()


	var user models.User
	err := db.Where("handle = ?", *handle).First(&user).Error

	if err == nil {

		user.IsCheater = *cheat
		db.Save(&user)
		fmt.Printf("✅ تم تحديث حالة المتسابق %s. غشاش: %v\n", *handle, *cheat)
	} else {

		newUser := &models.User{
			Handle:          *handle,
			DisplayName:     *name,
			BaseSolvedCount: *base,
			IsCheater:       *cheat,
		}
		userRepo.Create(ctx, newUser)
		fmt.Printf("✅ تم إضافة متسابق جديد وحالته غشاش: %v\n", *cheat)
	}
}
