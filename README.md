# cs425mp3

High level idea:
Idea 1:
Client connects to multiple servers and a single coordinator.
The coordinator is the one that arbitrates the locking.
Therefore, if a client wants to lock, it must request

server ignores all requests that are not for them
