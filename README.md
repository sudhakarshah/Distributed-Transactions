# Distributed-Transactions

Make : ```make```<br />
Start deadlock detector : ```python3 deadlock.py 9999```<br />
Start 5 servers : ```python3 run.py```<br />
Run client : ```./client num(0-9)```<br />

## Example <br />
```./client 0```<br />
```./client 1```<br />

## Features <br />
Atomic, Consistency, Isolation with deadlock detection <br/>

## Client Interface <br />
At start up, the client should automatically connect to all the necessary servers and start accepting commands typed in by the user. The user will execute the following commands:

BEGIN: Open a new transaction, and reply with “OK”.<br/>
SET server.key value: 
<br/>
Set the value of an object with the named key stored on the named server. E.g.:<br/>
```SET A.x 1```<br/>
```SET B.y 2```<br/>


GET server.key: Get the value of the object named by the key on the named server. <br/>
Reply with: server.key = value on a separate line. E.g.:<br/>
```GET A.x```<br/>
```A.x = 1```<br/>

If a query is made to an object that has not previously been SET, you should return NOT FOUND and abort the transaction.<br/>

### COMMIT: <br/>
Commit the transaction, making its results visible to other transactions. 
<br/>The client should reply either with COMMIT OK or ABORTED, in the case that the transaction had to be aborted during the commit process.<br/>
<br/>
### ABORT: <br/>
Abort the transaction. All updates made during the transaction must be rolled back. The client should reply with ABORTED to confirm that the transaction was aborted.<br/>


