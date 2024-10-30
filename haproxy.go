package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)

func main() {
	type Server struct {
			Name string
			IP string
			Port string
			Status string
			Weight string
	}
	myserver := []Server{}

    // Connect to HAProxy admin.socket
    conn, err := net.Dial("unix", "/var/run/haproxy/admin.sock") // Path to your HAProxy socket
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error connecting to HAProxy socket:", err)
        return
    }
    defer conn.Close()

    // Send the "show stat" command to retrieve server stats
    fmt.Fprintln(conn, "show stat")

    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        line := scanner.Text()

        // Skip empty lines and comments
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        // Split CSV fields
        fields := strings.Split(line, ",")
        if len(fields) < 34 {
            continue // Ensure there are enough fields based on HAProxy's output format
        }

        if fields[1] != "BACKEND" && fields[1] != "FRONTEND" {
        // Extracting relevant fields from the CSV output
        //backend := fields[0]
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
		if len(ipport) ==2 {
			ip = ipport[0]
			port = ipport[1]
		}
        //currentSessions := fields[4]
        //maxSessions := fields[5]
		newserver := Server{Name: server, IP: ip , Port: port, Status: status, Weight: weight}
		myserver = append(myserver, newserver)
      }
	  for _, srv := range myserver {
	  fmt.Println(srv)

	  }
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "Error reading from HAProxy socket:", err)
    }
}
