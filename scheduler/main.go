package main

import (
	"context"
	"fmt"
	"log"

	"github.com/copito/gocronscheduler/scheduler"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	ctx := context.Background()
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %s", err.Error())
	}

	// Start Watcher
	go scheduler.WatchCronJobs(ctx, clientSet)

	// Create schedule
	jobName := "example-cronjob"
	schedule := "*/5 * * * *"
	containerImage := "busybox"

	cronJob := scheduler.GenerateCronJobDefinition(jobName, containerImage, schedule)

	err = scheduler.CreateCronJob(ctx, clientSet, cronJob)
	if err != nil {
		log.Fatalf("Error creating cron job: %s", err.Error())
	}

	fmt.Println("Cron job created successfully")

	select {}
}
