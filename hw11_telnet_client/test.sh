#!/usr/bin/env bash
set -xeuo pipefail

go build -o go-telnet

python3 - >/tmp/nc.out <<'PY' &
import socket
import sys

server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
server.bind(("127.0.0.1", 4242))
server.listen(1)

conn, _ = server.accept()
conn.sendall(b"Hello\nFrom\nNC\n")

data = b""
while True:
    chunk = conn.recv(4096)
    if not chunk:
        break
    data += chunk

sys.stdout.buffer.write(data)
PY
NC_PID=$!

sleep 1
(echo -e "I\nam\nTELNET client\n" && cat 2>/dev/null) | ./go-telnet --timeout=5s localhost 4242 >/tmp/telnet.out &
TL_PID=$!

sleep 5
kill ${TL_PID} 2>/dev/null || true
kill ${NC_PID} 2>/dev/null || true
wait ${TL_PID} 2>/dev/null || true
wait ${NC_PID} 2>/dev/null || true

function fileEquals() {
  local fileData
  fileData=$(cat "$1")
  [ "${fileData}" = "${2}" ] || (echo -e "unexpected output, $1:\n${fileData}" && exit 1)
}

expected_nc_out='I
am
TELNET client'
fileEquals /tmp/nc.out "${expected_nc_out}"

expected_telnet_out='Hello
From
NC'
fileEquals /tmp/telnet.out "${expected_telnet_out}"

rm -f go-telnet
echo "PASS"
