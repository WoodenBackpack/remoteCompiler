from scapy.all import send,IP,UDP,Raw
import argparse
import socket

if __name__ == "__main__":
    print("start")
    parser = argparse.ArgumentParser()
    parser.add_argument("filename", type=str, nargs=1)
    args = parser.parse_args()
    print(args.filename[0])
    with open(args.filename[0]) as codeFile:
        content = codeFile.read()
        send(IP(dst="10.1.137.61")/UDP(dport=8000)/Raw(load=content))
        receiverSocket = socket.socket(family=socket.AF_INET, type=socket.SOCK_DGRAM)
        serverAddr = ("10.1.231.67", 8001)
        receiverSocket.bind(serverAddr)
        data, address = receiverSocket.recvfrom(4048)
        print(data.decode("utf-8"))
