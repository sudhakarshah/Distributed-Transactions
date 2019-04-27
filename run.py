import os
import sys
from time import sleep
from subprocess import Popen

server_port = 5000

commands = []
dir_path = os.path.dirname(os.path.realpath(__file__))


for i in range(5):
    cs = "%s/server %s %d" % (dir_path,chr(65+i), server_port+i)
    commands.append(cs)

procs = []
for i in commands:
    procs.append(Popen(i,shell=True))
for p in procs:
    p.wait()
    print(p)
