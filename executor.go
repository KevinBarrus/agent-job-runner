package main

import (
	"log"
	"time"
)

func (jh *JobHandler) runJob(id string) {
	// 1. 将 pending 转换为 running
	err := jh.store.UpdateStatus(id, JobRunning)
	if err != nil {
		log.Print("pending to running transition failed: ", err)
		return
	}

	// 2. 模拟执行 10 个步骤，每步约 1 秒
	for range 10 {
		time.Sleep(time.Second)
	}

	// 3. 执行完成后，将 running 转换为 succeeded
	err = jh.store.UpdateStatus(id, JobSucceeded)
	if err != nil {
		log.Print("running to succeeded transition failed: ", err)
		return
	}
}
