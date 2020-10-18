#!/bin/bash

set -euo pipefail

oc -n openshift-logging extract --confirm  secret/elasticsearch --to=. --keys=admin-ca
oc -n openshift-logging extract --confirm  secret/elasticsearch --to=. --keys=admin-cert
oc -n openshift-logging extract --confirm  secret/elasticsearch --to=. --keys=admin-key

cat <<EOF > "./ocp-es-route.yaml"
apiVersion: route.openshift.io/v1
kind: Route
metadata:
  name: elasticsearch
  namespace: openshift-logging
spec:
  host:
  to:
    kind: Service
    name: elasticsearch
  tls:
    termination: reencrypt
    destinationCACertificate: |
EOF

cat ./admin-ca | sed -e "s/^/      /" >> ./ocp-es-route.yaml

oc -n openshift-logging  delete route elasticsearch ||:
oc -n openshift-logging  create -f ./ocp-es-route.yaml
