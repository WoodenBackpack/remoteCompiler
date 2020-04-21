package main


import(
	"fmt"
	"net"
	"os/exec"
	"os"
  "errors"
  "bytes"
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

func createFilename(fType FileType) (string, error) {
	var filename string
	if (fType == CppFile) {
		filename = "main.cpp"
	} else if (fType == PythonFile) {
		filename = "main.py"
	} else {
    err := errors.New("Unrecognized FileType!")
    return "", err
  }
  return filename, nil
}

func createFileAndRun(fType FileType, content []byte) (string, string) {
  filename, err := createFilename(fType)
	if err != nil {
		return "", err.Error()
	}
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("cannot create file\n")
		return "", err.Error()
	}
  strContent := string(bytes.Trim(content, "\x00"))
	_, err = file.WriteString(strContent)
	if err != nil {
		fmt.Printf("cannot write to file\n")
		return "", err.Error()
	}
  var stdOutput, errOutput string = "", ""
	if (fType == CppFile) {
	  stdOut, stdErr := exec.Command("g++", "-g", "main.cpp", "-oout").CombinedOutput()
	  gdbOut, gdbErr := exec.Command("gdb", "-batch", "-ex", "run", "-ex", "bt", "./out").CombinedOutput()
    stdOutput = string(stdOut) + "\n" + string(gdbOut) + "\n"
    if (stdErr != nil) {
      errOutput += stdErr.Error() + "\n"
    }
    if (gdbErr != nil) {
      errOutput += gdbErr.Error() + "\n"
    }
    os.Remove("main.cpp")
    os.Remove("out")
	} else if (fType == PythonFile) {
	  pyStdOut, pyStdErr := exec.Command("python", "main.py").CombinedOutput()
    if (pyStdErr != nil) {
      errOutput += pyStdErr.Error() + "\n"
    }
    stdOutput += string(pyStdOut) + "\n"
    os.Remove("main.py")
	}
  file.Close()
  return stdOutput, errOutput
}

func compileAndRun(fType FileType, content []byte) (string) {
	stdOut, stdErr := createFileAndRun(fType, content)
	var output string
	if stdErr != "" {
	    output += string(stdOut) + "\nERROR:\n" +string(stdErr)
	} else {
	    output += string(stdOut)
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
    isCpp := packet[0] == 1
    var fileType FileType
    if (isCpp) {
      fileType = CppFile
    } else {
      fileType = PythonFile
    }
    outputStr := compileAndRun(fileType, packet[1:len(packet) - 1])
    sendResponseBack(addr, outputStr)
  }
}
