#!/usr/bin/env bash

set -e
set -o pipefail

cd /data/var/lib/icinga2/certs

if openssl x509 -noout -in ${ICINGA_CN}.crt -ext subjectAltName | grep -q "127.0.0.1"; then
  echo "Certificate for ${ICINGA_CN} already have IP SAN"
else
  echo "Generate new certificate for ${ICINGA_CN} with IP SAN for existent private key"
  rm -f ./${ICINGA_CN}.crt
  rm -f ./${ICINGA_CN}.csr

  cat <<EOF >extfile.cnf
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature,keyEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = ${ICINGA_CN}
IP.1 = 127.0.0.1
EOF

  openssl req -new -subj "/CN=${ICINGA_CN}" -key ${ICINGA_CN}.key -out ${ICINGA_CN}.csr
  openssl x509 -req -in ${ICINGA_CN}.csr -CA ../ca/ca.crt -CAkey ../ca/ca.key -CAcreateserial -out ${ICINGA_CN}.crt -days 365 -sha256 -extfile extfile.cnf
fi
