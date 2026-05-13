package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

const serverAddress = "localhost:9000"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filepath>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		os.Exit(1)
	}

	fileName := fileInfo.Name()
	fmt.Printf("Preparing to send file: %s (%d bytes)\n", fileName, fileInfo.Size())

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Connected to server at %s\n", serverAddress)

	err = writeFileName(conn, fileName)
	if err != nil {
		fmt.Println("Error sending file name:", err)
		os.Exit(1)
	}

	bytesSent, err := io.Copy(conn, file)
	if err != nil {
		fmt.Println("Error sending file:", err)
		os.Exit(1)
	}

	fmt.Printf("File sent successfully. %d bytes transferred.\n", bytesSent)
}

func writeFileName(writer io.Writer, fileName string) error {
	fileNameBytes := []byte(fileName)

	if len(fileNameBytes) > 255 {
		return fmt.Errorf("file name is too long")
	}

	_, err := writer.Write([]byte{byte(len(fileNameBytes))})
	if err != nil {
		return err
	}

	_, err = writer.Write(fileNameBytes)
	if err != nil {
		return err
	}

	return nil
}
