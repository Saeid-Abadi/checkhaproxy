package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"slices"
	"strings"
)

type Server struct {
	PXname string
	Name   string
	IP     string
	Port   string
	Status string
	Weight string
}

func main() {
	servers := fetchservers()
	fmt.Println(sortservers(servers))
}

func fetchservers() []Server {
	myserver := []Server{}

	// Connect to HAProxy admin.socket
	conn, err := net.Dial("unix", "/var/run/haproxy/admin.sock")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to HAProxy socket:", err)
		return nil
	}
	defer conn.Close()

	// Send the "show stat" command to retrieve server stats
	fmt.Fprintln(conn, "show stat")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Split(line, ",")
		if len(fields) < 34 {
			continue
		}

		//fmt.Println(fields)
		//if fields[1] != "BACKEND" && fields[1] != "FRONTEND" {
		// Extracting relevant fields from the CSV output
		pxname := fields[0]
		server := fields[1]
		status := fields[17]
		//lastChk := fields[18]
		//checkHealth := fields[39]
		//downtime := fields[13]
		weight := fields[6]
		addr := fields[73]
		ipport := strings.Split(addr, ":")
		ip := ""
		port := ""
		if len(ipport) == 2 {
			ip = ipport[0]
			port = ipport[1]
		}
		//currentSessions := fields[4]
		//maxSessions := fields[5]
		newserver := Server{PXname: pxname, Name: server, IP: ip, Port: port, Status: status, Weight: weight}
		myserver = append(myserver, newserver)

		//}

	}
	//for _, srv := range myserver {
	//    return srv
	//    fmt.Println(srv)
	//    }
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from HAProxy socket:", err)
		return nil
	}
	return myserver
}
func sortservers(server []Server) ([]Server, []Server, []Server) {
	frontends := []Server{}
	backends := []Server{}
	servers := []Server{}

	for _, s := range server {
		if s.Name == "FRONTEND" && s.PXname != "stats" {
			frontends = append(frontends, s)
		}
		if s.Name == "BACKEND" && s.PXname != "stats" {
			backends = append(backends, s)
		}
		if !slices.Contains(frontends, s) && !slices.Contains(backends, s) && s.PXname != "stats" {
			servers = append(servers, s)
		}
	}
	return frontends, backends, servers
}
