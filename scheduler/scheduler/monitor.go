package scheduler

import (
	"context"
	"fmt"
	"log"

	"github.com/copito/gocronscheduler/entity"
	batchV1 "k8s.io/api/batch/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

func WatchCronJobs(ctx context.Context, clientSet *kubernetes.Clientset) {
	watchInterface, err := clientSet.BatchV1().CronJobs(entity.Namespace).Watch(ctx, metaV1.ListOptions{})
	if err != nil {
		log.Fatalf("Error watching cron jobs: %s", err.Error())
	}
	for event := range watchInterface.ResultChan() {
		cronJob, ok := event.Object.(*batchV1.CronJob)
		if !ok {
			log.Fatalf("Unexpected type")
		}
		jobName := cronJob.Name
		jobUUID := cronJob.Labels["job_uuid"]

		job := entity.JobInfo{
			JobUUID: jobUUID,
			JobName: jobName,
			CronJob: cronJob,
		}

		switch event.Type {
		case watch.Added:
			// If they are added then nothing should be done (feedback loop maybe)
			handleCronJobAddition(jobName, &job)
		case watch.Deleted:
			handleCronJobDeletion(jobName, &job)
		case watch.Modified:
			handleCronJobModification(jobName, &job)
		case watch.Error:
			// If there is an error - service should be notified and change state in db
			handleCronJobError(jobName, &job)
		case watch.Bookmark:
			// No idea what this event is
		}
	}
}

func handleCronJobAddition(jobName string, job *entity.JobInfo) {
	fmt.Printf("Cron job %s was added.\n", jobName)
	fmt.Println("TODO: Update database with this cron_job new status (adding uuid)", job.CronJob.Status)
}

func handleCronJobDeletion(jobName string, job *entity.JobInfo) {
	fmt.Printf("Cron job %s was deleted.\n", jobName)
	fmt.Println("TODO: Update database with this cron_job new status (removing uuid)", job.CronJob.Status)
}

func handleCronJobModification(jobName string, job *entity.JobInfo) {
	fmt.Printf("Cron job %s status changed: %+v\n", jobName, job.CronJob.Status)
}

func handleCronJobError(jobName string, job *entity.JobInfo) {
	fmt.Printf("Cron job %s was error.\n", jobName)
	fmt.Println("TODO: Update database with this cron_job new status", job.CronJob.Status)
}
