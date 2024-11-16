#!/bin/bash
echo "Usage: generate public and private key on openssl v3 (run \"openssl version\" to check your version)"
echo $(pwd)
openssl genrsa -traditional -out private.pem 2048 && \
openssl rsa -traditional -in private.pem -out private_unencrypted.pem -outform PEM && \
openssl rsa -traditional -in private_unencrypted.pem -pubout -out public.pem && \
mv private_unencrypted.pem $(pwd)/key/my_private.pem && \
mv public.pem $(pwd)/key/my_public.pem