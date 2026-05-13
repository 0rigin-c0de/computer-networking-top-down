package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

const port = ":9000"
const saveDir = "./received_files"

func main() {
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		fmt.Println("Error creating save directory:", err)
		os.Exit(1)
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("File transfer server started on port 9000")
	fmt.Println("Saving received files to:", saveDir)
	fmt.Println("Waiting for connections...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("New connection from: %s\n", clientAddr)

	fileName, err := readFileName(conn)
	if err != nil {
		fmt.Println("Error reading file name:", err)
		return
	}

	savePath := buildSavePath(fileName)

	fmt.Printf("Receiving file: %s\n", fileName)

	outFile, err := os.Create(savePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outFile.Close()

	bytesReceived, err := io.Copy(outFile, conn)
	if err != nil {
		fmt.Println("Error receiving file data:", err)
		return
	}

	fmt.Printf("File saved: %s (%d bytes)\n\n", savePath, bytesReceived)
}

func readFileName(reader io.Reader) (string, error) {
	fileNameLengthBuffer := make([]byte, 1)
	_, err := io.ReadFull(reader, fileNameLengthBuffer)
	if err != nil {
		return "", err
	}

	fileNameLength := int(fileNameLengthBuffer[0])
	fileNameBuffer := make([]byte, fileNameLength)

	_, err = io.ReadFull(reader, fileNameBuffer)
	if err != nil {
		return "", err
	}

	return string(fileNameBuffer), nil
}

func buildSavePath(fileName string) string {
	return filepath.Join(saveDir, filepath.Base(fileName))
}
