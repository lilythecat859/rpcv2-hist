#!/usr/bin/env bash
# AGPL-3.0
set -euo pipefail

NAMESPACE="rpcv2-hist"
DEPLOYMENT="rpcv2-hist"

echo "=> Starting chaos mesh tests"
kubectl apply -f https://mirrors.chaos-mesh.org/latest/install.yaml
kubectl wait --for=condition=ready pod -l app.kubernetes.io/instance=chaos-mesh -n chaos-mesh --timeout=60s

cat <<EOF | kubectl apply -f -
apiVersion: chaos-mesh.org/v1alpha1
kind: PodChaos
metadata:
  name: pod-kill
  namespace: $NAMESPACE
spec:
  action: pod-kill
  mode: random-max-percent
  value: "30"
  duration: "60s"
  selector:
    namespaces:
      - $NAMESPACE
    labelSelectors:
      app: rpcv2-hist
EOF

echo "=> Chaos test running 60s"
sleep 70
kubectl delete podchaos pod-kill -n $NAMESPACE
echo "=> Done"