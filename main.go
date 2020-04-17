package main


import(
	"fmt"
	"net"
	"os/exec"
	"os"
)

func getPortForListening() (net.PacketConn, error) {
	fmt.Printf("listening for packets....\n")
	pc, err := net.ListenPacket("udp", ":8000")
	if err != nil {
		return nil, err
	}
	return pc, nil
}

func receivePacket(pc net.PacketConn) ([]byte, net.Addr, error) {
	buf := make([]byte, 1024)
	_, addr, err := pc.ReadFrom(buf)
	if (err != nil) {
		return nil, nil, err
	}
	fmt.Printf("received packet!\n")
	pc.WriteTo(buf, addr)
	return buf, addr, nil
}

type FileType string

const (
	CppFile FileType = "CppFile"
	PythonFile = "PythonFile"
	LOCAL_VLAN_SOCKET_ADDRESS string = "10.1.137.61:8001"
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

func compileMainAndRun() (string) {
	stdOut, errStr := exec.Command("g++", "-g", "main.cpp", "-oout").CombinedOutput()
	var output string
	if errStr != nil {
		output += "compilation status:\n" + errStr.Error() + "\n"
		output = "compilation errors:\n" + string(stdOut) + "\n"
	} else {
		output = "compilation success !\n"
	}
	execOut, errExecStr := exec.Command("./out").CombinedOutput()
	if errStr == nil {
	    //output += "Execution output:\n" + string(execOut) + "\n"
	    gdbOut, _ := exec.Command("gdb", "-batch", "-ex", "run", "-ex", "bt", "./out").CombinedOutput()
	    output += string(gdbOut)
	} else {
	    output += "Execution err output:\n" + errExecStr.Error() + "\n"
	    output += "Execution output:\n" + string(execOut) + "\n"
	}
	return output
}

func sendResponseBack(address net.Addr, output string) {
  responseAddr, err := net.ResolveUDPAddr("udp", address.String())
  if err != nil {
      fmt.Printf("cannot resolve addr! %s", err)
  }
  localAddr, _ := net.ResolveUDPAddr("udp", LOCAL_VLAN_SOCKET_ADDRESS)
  responseAddr.Port = 8001
  connection, errorConnection := net.DialUDP("udp", localAddr ,responseAddr)
  if (errorConnection != nil) {
      fmt.Printf("error in listening! %s\n", errorConnection)
  }
  _, errorWriter := connection.Write([]byte(output))
  if (errorWriter != nil) {
      fmt.Printf("error in writing! %s\n", errorWriter)
  }
  connection.Close()
}

func main() {
  listenConfig, err := getPortForListening()
  if err != nil {
    fmt.Printf("err: %s\n", err)
    return
  }
  for {
    packet, addr, err := receivePacket(listenConfig)
    if (err != nil) {
      fmt.Printf("error with receiving packets: %s\n", err)
      return
    }
    file, err := createFileOfTypeFromBuffer(CppFile, packet)
    if err == nil {
      output := compileMainAndRun()
      sendResponseBack(addr, output)
    } else {
      fmt.Printf("cant write to file err: %s", err)
    }
    file.Close()
    os.Remove("main.cpp")
    os.Remove("out")
  }
}
