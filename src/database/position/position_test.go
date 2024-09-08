package position

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	db "go_telegram_bot/src/database"
	"os"
	"path"
	"runtime"
	"testing"
)

func TestPosition_Insert(t *testing.T) {
	db.Pool(context.Background())
	testPosition := Position{
		Id:           0,
		PositionName: "Test_position",
	}
	err := testPosition.Insert(context.Background())
	fmt.Println("1111111111111111111111111111111111111 ", err)
	assert.Nil(t, err)
}

func init() {
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println("11111111111111111111111111 filename ", filename)
	// The ".." may change depending on you folder structure
	dir := path.Join(path.Dir(filename), "../..")
	fmt.Println("11111111111111111111111111 dir ", dir)

	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
