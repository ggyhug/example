# example
## 修改业务逻辑，增加prometheus Exporter
### 具体思路
在原有指标基础上新增3个指标，均为Gauge类型，其特点是Gauge类型的指标侧重于反应系统的当前状态，这类指标的样本数据可增可减。3个指标分别指向cpu当前利用率，cpu当前负载，以及系统memory当前的使用情况，prometheus来监控系统这三项数据并且以Service形式反馈用户。具体参考github.com/shirou/gopsutil/库中各类指标参考定义函数，该库是使用go语言返回系统各类参数的现成库，实验中通过调用定义的函数返回指标对应值。
### 具体代码
定义指标
```
	//系统cpu利用率
	cpu_usage = prometheus.NewGauge(
		prometheus.GaugeOpts{
		Name:      "cpu_usage",
		Help:      "system cpu usage.",
	})
	//系统cpu负载率
	cpu_load = prometheus.NewGauge(
		prometheus.GaugeOpts{
		Name:      "cpu_load",
		Help:      "system cpu load.",
	})
	//系统mem使用情况
	Mem = prometheus.NewGauge(
		prometheus.GaugeOpts{
		Name:      "mem",
		Help:      "system mem usage.",
	})
```
注册指标
```
	prometheus.MustRegister(cpu_usage)
	prometheus.MustRegister(cpu_load)
	prometheus.MustRegister(Mem)
```
调用库函数并且赋值
```
	cpu1,_ :=cpu.Percent(time.Second, false)
	cpu2,_ :=load.Avg()
	mem_,_ :=mem.VirtualMemory()
	cpu_usage.Set(cpu1[0])
	cpu_load.Set(cpu2.Load1)
	Mem.Set(mem_.UsedPercent)
```  
