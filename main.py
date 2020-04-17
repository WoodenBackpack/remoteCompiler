from scapy.all import send,IP,UDP,Raw
import argparse
import socket
import netifaces

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("filename", type=str, nargs=1)
    args = parser.parse_args()
    with open(args.filename[0]) as codeFile:
        content = codeFile.read()
        send(IP(dst="10.1.137.61")/UDP(dport=8000)/Raw(load=content))
        receiverSocket = socket.socket(family=socket.AF_INET, type=socket.SOCK_DGRAM)
        ifaces = netifaces.interfaces()
        ifacesAddr = netifaces.ifaddresses(ifaces[0])[netifaces.AF_INET][0]['addr']
        serverAddr = (ifacesAddr, 8001)
        receiverSocket.bind(serverAddr)
        data, address = receiverSocket.recvfrom(2048)
        print(data.decode("utf-8"))
