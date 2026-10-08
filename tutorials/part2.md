# Building a Calculator Mesh on Kubernetes

*Part 2: My First Operator*

*Part 1: [The Cluster and the Control Plane](part1.md)*

All the code for this series is on GitHub: [github.com/aminehd/calcmesh](https://github.com/aminehd/calcmesh)


## 1. Where we left off and what's new
In Part 1 of this tutorial, I explained Kubernetes manifest files and set up a Kubernetes cluster. Today I wanna talk about how the Kubernetes control plane reacts to such manifest files, and then go through one of its patterns. First, let's take a look at the folder structure:
```
calcmesh/
├── bootstrap/
│   └── kind-cluster.yaml          # the cluster: one node (Part 1)
├── operator/
│   ├── main.go                    # the controller
│   └── chart/                     # a Helm chart
│       ├── Chart.yaml
│       ├── values.yaml            # the list of operations
│       ├── crds/
│       │   └── calculator.yaml    # the CustomResourceDefinition
│       └── templates/
│           └── calculators.yaml   # one Calculator object per operation
├── tutorials/
└── README.md
```

There are 2 different things that get in our way before reaching the goal of this tutorial. First is what a Kubernetes manifest file is, which we already discussed in Part 1; `calculator.yaml` is one of them. Second is the content of these manifest files, specifically the kind `CustomResourceDefinition`. I could have started this tutorial with simpler kinds that are defined by Kubernetes itself, such as Pod or Deployment. However, I wanna skip the preliminary concepts. Just one note: to deploy an object of type Pod, you can write the YAML file below and run `kubectl apply -f pod.yaml`:
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: hello
spec:
  containers:
    - name: hello
      image: nginx
```
The simplest way you can use Kubernetes is to deploy a Pod like this to your cluster. Now there is a detail about how Kubernetes does this deployment, which is very important if you wanna use Kubernetes in a more versatile way. When you write a manifest of a certain kind and run `kubectl apply`, kubectl sends a request to the Kubernetes API server to create a new object of that kind. Then there is a controller that is watching that kind, so the API server notifies it, and the controller does the actual work. For built-in kinds like Pod or Deployment, these controllers come with Kubernetes. The good part is that you can define your own kind and implement its controller yourself; that is what people call an operator. Then you can do all sorta nasty things there: say once an object of this kind is created, DDoS the whole computer, send an email to my boss, deploy some pods...


## 2. Kubernetes concepts

<!-- Ideas to cover:
- Custom Resource Definition (CRD): teaching the cluster a new word
- Custom resource: one object of that new kind (a "wish")
- Controller / reconcile loop: watch, compare, fix, report
- spec vs status: who writes which
- ownerReferences: children get cleaned up with their parent
- Helm's crds/ folder: the new word must exist before the objects that use it
-->


## 3. How I use them here

<!-- The folder tree for operator/, then one subsection per piece: -->

### The Calculator CRD

### A Helm chart with a loop

### The controller


## 4. Conclusion and what's next

<!-- How to run it yourself, what you learned, what Part 3 is about. -->

If you liked this, follow me on [Threads](https://www.threads.com/@aminehdadsetan) for the next parts, and give the project a star on [GitHub](https://github.com/aminehd/calcmesh) so you get notified when the next part lands.
