package tests

import (
	"fmt"
	"strconv"
	"time"

	"github.com/tofa-project/client-daemon/db"
	"github.com/tofa-project/client-daemon/glob"
)

func DBCruds() {
	rowID := db.MakeApp(glob.J{"some": "DATA"})

	fmt.Println(db.GetAppByID(rowID))

	db.UpdateApp(rowID, glob.J{"meme": "love you" + strconv.Itoa(time.Now().Nanosecond())})

	fmt.Println(db.GetAppByID(rowID))
}
