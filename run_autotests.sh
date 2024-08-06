#!/bin/bash

PWD=`pwd`
CI_PROG_DIR=${GITHUB_WORKSPACE:-${PWD}}
echo "Start autotests in directory $CI_PROG_DIR"

SERVER_SECRET="ServerKey!12>Au{mL736}"

build/autotests -test.v -test.run=^TestRegistration$ -gophkeeper-binary-path=${CI_PROG_DIR}/build/server \
    -gophkeeper-tls-key=${CI_PROG_DIR}/keys/server-key.pem \
    -gophkeeper-tls-cert=${CI_PROG_DIR}/keys/server-cert.pem \
    -gophkeeper-tls-ca-cert=${CI_PROG_DIR}/keys/ca-cert.pem \
    -gophkeeper-server-secret=${SERVER_SECRET}

res=$?

exit $res



