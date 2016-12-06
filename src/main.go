package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gofilemon/daemon"
	"gofilemon/logger"
	"gofilemon/protocol"
	"gofilemon/server"
	"gofilemon/socket"
)

var BuildVersion = "unknown"

func main() {

	// print version and exit
	if len(os.Args) == 2 && (os.Args[1] == `-v` || os.Args[1] == `--version`) {
		fmt.Printf("%s version: %s\n", os.Args[0], BuildVersion)
		os.Exit(1)
	}

	// print help and exit
	if len(os.Args) == 2 && (os.Args[1] == `-h` || os.Args[1] == `--help`) {
		fmt.Printf("%s <path-to-file> <regexp> <command> <command-args...>\n", os.Args[0])
		os.Exit(1)
	}

	// if socket is not alive, create daemon
	if !socket.Alive() {
		fmt.Fprintf(os.Stderr, "Socket %s is not alive, start as daemon\n", socket.FilePath())
		if err := daemon.Daemonize("/"); err != nil {
			fmt.Fprintf(os.Stderr, "Daemonize error: %s\n", err.Error())
			os.Exit(2)
		}
	}

	// configure Log
	fd, err := logger.GetLogFD()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't open log file: %s\n", err.Error())
		os.Exit(3)
	}
	log.SetOutput(fd)

	// if socket, start listen server, else connect.
	if daemon.IsDaemon() {

		// run as daemon
		ss, err := socket.NewServer()
		if err != nil {
			log.Printf("[ERROR] Start socket: %s\n", err.Error())
			return
		}

		// run main daemon loop
		server.Run(ss)

	} else {

		// run as client
		client, err := socket.NewClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(4)
		}
		// send command args
		data := protocol.StringSliceToByte(os.Args[1:])
		if err := client.Write(data); err != nil {
			fmt.Fprintf(os.Stderr, "Can't write to socket %s: %s\n", socket.FilePath(), err.Error())
			os.Exit(5)
		}
		// read response from server
		if data, err := client.Read(); err != nil {
			fmt.Fprintf(os.Stderr, "Can't read from socket %s: %s\n", socket.FilePath(), err.Error())
			os.Exit(6)
		} else {
			// print result
			fmt.Printf("%s", data)
			if _, err := strconv.ParseFloat(string(data), 64); err != nil {
				// error
				os.Exit(protocol.ERROR_COMMON_CODE)
			} else {
				// no error
				os.Exit(0)
			}
		}

	}

}
