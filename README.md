# tftp-proxy-server
TFTP server for proxying from tftp to http/https 

cidrfile.yaml is a white list for tftp access.  This causes the tftp server to not distribute content if the IP is not in the cidr.

#Build
```
go get github.com/pin/tftp
go get gopkg.in/yaml.v2
go build
```
#Run
As root
```
./tftp-proxy-server
```

