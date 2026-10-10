# Building a Calculator Mesh on Kubernetes

*Part 2: My First Operator*

*Part 1: [The Cluster and the Control Plane](part1.md)*

All the code for this series is on GitHub: [github.com/aminehd/calcmesh](https://github.com/aminehd/calcmesh)


## 1. Where we left off and what's new
In Part 1 of this tutorial, I explained Kubernetes manifest files and set up a Kubernetes cluster. Today I wanna talk about what happens after running `kubectl apply ...`. Here is the folder structure:
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

Since writing Kubernetes manifests can be repetitive, we use Helm syntax to write higher-level code (loops, lookups...) and render the bare-bones manifest files from it.
However, the focus of this tutorial is more on the flow of things that happen after the `kubectl apply` command. Let's take a very simple manifest file:
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

[![What happens on kubectl apply](https://raw.githubusercontent.com/aminehd/calcmesh/main/tutorials/images/kubectl-apply.gif)](https://github.com/aminehd/calcmesh/blob/main/tutorials/images/kubectl-apply.gif)
By running `kubectl apply -f pod.yaml`, you send a request to the API server to create a new Pod. The API server itself does not attempt to create the Pod on a node.
Instead, the API server stores the Pod in etcd and notifies the scheduler and the internal controllers. Once the scheduler picks a node for it, it writes that choice back to the API server, which saves it in etcd again. Now the kubelet on that node gets notified that a Pod was assigned to it. The kubelet asks the container runtime to pull the image and start the container, and then reports the Pod's status (Running) back to the API server. Notice that nobody calls anybody directly: every step is "read from the API server, do my part, write back to the API server".

There is another pattern you can use Kubernetes for. You can define new kinds with CRDs (CustomResourceDefinitions).

Then you can run watchers called operators. Operators can run outside your cluster or as a pod in your cluster. Then you can create a manifest of that type, for example:
```yaml
apiVersion: example.com/v1
kind: MyCustomType
metadata:
  name: hellomytype
spec:
  containers:
    - name: custom
      customInput: adflasdfj
```
After applying, the API server will save the object in its memory (etcd). Since this is a new type, it won't go to the built-in controllers; instead it goes to your operator, and then your operator may ask the API server to create pods or other custom resources. If it creates a pod, the API server again saves it in etcd and notifies the built-in controllers.

See the image:

[![kubectl apply with your own operator](https://raw.githubusercontent.com/aminehd/calcmesh/main/tutorials/images/kubectl-apply-operator.gif)](https://github.com/aminehd/calcmesh/blob/main/tutorials/images/kubectl-apply-operator.gif)

Now, what is the point of this back and forth? Why doesn't the API server directly do what you want, instead of a controller? The Kubernetes system architecture is designed to handle failures and to make requests fast, non-blocking and recoverable. These patterns are the state of the art for doing so. Also, Kubernetes was the result of years of building internal tools at Google (Borg, ...).



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
