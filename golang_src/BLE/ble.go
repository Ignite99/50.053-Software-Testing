package BLE

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func BLE_Zephyr_Handler() {
	// Create a pipe for communicating with the Python script
	pipeOut, pipeIn, err := os.Pipe()
	if err != nil {
		fmt.Println("Error creating pipe:", err)
		return
	}

	cmd := exec.Command("/usr/bin/python3", "./BLE/ble_tester.py", "tcp-server:127.0.0.1:9000")

	cmd.Stdin = pipeIn
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting Python script:", err)
		return
	}

	pipeIn.Close()
	buf := make([]byte, 4096)
	for {
		n, err := pipeOut.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading from pipe:", err)
			return
		}
		fmt.Print(string(buf[:n]))
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("Python script exited with error:", err)
	}
}
