# go_sensor

### Installing

```Shell
go mod init sensorproject
sudo apt-get install libpcap-dev
go get github.com/google/gopacket/pcap
go get github.com/gookit/config/v2
sudo groupadd pcap
sudo usermod -a -G pcap $USER
sudo chgrp pcap /usr/sbin/tcpdump
sudo chmod 750 /usr/sbin/tcpdump
```