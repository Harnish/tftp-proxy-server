package main

import (
	"fmt"
	"github.com/pin/tftp"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

type conf struct {
	Port    int    `yaml:"port"`
	UrlBase string `yaml:"urlbase"`
	IpFile  string `yaml:"ipfile"`
	Cidrs   []string
}

var c conf

func main() {
	c.getConf()
	s := tftp.NewServer(readHandler, nil)
	s.SetTimeout(5 * time.Second)                       // optional
	err := s.ListenAndServe(":" + strconv.Itoa(c.Port)) // blocks until s.Shutdown() is called
	if err != nil {
		fmt.Fprintf(os.Stdout, "server: %v\n", err)
		os.Exit(1)
	}
}

func readHandler(filename string, rf io.ReaderFrom) error {
	raddr := rf.(tftp.OutgoingTransfer).RemoteAddr()
	laddr := rf.(tftp.RequestPacketInfo).LocalIP()
	log.Println("RRQ from", raddr.String(), "To ", laddr.String(), " For ", filename)
	if !CheckIP(raddr.String(), c.Cidrs) {
		log.Println("Connection Rejected.  Not in List")
		return nil
	}
	resp, err := http.Get(c.UrlBase + filename)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	n, err := rf.ReadFrom(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	fmt.Printf("%d bytes sent\n", n)
	return nil
}

func (c *conf) getConf() *conf {

	yamlFile, err := ioutil.ReadFile("tftp-proxy-server.yaml")
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	yamlFile2, err := ioutil.ReadFile(c.IpFile)
	if err != nil {
		fmt.Printf("No Cidr file.  All will be blocked")
	}
	var cidrs []string
	err = yaml.Unmarshal(yamlFile2, &cidrs)
	if err != nil {
		fmt.Printf("Nothing in the file")
	}
	c.Cidrs = cidrs
	fmt.Println(c)
	return c
}

func CheckIP(ip string, cidrs []string) bool {

	for _, cidr := range cidrs {
		_, cidrnet, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err) // assuming I did it right above
		}
		myaddr := net.ParseIP(ip)
		if cidrnet.Contains(myaddr) {
			return true
		}
	}
	return false
}
