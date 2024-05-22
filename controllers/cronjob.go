package controllers

import (
	backupconfigv1alpha1 "github.com/SDBrett/etcd-backup-operator/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *EtcdBackupReconciler) etcdBackupCronJob(m *backupconfigv1alpha1.EtcdBackup) *batchv1.CronJob {
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
	// Set EtcdBackup instance as the owner and controller
	ctrl.SetControllerReference(m, cj, r.Scheme)
	return cj
}

func labelsForEtcdBackup(name string) map[string]string {
	return map[string]string{"app": "etcdbackup"}
}

// // getPodNames returns the pod names of the array of pods passed in
// func getPodNames(pods []corev1.Pod) []string {
// 	var podNames []string
// 	for _, pod := range pods {
// 		podNames = append(podNames, pod.Name)
// 	}
// 	return podNames
// }
