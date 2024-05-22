package controllers

import (
	"reflect"
	"testing"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCronJob(t *testing.T) {

	// Test scenarios
	// Create cronjob as specified
	// Detect cronjob config diff
	// Update cronjob if different

	t.Run("create cronjob", func(t *testing.T) {
		got := NewCronJob()
		want := batchv1.CronJob{
			ObjectMeta: mockObjectMetadata(),
		}
		if !reflect.DeepEqual(got.ObjectMeta, want.ObjectMeta) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

}

func mockObjectMetadata() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      "mockCronJob",
		Namespace: "testnamespace",
	}
}
