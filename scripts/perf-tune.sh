#!/usr/bin/env bash
# AGPL-3.0
set -euo pipefail

echo "=> OS-level tuning for AMD 16-core / 256 GB"

# CPU governor
echo performance | tee /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor

# IRQ affinity
systemctl stop irqbalance
set_irq_affinity.sh eth0

# Sysctl
cat >> /etc/sysctl.conf <<EOF
net.core.rmem_max = 134217728
net.core.wmem_max = 134217728
net.ipv4.tcp_rmem = 4096 87380 134217728
net.ipv4.tcp_wmem = 4096 65536 134217728
net.core.netdev_max_backlog = 5000
net.ipv4.tcp_congestion_control = bbr
EOF
sysctl -p

# Hugepages
echo 512 > /sys/kernel/mm/hugepages/hugepages-2048kB/nr_hugepages

# Disk scheduler
echo noop > /sys/block/nvme0n1/queue/scheduler
echo noop > /sys/block/nvme1n1/queue/scheduler

echo "=> Done"