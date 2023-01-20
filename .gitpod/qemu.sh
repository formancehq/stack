#!/bin/bash

set -xeuo pipefail

script_dirname="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
outdir="${script_dirname}/_output"

sudo qemu-system-x86_64 -kernel "/boot/vmlinuz" \
-boot c -m 8192M -hda "${outdir}/rootfs/jammy-server-cloudimg-amd64.img" \
-net user \
-smp 8 \
-append "root=/dev/sda rw console=ttyS0,115200 acpi=off nokaslr" \
-nic user,hostfwd=tcp::2222-:22,hostfwd=tcp::80-:80,hostfwd=tcp::443-:443 \
-serial mon:stdio -display none
