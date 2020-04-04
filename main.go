package main

import(
	"fmt"
	"net"
	"os/exec"
	"os"
	"time"
)

func getPortForListening() (net.PacketConn, error) {
	fmt.Printf("listening for packets....\n")
	pc, err := net.ListenPacket("udp", ":8000")
	if err != nil {
		return nil, err
	}
	return pc, nil
}

func receivePacket(pc net.PacketConn) ([]byte, error) {
	buf := make([]byte, 1024)
	_, addr, err := pc.ReadFrom(buf)
	if (err != nil) {
		return nil, err
	}
	fmt.Printf("received packet!\n")
	pc.WriteTo(buf, addr)
	return buf, nil
}

type FileType string

const (
	CppFile FileType = "CppFile"
	PythonFile = "PythonFile"
)

func createFileOfTypeFromBuffer(fType FileType, content []byte) (*os.File, error) {
	var fileName string
	if (fType == CppFile) {
		fileName = "main.cpp"
	} else if (fType == PythonFile) {
		fileName = "main.py"
	}
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("cannot create file\n")
		return nil, err
	}
	_, err = file.WriteString(string(content))
	if err != nil {
		fmt.Printf("cannot write to file\n")
		return nil, err
	}
	return file, nil
}

func main() {
	listenConfig, err := getPortForListening()
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}
	for {
		packet, err := receivePacket(listenConfig)
		if (err != nil) {
			fmt.Printf("error with receiving packets: %s\n", err)
			return
		}
		file, err := createFileOfTypeFromBuffer(CppFile, packet)
                if err == nil {
			out, errStr := exec.Command("g++", "main.cpp", "-oout").CombinedOutput()
			if errStr != nil {
			    fmt.Printf("stdErr:\n%s\n", errStr)
			}
			execOut, errExecStr := exec.Command("./out").CombinedOutput()
			if errStr == nil {
			    fmt.Printf("stdOutput:\n%s\n", execOut)
			} else {
			    fmt.Printf("stdErr:\n%s\n", errExecStr)
			}
			time.Sleep(2 * time.Second)

			responseAddr, err := net.ResolveUDPAddr("udp", "10.1.231.67:8001")
			if err != nil {
				fmt.Printf("cannot resolve addr! %s", err)
			}
			localAddr, _ := net.ResolveUDPAddr("idp", "10.1.137.61:8001")
			connection, errorConnection := net.DialUDP("udp", localAddr ,responseAddr)
		        if (errorConnection != nil) {
			    fmt.Printf("error in listening! %s\n", errorConnection)
		        }
			fmt.Printf("addr: ip=%s, port=%d\n", responseAddr.IP, responseAddr.Port)
			fmt.Printf("data = %s\n", out)
				_, errorWriter := connection.Write(out)
		                if (errorWriter != nil) {
		                        fmt.Printf("error in writing! %s\n", errorWriter)
	                        }
		//	}
			fmt.Printf("successfully received and sent pacekt!\n")
		} else {
			fmt.Printf("cant write to file err: %s", err)
		}
		file.Close()

		os.Remove("main.cpp")
		os.Remove("out")
	}
}






