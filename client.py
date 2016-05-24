#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import socket


client = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
client.connect("/tmp/pim")
client.send(b"Hello world")
r = client.recv(1024)
print(r)
client.close()
