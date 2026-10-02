# Building a Calculator Mesh on Kubernetes, Part 1: The Cluster and the Control Plane

## 1. The project and why

I am setting out to build the simplest possible service mesh. A lot of big enterprise applications are, in some way, a service mesh. My definition of a service mesh: services that run in isolation, but can still talk to each other, either directly or through a coordinator. The only purpose of this project is learning: first, how to build a service mesh, and second, the technologies that the industry uses for systems like this.

For this project, I use Kubernetes to run the containers, and Envoy as the L7 proxy. The services are as simple as possible: one calculator for each operator, addition, subtraction, multiplication and division. A client calls one main endpoint, and Envoy routes the request to the right calculator. This first tutorial is about bootstrapping the Kubernetes cluster, with only the Kubernetes concepts we need for that.

All the code for this series is on GitHub: [github.com/aminehd/calcmesh](https://github.com/aminehd/calcmesh)


## 2. Kubernetes concepts

The philosophy of Kubernetes in one sentence, from what I have learned: you write down the state you want, and Kubernetes keeps running loops that make the real state match it.

Our desired state lives in our code base, as manifest files. At runtime, Kubernetes has two main parts.

The first part is the control plane. One thing that helped me understand Kubernetes is that the control plane is not our code; it is Kubernetes' own code. It has three building blocks:

- **Controllers** run the loops that reconcile the desired state with the current state.
- **The scheduler** assigns Pods to real nodes.
- **The API server** stores our objects. It is the only door into the cluster: developers, controllers, the scheduler and the nodes all talk to it.

The second part is the worker nodes, where a small agent called the kubelet actually starts our containers.

This control-plane pattern shows up in many other infrastructure systems too. Envoy, for example, also gets its config from a control plane.

When the real state drifts, for example when a Pod dies, the loops notice and fix the difference by creating or deleting Pods. Pods are the building blocks that actually run our containers. Objects like Deployments end up as Pods, while objects like Services and ConfigMaps are not Pods; they describe the networking and config around them.

I think this is a very good mental model to start with. From here, the next step is to learn which building blocks Kubernetes offers, and what people in the area build with them.

If you want to see these ideas explained visually, watch this playlist: [Kubernetes on YouTube](https://www.youtube.com/playlist?list=PLy7NrYWoggjwPggqtFsI_zMAwvG0SqYCb)

Some more concepts:

- **Node:** one machine, physical or virtual. Think of it as a running OS with its own hardware. In kind, a node is just a Docker container pretending to be a machine.
- **Cluster:** the whole thing running: the control plane and the worker nodes.
- **Control plane:** the node and containers that run Kubernetes itself. All the code that the Kubernetes core developers wrote runs here. You set it up once.
- **API server:** one process in the control plane that everyone talks to.
- **Manifest:** a YAML file that describes the Kubernetes objects we need. You can write manifests by hand, or use Helm to generate them from templates with loops and variables. I write them by hand for now; Helm comes in the next part.


## 3. How I use them here

This is the folder structure:

```
calcmesh/
├── bootstrap/
│   └── kind-cluster.yaml        # the cluster: one node
├── calculators/                 # one folder per calculator (placeholders for now)
│   ├── addition/
│   ├── subtraction/
│   ├── multiplication/
│   └── division/
├── envoy/                       # the proxy config goes here (placeholder for now)
├── deployment/
│   └── calcmesh/                # the Helm chart (Part 2)
├── tutorials/
│   └── part1.md
└── README.md
```

To run Kubernetes on my PC, I use kind, which stands for Kubernetes IN Docker. Real clusters run on many machines; kind fakes them. Every node is a Docker container, and inside it runs the real Kubernetes: the same API server, etcd, scheduler and controllers as a cloud cluster. So one PC can have one node or five, and if you break it, you delete it and make a new one in a minute.

In this tutorial there is only one file, kind-cluster.yaml, and it creates a cluster with one node. Notice that this file is not a Kubernetes manifest. It is kind's own config (look at the apiVersion). It never goes to the API server; kind only reads it to know how many nodes to make. Then I create the cluster with kind:

```yaml
# One node: it runs the control plane (API server, etcd, scheduler,
# controller manager) and your pods.
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
```
```bash
kind create cluster --name calcmesh --config bootstrap/kind-cluster.yaml
kubectl get nodes
```

kind also points kubectl at the new cluster, so `kubectl get nodes` shows one node: calcmesh-control-plane. With only one node, it runs both the control plane and, later, our pods.

For the services and for Envoy, there are only placeholders for now. The main point for me is to keep the services isolated from the networking code, which lives in Envoy.


## 4. Conclusion and what's next

To run it yourself, you only need Docker, kind and kubectl. Clone the repo and bring up the cluster:

```bash
git clone https://github.com/aminehd/calcmesh.git
cd calcmesh
kind create cluster --name calcmesh --config bootstrap/kind-cluster.yaml
```

Then look at the control plane. It runs as pods inside our one node:

```bash
kubectl get pods -n kube-system
```

You should see kube-apiserver, etcd, kube-scheduler and kube-controller-manager, the same pieces from section 2. When you are done, `kind delete cluster --name calcmesh` removes everything.

What I take away from this part:

1. This project builds a service mesh of calculators that talk to each other. Kubernetes runs the containers, and Envoy will be the network proxy.
2. A Kubernetes cluster has a control plane. The control plane always watches the objects we sent to the API server, and keeps the running state matched to our desired state.
3. kind runs a real Kubernetes cluster on one PC, where every node is a Docker container.
4. We saw the structure of the code.

In the next part, I use Helm to put the calculators in the cluster. Helm is a package manager for Kubernetes. Instead of writing one Deployment and one Service for each calculator by hand, I write one template with a loop and a values.yaml with the list of calculators. Helm generates all the YAML and sends it to the API server. Adding a new calculator becomes one line, and if something breaks, Helm can roll back to the last version.

If you liked this, follow me on [social media links here] for the next parts, and give the project a star on [GitHub](https://github.com/aminehd/calcmesh) so you get notified when the next part lands.
