#!/bin/bash

echo "Creating server.key (ECDSA 384-bit)"
openssl ecparam -genkey -name secp384r1 -out server.key

echo "Creating server.crt (self-signed, valid for 365 days)"
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365 -batch
