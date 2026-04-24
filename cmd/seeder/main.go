package main

import (
	"fmt"
	"github.com/zmashaly/cp-leaderboard/internal/database"
	"github.com/zmashaly/cp-leaderboard/internal/models"
)

func main() {
	db, err := database.Connect(database.Config{})
	if err != nil {
		panic("مش قادر أتصل بقاعدة البيانات: " + err.Error())
	}

	// لازم نضيف العمود الجديد في الداتا بيز
	database.AutoMigrate(db, &models.User{})
	db.Exec("TRUNCATE users CASCADE;")

	usersToSeed := []models.User{
		{Handle: "TryHardZ", DisplayName: "زياد مشالي", AtCoderHandle: "TryHardZ", BaseSolvedCount: 173},
		{Handle: "AhmedEmadl", DisplayName: "أحمد عماد", AtCoderHandle: "AhmedEmad", BaseSolvedCount: 1012},
		{Handle: "Ahmed_Gawish", DisplayName: "أحمد جاويش", AtCoderHandle: "", BaseSolvedCount: 330},
		{Handle: "Armo38", DisplayName: "محمد عرفة", AtCoderHandle: "", BaseSolvedCount: 36},
		{Handle: "Mariam--Gamal", DisplayName: "مريم جمال", AtCoderHandle: "", BaseSolvedCount: 390},
		{Handle: "Aurora_10", DisplayName: "منة شتات", AtCoderHandle: "mennashatat", BaseSolvedCount: 328},
		{Handle: "emankassab2006", DisplayName: "إيمان كساب", AtCoderHandle: "emankassab2006", BaseSolvedCount: 312},
		{Handle: "S8ti", DisplayName: "محمد جاويش", AtCoderHandle: "", BaseSolvedCount: 602},
		{Handle: "RoqayaWaleed", DisplayName: "رقية وليد", AtCoderHandle: "RoqayaWaleed", BaseSolvedCount: 359},
		{Handle: "mohamedmagdy4920", DisplayName: "محمد مجدي", AtCoderHandle: "mohamedmagdy279", BaseSolvedCount: 609},
		{Handle: "Khaled_Zalama", DisplayName: "خالد عادل", AtCoderHandle: "Khaled_Zalama", BaseSolvedCount: 577},
		{Handle: "TryAgain0", DisplayName: "عبدالله الخولي", AtCoderHandle: "", BaseSolvedCount: 18},
		{Handle: "Maryam_zz", DisplayName: "مريم السعيد", AtCoderHandle: "", BaseSolvedCount: 152},
		{Handle: "RunTimeTerror00", DisplayName: "يوسف مجدي", AtCoderHandle: "RunTimeTerror00", BaseSolvedCount: 64},
		{Handle: "yousef1234556", DisplayName: "يوسف محمد", AtCoderHandle: "", BaseSolvedCount: 45},
	}

	for _, u := range usersToSeed {
		u.TotalSolved = u.BaseSolvedCount
		db.Create(&u)
		fmt.Printf("✅ تم إضافة: %-15s | %s\n", u.Handle, u.DisplayName)
	}
}
