package scheduler

import (
	"context"

	"github.com/copito/gocronscheduler/entity"
	"github.com/google/uuid"
	batchV1 "k8s.io/api/batch/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GenerateCronJobDefinition(jobName string, containerImage string, schedule string) *batchV1.CronJob {
	// Set clean up of jobs to 2 seconds irrespective of failing or succeding
	ttlInSeconds := int32(2)
	// Set the activeDeadlineSeconds to 2 minutes (120 seconds)
	maxRunTimeInSeconds := int64(120)

	// BackoffLimit
	BackoffLimit := int32(3)

	// Set limits for job history
	successfulJobsHistoryLimit := int32(3)
	failedJobsHistoryLimit := int32(1)

	cronJob := &batchV1.CronJob{
		ObjectMeta: metaV1.ObjectMeta{
			Name: jobName,
			Labels: map[string]string{
				"job_uuid": uuid.NewString(),
				"type":     "metric",
			},
		},
		Spec: batchV1.CronJobSpec{
			Schedule:                   schedule,
			SuccessfulJobsHistoryLimit: &successfulJobsHistoryLimit,
			FailedJobsHistoryLimit:     &failedJobsHistoryLimit,
			JobTemplate: batchV1.JobTemplateSpec{
				Spec: batchV1.JobSpec{
					ActiveDeadlineSeconds:   &maxRunTimeInSeconds,
					TTLSecondsAfterFinished: &ttlInSeconds,
					BackoffLimit:            &BackoffLimit,
					Template: coreV1.PodTemplateSpec{
						Spec: coreV1.PodSpec{
							Containers: []coreV1.Container{
								{
									Name:  jobName,
									Image: containerImage,
									// Resources: coreV1.ResourceRequirements{
									// 	Requests: coreV1.ResourceList{
									// 		coreV1.ResourceCPU:    *resource.Quantity("100m"),
									// 		coreV1.ResourceMemory: resourceQuantity("128Mi"),
									// 	},
									// 	Limits: coreV1.ResourceList{
									// 		coreV1.ResourceCPU:    resourceQuantity("500m"),
									// 		coreV1.ResourceMemory: resourceQuantity("512Mi"),
									// 	},
									// },
								},
							},
							RestartPolicy: coreV1.RestartPolicyOnFailure,
						},
					},
				},
			},
		},
	}
	return cronJob
}

func CreateCronJob(ctx context.Context, clientSet *kubernetes.Clientset, cronJob *batchV1.CronJob) error {
	_, err := clientSet.BatchV1().CronJobs(entity.Namespace).Create(ctx, cronJob, metaV1.CreateOptions{})
	return err
}

func DeleteCronJob(ctx context.Context, clientSet *kubernetes.Clientset, jobName string) error {
	return clientSet.BatchV1().CronJobs(entity.Namespace).Delete(ctx, jobName, metaV1.DeleteOptions{})
}

func UpdateCronJob(ctx context.Context, clientset *kubernetes.Clientset, jobName string, adjustedCronJob *batchV1.CronJob) (*batchV1.CronJobStatus, error) {
	cronJob, err := clientset.BatchV1().CronJobs(entity.Namespace).Update(ctx, adjustedCronJob, metaV1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return &cronJob.Status, nil
}

func UpdateCronJobName(ctx context.Context, clientset *kubernetes.Clientset, previousJobName string, newJobName string) (*batchV1.CronJobStatus, error) {
	cronJob, err := clientset.BatchV1().CronJobs(entity.Namespace).Get(ctx, previousJobName, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// Change Job Name
	cronJob.Name = newJobName
	cronJob, err = clientset.BatchV1().CronJobs(entity.Namespace).Update(ctx, cronJob, metaV1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return &cronJob.Status, nil
}

func GetCronJobStatus(ctx context.Context, clientset *kubernetes.Clientset, jobName string) (*batchV1.CronJobStatus, error) {
	cronJob, err := clientset.BatchV1().CronJobs(entity.Namespace).Get(ctx, jobName, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &cronJob.Status, nil
}
