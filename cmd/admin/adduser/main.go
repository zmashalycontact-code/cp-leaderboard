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
	name := flag.String("name", "", "Display Name (مطلوب)")
	atcoder := flag.String("atcoder", "", "AtCoder Handle (اختياري)")
	baseCount := flag.Int("base", 0, "Base Solved Count as of April 14 (مهم جداً)")

	flag.Parse()

	if *handle == "" || *name == "" {
		log.Fatal("❌ خطأ: لازم تدخل الـ Handle والـ Name.\nمثال: go run cmd/admin/adduser/main.go -handle=tourist -name=\"Ziad Mashaly\" -base=485")
	}

	fmt.Printf("⏳ جاري إضافة المتسابق: %s (@%s) بأساس مسائل: %d...\n", *name, *handle, *baseCount)

	db, err := database.Connect(database.Config{})
	if err != nil {
		log.Fatalf("❌ فشل الاتصال بقاعدة البيانات: %v", err)
	}

	userRepo := repository.NewUserRepository(db)

	newUser := &models.User{
		Handle:          *handle,
		AtCoderHandle:   *atcoder,
		DisplayName:     *name,
		BaseSolvedCount: *baseCount,
	}

	err = userRepo.Create(context.Background(), newUser)
	if err != nil {
		log.Fatalf("❌ فشل حفظ المتسابق: %v", err)
	}

	fmt.Println("✅ تم إضافة المتسابق بنجاح لليدربورد!")
}
