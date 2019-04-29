import asyncio
import sys
import random
import re
import time

if len(sys.argv) != 2:
    print("Usage: deadlock.py <port_number>")
    sys.exit(1)

port = int(sys.argv[1])

lock = asyncio.Lock()

# let g be a dict client -> set(obj)
g_client = {}
g_obj = {}

def will_deadlock(g_client, g_obj, client, obj, lock_type):
    if obj not in g_obj:
        return False
    if client not in g_client:
        g_client[client] = set()
    client_objs = g_client[client]

    cs = g_obj[obj]
    obj_set = set()
    for c, ltype in cs:
        if "W" in (lock_type + ltype):
            obj_set = obj_set.union(g_client[c])
    del_list = []
    for obj_tup in obj_set:
        if obj_tup[0] == obj:
            del_list.append(obj_tup)
    for t in del_list:
        obj_set.remove(t)
    return conflict(client_objs, obj_set)

def conflict(set1, set2):
    # convert to dict
    for obj1, lt1 in set1:
        for obj2, lt2 in set2:
            if obj1 == obj2 and (lt1 == "W" or lt2 == "W"):
                return True
    return False
    

def remove_client(g_client, g_obj, client):
    if client not in g_client:
        return
    rlist = g_client[client]
    g_client[client] = set()
    for obj in rlist:
        remove_list = []
        obj_name = obj[0]
        if obj_name in g_obj:
            r_varient = (client, "R")
            w_varient = (client, "W")
            if r_varient in g_obj[obj_name]:
                g_obj[obj_name].remove(r_varient)
            if w_varient in g_obj[obj_name]:
                g_obj[obj_name].remove(w_varient)

def add_edge(g_client, g_obj, client, obj, ltype):
    if client not in g_client:
        g_client[client] = set()
    if obj not in g_obj:
        g_obj[obj] = set()
    g_client[client].add((obj, ltype))
    g_obj[obj].add((client, ltype))
    return

async def handle_connection(reader, writer):
    addr = writer.get_extra_info('peername')
    print(f"Received connection from {addr}")

    try:
        while True:
            connect_line = await reader.readline()
            connect_line = connect_line.strip().decode()
            print(connect_line)
            token = connect_line.split(" ")

            await lock.acquire()
            try:
                if token[0] == "ADD":
                    if will_deadlock(g_client, g_obj, token[1], token[2], token[3]): # token[1] = client, token[2] = obj
                        writer.write("NO\n".encode())
                        remove_client(g_client, g_obj, token[1]) # token[1] = client
                    else:
                        add_edge(g_client, g_obj, token[1], token[2], token[3])
                        writer.write("YES\n".encode())

                elif token[0] == "REMOVE":
                    remove_client(g_client, g_obj, token[1]) # token[1] = client
                    writer.write("YES\n".encode())
                print(g_client)
                print(g_obj)
            finally:
                lock.release()
            await writer.drain()
    except ConnectionError:
        print(f"Error in connection with {addr}")
    finally:
        writer.close()


loop = asyncio.get_event_loop()
coro = asyncio.start_server(handle_connection, None, port, loop=loop)
server = loop.run_until_complete(coro)

# Serve requests until Ctrl+C is pressed
print('Serving on {}'.format(server.sockets[0].getsockname()))
try:
    loop.run_forever()
except KeyboardInterrupt:
    pass

# Close the server
server.close()
loop.run_until_complete(server.wait_closed())
loop.close()
