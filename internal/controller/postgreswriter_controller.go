/*
Copyright 2024.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	testleryv1 "github.com/nambhd/postgres-writer-operator/api/v1"
)

// PostgresWriterReconciler reconciles a PostgresWriter object
type PostgresWriterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	PostgresDBClient *psql.PostgresDBClient // newly added by me, this will be used by our operator to talk
}

// +kubebuilder:rbac:groups=testlery.testlery.com,resources=postgreswriters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=testlery.testlery.com,resources=postgreswriters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=testlery.testlery.com,resources=postgreswriters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PostgresWriter object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *PostgresWriterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	logger := log.FromContext(ctx)

	// parsing the incoming postgres-writer resource
	postgresWriterObject := &demov1.PostgresWriter{}
	err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, postgresWriterObject)
	if err != nil {
		// if the resource is not found, then just return (might look useless as this usually happens in case of Delete events)
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Error occurred while fetching the PostgresWriter resource")
		return ctrl.Result{}, err
	}

	// parsing the table, name, age, country fields from the spec of the incoming postgres-writer resource
	table, name, age, country := postgresWriterObject.Spec.Table, postgresWriterObject.Spec.Name, postgresWriterObject.Spec.Age, postgresWriterObject.Spec.Country

	// forming a unique id corresponding to the incoming resource
	id := postgresWriterObject.Namespace + "/" + postgresWriterObject.Name

	logger.Info(fmt.Sprintf("Writing name: %s age: %d country: %s into table: %s with id as %s ", name, age, country, table, id))

	// performing the `INSERT` to the DB with the provided name, age, country on the provided table
	if err = r.PostgresDBClient.Insert(id, table, name, age, country); err != nil {
		logger.Error(err, "error occurred while inserting the row in the Postgres DB")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PostgresWriterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&testleryv1.PostgresWriter{}).
		Complete(r)
}
