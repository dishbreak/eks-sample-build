#!/bin/bash

set -eux

mkdir -p ../generated

cd ../generated

config_file=$(aws ec2 \
    export-client-vpn-client-configuration \
    --client-vpn-endpoint $CLIENT_ENDPOINT_ID \
    --output text)

cat <<EOF > ./client.ovpn
$config_file

<cert>
$CERT
</cert>

<key>
$PRIVATE_KEY_PEM
</key>
EOF
