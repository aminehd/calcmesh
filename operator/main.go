// A tiny operator: for every Calculator object, keep a Deployment and a
// Service running for its op.
//
// The standard operator shape (the same one kubebuilder generates):
//   SetupWithManager: For(Calculator), Owns(Deployment)
//   Reconcile:        read the wish -> create what is missing -> write status
package main

import (
	"context"
	"fmt"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

const (
	image = "hashicorp/http-echo:1.0" // placeholder: answers with its own name
	port  = 8080
)

// calculatorGVK names our CRD's type. We use an "unstructured" object (a map)
// instead of a Go struct to skip code generation. Real operators use typed
// structs (a *_types.go file) generated with kubebuilder; that's the next step.
var calculatorGVK = schema.GroupVersionKind{Group: "calcmesh.io", Version: "v1", Kind: "Calculator"}

func newCalculator() *unstructured.Unstructured {
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(calculatorGVK)
	return u
}

type CalculatorReconciler struct {
	client.Client
}

// Reconcile runs every time a Calculator, or a Deployment it owns, changes.
func (r *CalculatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := ctrl.LoggerFrom(ctx)

	// 1. Read the wish.
	calc := newCalculator()
	if err := r.Get(ctx, req.NamespacedName, calc); err != nil {
		// Deleted: nothing to do (ownerReferences will clean up its children).
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	op, _, _ := unstructured.NestedString(calc.Object, "spec", "op")
	replicas, found, _ := unstructured.NestedInt64(calc.Object, "spec", "replicas")
	if !found {
		replicas = 1
	}
	logger.Info("reconciling", "op", op, "replicas", replicas)

	// 2. Make the world match the wish.
	// TODO(Amineh): ensure the Deployment exists.
	//   - Try r.Get(ctx, req.NamespacedName, &appsv1.Deployment{})
	//   - If apierrors.IsNotFound(err): build one with
	//       desiredDeployment(req.Name, req.Namespace, op, int32(replicas))
	//     mark it as owned by calc with
	//       ctrl.SetControllerReference(calc, dep, r.Scheme())
	//     and create it with r.Create(ctx, dep)
	//   - Any other error: return it (controller-runtime will retry)
	// Then do the same for the Service with desiredService(...).
	// It must be safe to run many times: only create what is missing.

	// 3. Report back in status.
	// TODO(Amineh, later): set status.ready / status.message and call r.Status().Update(ctx, calc)

	return ctrl.Result{}, nil
}

func desiredDeployment(name, namespace, op string, replicas int32) *appsv1.Deployment {
	labels := map[string]string{"app": name}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: labels},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{Containers: []corev1.Container{{
					Name:  "calculator",
					Image: image,
					Args:  []string{fmt.Sprintf("-listen=:%d", port), "-text=" + op},
					Ports: []corev1.ContainerPort{{ContainerPort: port}},
				}}},
			},
		},
	}
}

func desiredService(name, namespace string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": name},
			Ports:    []corev1.ServicePort{{Port: 80, TargetPort: intstr.FromInt32(port)}},
		},
	}
}

// SetupWithManager is the subscription: wake me for Calculators,
// and for Deployments/Services that a Calculator owns.
func (r *CalculatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(newCalculator()).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	_ = apierrors.IsNotFound // used in your TODO

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{}) // uses your kubeconfig, like kubectl
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := (&CalculatorReconciler{Client: mgr.GetClient()}).SetupWithManager(mgr); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("watching calculators...")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
