package server

import (
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"gofilemon/protocol"
)

type socketServer interface {
	Run(func(net.Conn))
	Close() error
	Alive() bool
}

type Server struct {
	ss socketServer
	// map[FileName]LinkToObject
	fileProcessors map[string]*fileProcessor
	lock           *sync.Mutex
}

func writeWithTimeout(conn net.Conn, data []byte) {
	conn.SetWriteDeadline(time.Now().Add(100 * time.Millisecond))
	conn.Write(data)
}

// принимаем сообщение, парсим его и даем ответ
func (s *Server) clientHandle(conn net.Conn) {

	defer conn.Close()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("[FATAL] clientHandle recover: %v\n", r)
		}
	}()

	// read request
	buf := make([]byte, 1024)
	n, err := conn.Read(buf[:])
	if err != nil {
		log.Printf("[ERROR] Can't read response: %s\n", err.Error())
		return
	}

	data := buf[0:n]

	// ping pong
	if string(data) == protocol.PING {
		conn.Write([]byte(protocol.PONG))
		return
	}

	// unmarshal request
	commands, err := protocol.ByteToStringSlice(data)
	if err != nil {
		log.Printf("[ERROR] Parse %s: %s\n", data, err.Error())
		return
	}

	if len(commands) < 3 {
		log.Printf("[ERROR] Command line is too small: %#v\n", commands)
		writeWithTimeout(conn, []byte(protocol.ERROR_COMMAND_TOO_SMALL))
		return
	}

	logFile, regexpLine, command := commands[0], commands[1], commands[2]
	commandArgs, commandKey := commands[3:], string(data)

	s.lock.Lock()
	defer s.lock.Unlock()

	var processor *fileProcessor
	if oldFileProcessor, ok := s.fileProcessors[logFile]; ok {
		processor = oldFileProcessor
	} else {
		processor, err = newFileProcessor(logFile)
		if err != nil {
			log.Printf("[ERROR] Create file processor: %s\n", err.Error())
			writeWithTimeout(conn, []byte(protocol.ERROR_CREATE_FILE_PROCESSOR))
			return
		}
		s.fileProcessors[logFile] = processor
	}

	if cmd, ok := processor.commands[commandKey]; ok {
		result := strconv.FormatFloat(cmd.Result(), 'f', 6, 64)
		// общий случай, принтуем накопленное значение
		writeWithTimeout(conn, []byte(result))
		return
	}

	log.Printf("[INFO] Register new command: %s\n", commandKey)
	if err := processor.registerCommand(commandKey, command, regexpLine, commandArgs); err != nil {
		log.Printf("[ERROR] Can't register command: %s\n", err.Error())
		writeWithTimeout(conn, []byte(protocol.ERROR_COMMAND))
		return
	}

	writeWithTimeout(conn, []byte(protocol.ERROR_ONLY_REGISTER))
	return

}

func Run(ss socketServer) {

	server := &Server{
		ss:             ss,
		lock:           &sync.Mutex{},
		fileProcessors: make(map[string]*fileProcessor, 0),
	}
	go server.ss.Run(server.clientHandle)

	healTickChan := time.NewTicker(time.Second).C
	for {
		select {
		case <-healTickChan:
			if !server.ss.Alive() {
				log.Printf("[ERROR] Server is not alive, exiting now...\n")
				server.ss.Close()
				return
			}
		}
	}

}
