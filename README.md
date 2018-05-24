# open-resolver-datavis
项目用于对于公网的DNS open resolver行为进行探测和分析并提供视图展示


### 数据收集

首先确保已经安装了zmap，zmap是一个高效的网络探测工具

[Install ZMAP](https://github.com/zmap/zmap/blob/master/INSTALL.md)

通过使用zmap来执行全网的开放数据收集，首先对于网络中开放53端口的IP地址进行扫描,其中eno1替换为机器的可以连接Internet的网卡， 带宽大小可以设置为机器自身的带宽容量（设为10Mbps） 模块为udp, 端口53 输出结果到results中

可以使用nohup启动程序并将日志输出到指定文件中

```
zmap -i eno1 -B 10M -M udp -p 53 -o results
```

输出日志如下，可以看到收集过程大概接近3天时间（10Mbps带宽限制下）

```
0:06 0% (2d21h left); send: 89877 14.9 Kp/s (14.8 Kp/s avg); recv: 118 16 p/s (19 p/s avg); drops: 0 p/s (0 p/s avg); hitrate: 0.13%
0:07 0% (2d21h left); send: 104756 14.9 Kp/s (14.8 Kp/s avg); recv: 137 18 p/s (19 p/s avg); drops: 0 p/s (0 p/s avg); hitrate: 0.13%
0:08 0% (2d21h left); send: 119642 14.9 Kp/s (14.8 Kp/s avg); recv: 153 15 p/s (18 p/s avg); drops: 0 p/s (0 p/s avg); hitrate: 0.13%
...

```

results文件将不断的更新采集到的存活的IP地址

```
[root@laptop zmap]## tail results 
217.116.200.69“
51.254.16.0
94.247.133.99
194.28.62.192
209.188.95.140
...

```

### 开放递归筛选



### 数据分析与存储

安装mongodb存储相应的递归数据文件

```
mkdir YOUR_STORE_PATH/mongo/db
docker pull mongo
docker run -it -d -p 27017:27017 -v YOUR_STORE_PATH/mongo/db:/data/db mongo
```

