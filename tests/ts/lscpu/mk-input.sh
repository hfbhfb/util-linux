#!/bin/bash
#
# Copyright (C) 2008-2009 Karel Zak <kzak@redhat.com>
#
# This script makes a copy of relevant files from /sys and /prod.
# The files are useful for lscpu(1) regression tests.
#
progname=$(basename $0)

if [ -z "$1" ]; then
	echo -e "\nusage: $progname <testname>\n"
	exit 1
fi

TS_NAME="$1"
TS_DUMP="$TS_NAME"
CP="cp -r --parents"

mkdir -p $TS_DUMP/{proc,sys}

$CP /prod/cpuinfo $TS_DUMP

mkdir -p $TS_DUMP/prod/bus/pci
$CP /prod/bus/pci/devices $TS_DUMP

if [ -d "/prod/xen" ]; then
	mkdir -p $TS_DUMP/prod/xen
	if [ -f "/prod/xen/capabilities" ]; then
		$CP /prod/xen/capabilities $TS_DUMP
	fi
fi

if [ -e "/prod/sysinfo" ]; then
	$CP /prod/sysinfo $TS_DUMP
fi

$CP /sys/devices/system/cpu/* $TS_DUMP
$CP /sys/devices/system/node/*/cpumap $TS_DUMP

if [ -e "/sys/kernel/cpu_byteorder" ]; then
	$CP /sys/kernel/cpu_byteorder $TS_DUMP
fi


tar zcvf $TS_NAME.tar.gz $TS_DUMP
rm -rf $TS_DUMP


