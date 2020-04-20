from scapy.all import send,IP,UDP,Raw
import argparse
import socket
import netifaces
from RemoteCodeCompilerProtocol import RCCP
from enum import Enum
from sys import exit
import select

FileExtension = {"py" : 1, "cpp" : 0}


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("filename", type=str, nargs=1)
    args = parser.parse_args()
    fileExtension = args.filename[0].split(".")[2]
    if (fileExtension == "cpp"):
        fileType = FileExtension["cpp"]
    elif (fileExtension == "py"):
        fileType = FileExtension["py"]
    else:
        print("Err")
        exit(1)
    with open(args.filename[0]) as codeFile:
        content = codeFile.read()
        send(IP(dst="10.1.137.61")/UDP(dport=8000)/RCCP(fileType)/Raw(load=content))
        receiverSocket = socket.socket(family=socket.AF_INET, type=socket.SOCK_DGRAM)
        ifaces = netifaces.interfaces()
        ifacesAddr = netifaces.ifaddresses(ifaces[0])[netifaces.AF_INET][0]['addr']
        serverAddr = (ifacesAddr, 8001)
        receiverSocket.bind(serverAddr)
        timeout = 10
        ready = select.select([receiverSocket], [], [], timeout)
        if ready[0]:
            data = receiverSocket.recv(2048)
            print(data.decode("utf-8"))
            exit(0)
        print("timeout!\n")
        exit(2)
