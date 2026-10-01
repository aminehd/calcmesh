# calcmesh

A service mesh built one step at a time: calculators that each do one
operation (+, -, *, /) and know nothing about the network.

## Step 0: a cluster with one node

```bash
kind create cluster --name calcmesh --config bootstrap/kind-cluster.yaml
kubectl get nodes
```
