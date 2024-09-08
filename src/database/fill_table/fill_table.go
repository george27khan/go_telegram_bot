package fill_table

import (
	"context"
	"fmt"
	emp "go_telegram_bot/src/database/employee"
	pstn "go_telegram_bot/src/database/position"
	"os"
	"path/filepath"
	"time"
)

var positionList []string = []string{"стажер", "специалист", "главный специалист"}

func AddPositions(ctx context.Context) error {
	for _, val := range positionList {
		position := pstn.Position{PositionName: val}
		if err := position.Insert(ctx); err != nil {
			return err
		}
	}
	return nil
}

func loadImg(path string) (img []byte) {
	img, err := os.ReadFile(path)
	fmt.Println(err)
	if err != nil {
		return nil
	}
	return
}

func AddEmployees(ctx context.Context) error {
	wd, _ := os.Getwd()
	position, err := pstn.GetByName(ctx, "специалист")
	fmt.Println(os.Getwd())
	if err != nil {
		return err
	}
	fmt.Println(filepath.Join(wd, "src", "database", "fill_table", "image", "1.jpeg"))
	empSlice := []emp.Employee{
		emp.Employee{
			FirstName:   "Иван",
			MiddleName:  "Иванов",
			LastName:    "Иванович",
			BirthDate:   time.Now().AddDate(20, 0, 0),
			Email:       "ivan@gmail.com",
			PhoneNumber: "+71234567892",
			Position:    position,
			HireDate:    time.Now().AddDate(1, 0, 0),
			Photo:       loadImg(filepath.Join(wd, "src", "database", "fill_table", "image", "1.jpeg")),
		},
		emp.Employee{
			FirstName:   "Елена",
			MiddleName:  "Иванова",
			LastName:    "Петровна",
			BirthDate:   time.Now().AddDate(22, 5, 0),
			Email:       "elena@gmail.com",
			PhoneNumber: "+71123334567",
			Position:    position,
			HireDate:    time.Now().AddDate(2, 3, 0),
			Photo:       loadImg(filepath.Join(wd, "src", "database", "fill_table", "image", "2.jpeg")),
		},
		emp.Employee{
			FirstName:   "Александр",
			MiddleName:  "Петров",
			LastName:    "Викторович",
			BirthDate:   time.Now().AddDate(25, 5, 0),
			Email:       "alex@gmail.com",
			PhoneNumber: "+71127774567",
			Position:    position,
			HireDate:    time.Now().AddDate(4, 3, 0),
			Photo:       loadImg(filepath.Join(wd, "src", "database", "fill_table", "image", "3.jpg")),
		},
	}
	for _, emp := range empSlice {
		if err := emp.Insert(ctx); err != nil {
			return err
		}
	}
	return nil
}
