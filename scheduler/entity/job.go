package entity

import (
	batchV1 "k8s.io/api/batch/v1"
)

type JobInfo struct {
	JobUUID string
	JobName string
	CronJob *batchV1.CronJob
}
