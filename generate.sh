#!/bin/bash
echo "Usage: generate public and private key on openssl v1 (run \"openssl version\" to check your version)"
echo $(pwd)
openssl genrsa -out private.pem 2048 && \
openssl rsa -in private.pem -pubout -out public.pem && \
mv private.pem $(pwd)/key/my_private.pem && \
mv public.pem $(pwd)/key/my_public.pem