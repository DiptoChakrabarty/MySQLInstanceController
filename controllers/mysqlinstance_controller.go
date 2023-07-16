/*
Copyright 2023 DiptoChakrabarty.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"math/rand"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	mysqlv1alpha1 "github.com/DiptoChakrabarty/MySQLInstanceController.git/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
)

// MySQLInstanceReconciler reconciles a MySQLInstance object
type MySQLInstanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=dipto.mysql.example.dipto.mysql.example,resources=mysqlinstances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dipto.mysql.example.dipto.mysql.example,resources=mysqlinstances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dipto.mysql.example.dipto.mysql.example,resources=mysqlinstances/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MySQLInstance object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (rtx *MySQLInstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	// Fetch MySQLInstance resource
	instance := &mysqlv1alpha1.MySQLInstance{}
	err := rtx.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Check if StatefulSet exists
	statefulset := &appsv1.StatefulSet{}
	err = rtx.Get(context.TODO(), types.NamespacedName{
		Namespace: req.Namespace,
		Name:      instance.Name,
	}, statefulset)

	if err != nil {
		err = rtx.CreateStatefulSet(instance)
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// Create StatefulSet method
func (rtx *MySQLInstanceReconciler) CreateStatefulSet(instance *mysqlv1alpha1.MySQLInstance) error {
	// Generate a random password for the MySQL root user
	rand.Seed(time.Now().UnixNano())
	password := generateRandomPassword()

	// Create StatefulSet
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    3,
			ServiceName: instance.Name,
		},
	}

	// Set password as an environment variable

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MySQLInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.MySQLInstance{}).
		Complete(r)
}
