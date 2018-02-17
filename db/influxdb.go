package db

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/influxdata/influxdb/client/v2"
	"jvole.com/monitor/util"
)

type Influxdb interface {
	/*
	   写influx数据库
	   @parm tags 标签相当于属性
	   @parm fields 存储的字段集合，key value
	   @parm precision  精度 h,m,s,ms,ns
	   @parm table 表明
	   @error
	*/
	WriteInflux(tags map[string]string, fields map[string]interface{}, table string) error
}

type influxdb struct {
	addr        string //连接地址
	user        string //用户名
	passwd      string //密码
	client      client.Client
	buff        int                //当前缓存数量
	batchPoints client.BatchPoints //数据行
}

var (
	Buffer = 10 //缓存，达到一定数量后写数据库，可提高效率。默认值：1000
)

/*
   写influx数据库
   @parm tags 标签相当于属性
   @parm fields 存储的字段集合，key value
   @parm precision  精度 h,m,s,ms,ns
   @parm table 表明
   @error
*/
func (idb *influxdb) WriteInflux(tags map[string]string, fields map[string]interface{}, table string) error {
	pt, err := client.NewPoint(table, tags, fields, time.Now())
	if err != nil {
		util.Log.WithFields(logrus.Fields{
			"name": "错误",
			"err":  err,
		}).Errorln("出错了")
	}
	idb.batchPoints.AddPoint(pt)
	idb.buff++
	if idb.buff >= Buffer { //缓存满了写数据库
		// pp.Println(len(idb.batchPoints.Points()))
		//写数据库
		if err := idb.client.Write(idb.batchPoints); err != nil {
			util.Log.WithFields(logrus.Fields{
				"name":     "错误",
				"lenPoint": len(idb.batchPoints.Points()),
				"err":      err,
			}).Errorln("出错了")
		}
		idb.buff = 0
		idb.batchPoints.ClearPoint() //自定义方法，每次清除point
	}
	return nil
}

func NewInfluxdb(addr, user, password, db, precision string) *influxdb {
	// fmt.Printf("addr:%s,user:%s,pwd:%s\n", addr, user, password)

	// 创建 point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db,
		Precision: precision,
	})
	if err != nil {
		util.Log.WithFields(logrus.Fields{
			"name": "错误",
			"err":  err,
		}).Errorln("出错了")
	}

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: user,
		Password: password,
	})
	if err != nil {
		util.Log.WithFields(logrus.Fields{
			"name": "错误",
			"err":  err,
		}).Errorln("出错了")
	}
	indb := &influxdb{addr, user, password, c, 0, bp}
	return indb
}
