#!/usr/bin/env bash
#Start swapper Server only if it is not running
if [ "$(ps -ef | grep -v grep | grep swapper | wc -l)" -le 0 ]
then
 "$Home"/.swapper/bin/swapper
 echo "swapper started"
else
 echo "swapper already Running"
fi