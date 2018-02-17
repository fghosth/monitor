package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli"
	"jvole.com/monitor/db"
	"jvole.com/monitor/serverInfo"
	"jvole.com/monitor/util"
)

var (
	ERRADDR     = errors.New("数据库地址不能为空")
	ERRUSERNAME = errors.New("用户名不能为空")
	ERRPASSWORD = errors.New("密码不能为空")
	ERRSERVER   = errors.New("服务器名称不能为空")
	ERRPATH     = errors.New("要监控的磁盘路径不能为空")
	ERRNOCONFIG = errors.New("请指定配置文件")
	indb        = "serverInfo"
	precision   = "s"
	buffer      = 10
	cfg         *config
)

type config struct {
	serverName string //服务器名称
	addr       string //infulx地址
	username   string //influx用户名
	password   string //influx密码
	path       string //监控的路径
	dbname     string //influx数据库名
	precision  string //精确度
	buffer     int    //缓存
}

func loadconfig(file string) {
	util.Viper.SetConfigType("yaml")
	// util.Viper.SetConfigName(".cfg")
	util.Viper.AddConfigPath(".")
	// util.Viper.AddConfigPath("/Users/derek/project/go/src/jvole.com/monitor/")
	util.Viper.SetConfigFile(file)

	err := util.Viper.ReadInConfig() // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		util.Log.WithFields(logrus.Fields{
			"name": "信息",
			"err":  err,
		}).Infoln("配置文件加载失败")
		return
	}

	cfg = &config{}
	mapConfig()
}
func mapConfig() {
	if influxcfg, ok := util.Viper.Get("inflxudb").(map[string]interface{}); ok {
		cfg.addr = influxcfg["addr"].(string)
		cfg.dbname = influxcfg["dbname"].(string)
		cfg.precision = influxcfg["precision"].(string)
		cfg.username = influxcfg["username"].(string)
		cfg.password = influxcfg["password"].(string)
	}
	if v, ok := util.Viper.Get("buffer").(int); ok {
		cfg.buffer = v
	}
	if v, ok := util.Viper.Get("serverName").(string); ok {
		cfg.serverName = v
	}
	if v, ok := util.Viper.Get("path").(string); ok {
		cfg.path = v
	}
}
func main() {
	app := cli.NewApp()
	app.Name = "RecordeServerInfo"
	app.Version = "0.1.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "derek",
			Email: "fghosth@163.com",
		},
	}
	app.Copyright = "(c) derek fan"
	app.Usage = "服务器监控等"
	app.UsageText = "用于服务器监控等"
	app.Commands = []cli.Command{
		{
			Name:    "RecordeServerInfo",
			Aliases: []string{"start"},
			Usage:   "启动服务",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, f",
					Value: "",
					Usage: "配置文件地址 -f /etc/server.yaml",
				},
				// cli.StringFlag{
				// 	Name:  "address, addr",
				// 	Value: "",
				// 	Usage: "数据库地址 -addr http://jvole.com:8086",
				// },
				// cli.StringFlag{
				// 	Name:  "username, u",
				// 	Value: "",
				// 	Usage: "用户名 -u derek",
				// },
				// cli.StringFlag{
				// 	Name:  "password, p",
				// 	Value: "",
				// 	Usage: "密码 -p 333333",
				// },
				// cli.StringFlag{
				// 	Name:  "serverName, s",
				// 	Value: "",
				// 	Usage: "服务器名称 -s newbidder",
				// },
				// cli.StringFlag{
				// 	Name:  "path",
				// 	Value: "",
				// 	Usage: "要监控的磁盘路径 -path /data",
				// },
			},
			Action: func(c *cli.Context) error {

				err := recordS(c)
				if err != nil {
					util.Log.WithFields(logrus.Fields{
						"name": "错误",
						"err":  err,
					}).Errorln("出错了")
					os.Exit(1)
				}
				return err
			},
		},
	}
	app.Run(os.Args)
	fmt.Println("运行中....")
	ch := make(chan int)
	<-ch // 阻塞main goroutine, 信道c被锁
}

func recordS(c *cli.Context) error {
	f := c.String("config")
	if f == "" {
		return ERRNOCONFIG
	}
	loadconfig(f)
	if cfg == nil {
		addr := c.String("address")
		u := c.String("username")
		p := c.String("password")
		s := c.String("serverName")

		path := c.String("path")
		if addr == "" {
			return ERRADDR
		}
		if u == "" {
			return ERRUSERNAME
		}
		if p == "" {
			return ERRPASSWORD
		}
		if s == "" {
			return ERRSERVER
		}
		if path == "" {
			return ERRPATH
		}
		cfg = &config{}
		cfg.addr = addr
		cfg.username = u
		cfg.password = p
		cfg.serverName = s
		cfg.path = path
		cfg.dbname = indb
		cfg.precision = precision
		cfg.buffer = buffer
	}

	// pp.Println(cfg)
	db.Buffer = cfg.buffer
	rs := serverInfo.NewRecordServer(cfg.addr, cfg.username, cfg.password, cfg.dbname, cfg.precision, cfg.serverName, cfg.path)
	rs.Run()

	//监控配置文件变化
	util.Viper.WatchConfig()
	util.Viper.OnConfigChange(func(e fsnotify.Event) {
		util.Log.WithFields(logrus.Fields{
			"name": "信息",
		}).Infoln("配置文件变更，重新生效")
		rs.Stop()
		mapConfig()
		db.Buffer = cfg.buffer
		// pp.Println(cfg)

		time.Sleep(time.Duration(serverInfo.Interval) * time.Second) //保证进程退出
		rs2 := serverInfo.NewRecordServer(cfg.addr, cfg.username, cfg.password, cfg.dbname, cfg.precision, cfg.serverName, cfg.path)
		rs2.Run()
	})
	return nil
}
