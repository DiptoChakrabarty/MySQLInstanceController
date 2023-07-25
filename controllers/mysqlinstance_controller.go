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

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	mysqlv1alpha1 "github.com/DiptoChakrabarty/MySQLInstanceController.git/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// MySQLInstanceReconciler reconciles a MySQLInstance object
type MySQLInstanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type MySQLInstanceConfig struct {
	Name      string
	Namespace string
}

//+kubebuilder:rbac:groups=dipto.mysql.example.dipto.mysql.example,resources=mysqlinstances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dipto.mysql.example.dipto.mysql.example,resources=mysqlinstances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dipto.mysql.example.dipto.mysql.example,resources=mysqlinstances/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

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

	mysqlInstanceConfig := MySQLInstanceConfig{Name: instance.Name, Namespace: instance.Namespace}

	// Check if StatefulSet exists
	statefulset := &appsv1.StatefulSet{}
	err = rtx.Get(context.TODO(), types.NamespacedName{
		Namespace: req.Namespace,
		Name:      instance.Name,
	}, statefulset)

	// Create the statefulset
	if err != nil {
		err = rtx.CreateMySQLStatefulset(mysqlInstanceConfig)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	// Check if secret present
	secret := &corev1.Secret{}
	err = rtx.Get(context.TODO(), types.NamespacedName{
		Namespace: req.Namespace,
		Name:      instance.Name + "-secret",
	}, secret)

	// Create the secret
	if err != nil {
		err = rtx.CreateMySQLSecret(mysqlInstanceConfig)
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MySQLInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.MySQLInstance{}).
		Complete(r)
}

func generateRandomPassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@!_$#&"
	const passwordLength = 16

	pwd := make([]byte, passwordLength)
	for i := range pwd {
		pwd[i] = charset[rand.Intn(len(charset))]
	}
	return string(pwd)
}
