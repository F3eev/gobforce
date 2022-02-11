#!/bin/bash

cat ./ip.txt | while read LINE; do

rustscan -b 50000  -t 2500 -a $LINE -r 1-65535 --scan-order random  -- -sV -oX ./out/$LINE.xml

done