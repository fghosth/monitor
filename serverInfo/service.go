package serverInfo

import (
	"encoding/json"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/net"
	indb "jvole.com/monitor/db"
	"jvole.com/monitor/util"
)

type RecordServer struct {
	CpuInterval time.Duration
	Influxdb    indb.Influxdb
	Db          string //数据库名
	Precision   string //精度 h,m,s,ms,ns
	ServerName  string //服务器名称
	Path        string //要监控的路径，只支持 一个
}

var (
	Interval  = 1 //记录数据的时间间隔
	cupTable  = "cpu"
	diskTable = "disk"
	loadTable = "load"
	netTable  = "net"
	breakgo   = false //退出进程型号
)

func NewRecordServer(addr, user, password, db, precision, servername, path string) RecordServer {
	rs := &RecordServer{}
	rs.Db = db
	rs.Precision = precision
	rs.Influxdb = indb.NewInfluxdb(addr, user, password, db, precision)
	rs.ServerName = servername
	rs.Path = path
	return *rs
}

func (rd RecordServer) Stop() {
	breakgo = true
}
func (rd RecordServer) Back() {
	breakgo = false
}
func (rd RecordServer) Run() {
	breakgo = false
	go func() {
		for range time.Tick(time.Duration(Interval) * time.Second) {
			if breakgo {
				break
			}
			err := rd.CpuInfo()
			if err != nil {
				util.Log.WithFields(logrus.Fields{
					"name": "错误",
					"err":  err,
				}).Errorln("出错了")
			}
			err = rd.DiskInfo()
			if err != nil {
				util.Log.WithFields(logrus.Fields{
					"name": "错误",
					"err":  err,
				}).Errorln("出错了")
			}

			err = rd.LoadInfo()
			if err != nil {
				util.Log.WithFields(logrus.Fields{
					"name": "错误",
					"err":  err,
				}).Errorln("出错了")
			}
			err = rd.NetInfo()
			if err != nil {
				util.Log.WithFields(logrus.Fields{
					"name": "错误",
					"err":  err,
				}).Errorln("出错了")
			}

		}
		// pp.Println("======Out")
	}()
}

/*
   记录cpu信息到Influxdb
   @return error
*/
func (rd RecordServer) CpuInfo() error {
	per, err := cpu.Percent(time.Second*1, false) //第二个参数是多核cpu是否分开显示
	if err != nil {
		return err
	}
	tags := map[string]string{"server": rd.ServerName, "type": "cpu_per"}
	fields := map[string]interface{}{
		"percent": per[0],
	}
	return rd.Influxdb.WriteInflux(tags, fields, cupTable)
}

func (rd RecordServer) MemInfo() error {

	return nil
}
func (rd RecordServer) DiskInfo() error {
	usageStat, err := disk.Usage(rd.Path)
	if err != nil {
		return err
	}
	ustring := usageStat.String()
	var dat map[string]interface{}
	json.Unmarshal([]byte(ustring), &dat)
	tags := map[string]string{
		"server": rd.ServerName,
		"type":   "disk_usageStat",
		"path":   rd.Path,
	}
	total := dat["total"].(float64) / 1024 / 1024 / 1024
	free := dat["free"].(float64) / 1024 / 1024 / 1024
	used := dat["used"].(float64) / 1024 / 1024 / 1024
	fields := map[string]interface{}{
		"total": total,
		"free":  free,
		"used":  used,
	}
	// pp.Println(dat["total"].(float64) >> 10)
	// pp.Println(dat["total"].(float64) / 1024 / 1024 / 1024)
	return rd.Influxdb.WriteInflux(tags, fields, diskTable)
}
func (rd RecordServer) LoadInfo() error {
	load, err := load.Avg()
	if err != nil {
		return err
	}
	tags := map[string]string{"server": rd.ServerName, "type": "load_avg"}
	// pp.Println(tags)
	fields := map[string]interface{}{
		"load1":  load.Load1,
		"load5":  load.Load5,
		"load15": load.Load15,
	}
	return rd.Influxdb.WriteInflux(tags, fields, loadTable)
}
func (rd RecordServer) NetInfo() error {

	nets, err := net.IOCounters(false)
	if err != nil {
		return err
	}
	netjson := nets[0].String()
	var dat map[string]interface{}
	json.Unmarshal([]byte(netjson), &dat)
	tags := map[string]string{"server": rd.ServerName, "type": "load_avg"}
	fields := dat
	return rd.Influxdb.WriteInflux(tags, fields, netTable)
}
func (rd RecordServer) ProcessInfo() error {
	return nil
}
