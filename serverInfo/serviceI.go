package serverInfo

type RecordInfo interface {
	/*
	   记录cpu信息到influxdb
	   @return error
	*/
	CpuInfo() error
	/*
	   记录内存信息到influxdb
	   @return error
	*/
	MemInfo() error
	/*
	   记录磁盘信息到influxdb
	   @return error
	*/
	DiskInfo() error
	/*
	   记录负载信息到influxdb
	   @return error
	*/
	LoadInfo() error
	/*
	   记录网络信息到influxdb
	   @return error
	*/
	NetInfo() error
	/*
	   记录进程信息到influxdb
	   @return error
	*/
	ProcessInfo() error
}
