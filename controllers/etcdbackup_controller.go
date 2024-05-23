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

package controllers

import (
	"context"

	backupconfigv1alpha1 "github.com/SDBrett/etcd-backup-operator/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctrl "sigs.k8s.io/controller-runtime"

	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

// EtcdBackupReconciler reconciles a EtcdBackup object
type EtcdBackupReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=backupconfig.sdbrett.com,resources=etcdbackups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=backupconfig.sdbrett.com,resources=etcdbackups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=backupconfig.sdbrett.com,resources=etcdbackups/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete

func (r *EtcdBackupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrllog.FromContext(ctx)

	// Fetch the Memcached instance
	etcdBackup := &backupconfigv1alpha1.EtcdBackup{}
	err := r.Get(ctx, req.NamespacedName, etcdBackup)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("etcdbackup resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get etcdbackup")
		return ctrl.Result{}, err
	}

	// Check if the deployment already exists, if not create a new one
	found := &batchv1.CronJob{}
	err = r.Get(ctx, types.NamespacedName{Name: etcdBackup.Name, Namespace: etcdBackup.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		cj := r.cronJobForEtcdBackup(etcdBackup)
		log.Info("Creating a new Cronjob", "CronJob.Namespace", cj.Namespace, "CronJob.Name", cj.Name)
		err = r.Create(ctx, cj)
		if err != nil {
			log.Error(err, "Failed to create new CronJob", "CronJob.Namespace", cj.Namespace, "CronJob.Name", cj.Name)
			return ctrl.Result{}, err
		}
		// CronJob created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get CronJob")
		return ctrl.Result{}, err
	}

	// Ensure the deployment size is the same as the spec
	// size := memcached.Spec.Size
	// if *found.Spec.Replicas != size {
	// 	found.Spec.Replicas = &size
	// 	err = r.Update(ctx, found)
	// 	if err != nil {
	// 		log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
	// 		return ctrl.Result{}, err
	// 	}
	// 	// Ask to requeue after 1 minute in order to give enough time for the
	// 	// pods be created on the cluster side and the operand be able
	// 	// to do the next update step accurately.
	// 	return ctrl.Result{RequeueAfter: time.Minute}, nil
	// }

	// Update the Memcached status with the pod names
	// List the pods for this memcached's deployment
	// podList := &corev1.PodList{}
	// listOpts := []client.ListOption{
	// 	client.InNamespace(memcached.Namespace),
	// 	client.MatchingLabels(labelsForMemcached(memcached.Name)),
	// }
	// if err = r.List(ctx, podList, listOpts...); err != nil {
	// 	log.Error(err, "Failed to list pods", "Memcached.Namespace", memcached.Namespace, "Memcached.Name", memcached.Name)
	// 	return ctrl.Result{}, err
	// }
	// podNames := getPodNames(podList.Items)

	// // Update status.Nodes if needed
	// if !reflect.DeepEqual(podNames, memcached.Status.Nodes) {
	// 	memcached.Status.Nodes = podNames
	// 	err := r.Status().Update(ctx, memcached)
	// 	if err != nil {
	// 		log.Error(err, "Failed to update Memcached status")
	// 		return ctrl.Result{}, err
	// 	}
	// }

	return ctrl.Result{}, nil
}

// cronJobForEtcdBackup returns a etcdbackup CronJob object
func (r *EtcdBackupReconciler) cronJobForEtcdBackup(m *backupconfigv1alpha1.EtcdBackup) *batchv1.CronJob {
	ls := labelsForEtcdBackup(m.Name)
	cj := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: batchv1.CronJobSpec{

			Schedule: "",
			JobTemplate: batchv1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: ls,
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{{
								Name:    "etcdbackup",
								Image:   "quay.io/coreos/etcd:v3.5.13-amd64",
								Command: []string{"/bin/bash", "-c", "sleep 360"},
							}},
						},
					},
				},
			},
		},
	}
	// Set Memcached instance as the owner and controller
	ctrl.SetControllerReference(m, cj, r.Scheme)
	return cj
}

// labelsForEtcdBackup returns the labels for selecting the resources
// belonging to the given etcdbackup CR name.
func labelsForEtcdBackup(name string) map[string]string {
	return map[string]string{"app": "etcdbackup", "etcdbackup_cr": name}
}

// SetupWithManager sets up the controller with the Manager.
func (r *EtcdBackupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&backupconfigv1alpha1.EtcdBackup{}).
		Owns(&batchv1.CronJob{}).
		Complete(r)
}
