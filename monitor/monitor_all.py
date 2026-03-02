import subprocess
import socket
import sys

def run(cmd):
    try:
        result = subprocess.run(
            cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE
        )
        return result.returncode == 0, result.stdout.decode().strip()
    except Exception:
        return False, ""

def check_port(host, port):
    try:
        with socket.create_connection((host, port), timeout=2):
            return True
    except Exception:
        return False

print("============== SYSTEM HEALTH CHECK ==============\n")

# ---------- Docker Engine ----------
print("[Docker Engine]")
ok, _ = run("docker info")
if ok:
    print("✔ Docker engine running")
else:
    print("✘ Docker engine NOT running")
    sys.exit(1)
print()

# ---------- Docker Containers ----------
print("[Docker Containers]")
ok, containers = run('docker ps --format "{{.Names}}"')
container_list = containers.splitlines() if containers else []

if container_list:
    for c in container_list:
        print(f"✔ {c} running")
else:
    print("✘ No containers running")
print()

# ---------- Kafka ----------
print("[Kafka]")
if any("kafka" in c.lower() for c in container_list):
    print("✔ Kafka container running")
    if check_port("127.0.0.1", 9092):
        print("✔ Kafka port 9092 reachable")
    else:
        print("✘ Kafka port 9092 NOT reachable")
else:
    print("✘ Kafka container NOT running")
print()

# ---------- MySQL ----------
print("[MySQL]")
if any("mysql" in c.lower() for c in container_list):
    print("✔ MySQL container running")
    ok, _ = run("docker exec local-mysql mysqladmin ping -uroot -proot")
    if ok:
        print("✔ MySQL responding")
    else:
        print("✘ MySQL NOT responding")
else:
    print("✘ MySQL container NOT running")
print()

# ---------- Producer API ----------
print("[Producer API]")
if check_port("127.0.0.1", 8080):
    print("✔ Producer API running on port 8080")
else:
    print("✘ Producer API NOT running")
print()

# ---------- Consumer ----------
print("[Kafka Consumer]")
ok, procs = run("tasklist")
consumer_found = False

if ok:
    for line in procs.splitlines():
        if "go.exe" in line.lower():
            consumer_found = True
            print("✔ Consumer process running")

if not consumer_found:
    print("✘ Consumer process NOT running")

print("\n============== CHECK COMPLETE ====================")