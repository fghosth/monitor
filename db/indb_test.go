package db_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"jvole.com/monitor/db"
)

var dbclinet = db.NewInfluxdb("http://localhost:8086", "derek", "123456", "serverInfo", "ns")

func TestWrite(t *testing.T) {
	table := "cpu"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 10; i++ {
		tags := map[string]string{"server": "test", "type": "cpu_per"}
		fields := map[string]interface{}{
			"percent": r.Float64(),
		}
		err := dbclinet.WriteInflux(tags, fields, table)
		if err != nil {
			fmt.Println(err)
		}
	}

	// defer dbclinet.Close()
	// pp.Println(result[0])
}
