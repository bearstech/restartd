#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import socket
import sys

"""
First is socker name, trailing is the message

./client.py pim Hello World

"""


client = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
client.connect("/tmp/%s" % sys.argv[1])
client.send(bytes(" ".join(sys.argv[2:]), 'utf8'))
r = client.recv(1024)
print(r)
client.close()
