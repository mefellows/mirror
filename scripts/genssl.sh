#!/bin/bash -e

# Required
domain="*.talent.seek.com.au.dev"
commonname=$domain
SERVER_KEY="server.key"
password="dummypassword"
 
# Company Details
country="AU"
state="VIC"
locality="Melbourne"
organization="SEEK"
organizationalunit="IT"
email="mario@seek.com.au"

echo "Generating a custom CA for $domain"
openssl genrsa -des3 -passout pass:$password -out ca.key 1024
openssl req -new -key ca.key -out ca.csr -passin pass:$password \
    -subj "/C=$country/ST=$state/L=$locality/O=$organization/OU=$organizationalunit/CN=$commonname/emailAddress=$email"
openssl x509 -req -days 365 -in ca.csr -out ca.crt -signkey ca.key -CAserial file.srl -passin pass:$password


echo "Generating key request for $domain"
 
# Generate a key
openssl genrsa -des3 -passout pass:$password -out $SERVER_KEY 1024 -noout

# Remove passphrase from the key.
echo "Removing passphrase from key"
cp $SERVER_KEY $SERVER_KEY.orig
openssl rsa -in $SERVER_KEY.orig -passin pass:$password -out $SERVER_KEY

# Create the request
echo "Creating CSR"
openssl req -new -key server.key -out server.csr -passin pass:$password \
    -subj "/C=$country/ST=$state/L=$locality/O=$organization/OU=$organizationalunit/CN=$commonname/emailAddress=$email"

echo "Creating Cert"
openssl x509 -req -days 10000 -in server.csr -signkey server.key -out server.crt

echo "---------------------------"
echo "-----Below is your CSR-----"
echo "---------------------------"
echo
cat server.csr
 
echo
echo "---------------------------"
echo "-----Below is your Key-----"
echo "---------------------------"
echo
cat $SERVER_KEY

echo "---------------------------"
echo "-----Below is your Cert----"
echo "---------------------------"
echo
cat server.crt



# Client Certs - need to be able to script this too...
# openssl genrsa -des3 -out client.key 1024 
# If encrypted?
# # openssl genrsa -out client.key 1024
# openssl req -key client.key -new -out client.req
# openssl x509 -req -in client.req -CA ca.crt -CAkey ca.key -CAserial file.srl -out client.pem