#!/bin/bash -eux

xe vm-import filename=/opt/packer/"$1"-iso/"$1".xva
UUID_VM=$(xe vm-list | grep -B 1 "$1" | grep uuid | awk '{print $5}')
xe vm-param-set name-label="$2" ha-restart-priority=restart ha-always-run=true uuid=$UUID_VM
xe vm-start vm="$2"
