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

So what happens after you run `kubectl apply ...`? You request the Kubernetes API server to create a new object of kind Pod. The API server creates the object and then notifies the controllers watching Pods (the scheduler and the kubelet), and they run their reconcile loops. See the picture:

[![What happens on kubectl apply](https://raw.githubusercontent.com/aminehd/calcmesh/main/tutorials/images/kubectl-apply.gif)](https://github.com/aminehd/calcmesh/blob/main/tutorials/images/kubectl-apply.gif)

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

now what is the point of this back and forth, why not directly api server do what you wnat instead of controller. The kubernetes system archticuture is to handle failures and make request fasts and non blocking and recoverable. These pattersn are state of the art for doing so. Also Kubernetes was the result of years of building internal tools at Google (Borg, ...). 



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
