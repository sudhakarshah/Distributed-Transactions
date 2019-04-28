import socket
import sys

clientsocket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
clientsocket.connect(('localhost', 9999))
while True:
    s = sys.stdin.readline().strip()
    s+="\n"
    print(s)
    clientsocket.send(s.encode())
    print(clientsocket.recv(1024))
