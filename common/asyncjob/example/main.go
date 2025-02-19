package main

import (
	"GoTodo/common/asyncjob"
	"context"
	"log"
	"time"
)

func main() {
	job1 := asyncjob.NewJob(func(ctx context.Context) error {
		time.Sleep(2 * time.Second)
		return nil
	}, asyncjob.WithName("job1"))

	job2 := asyncjob.NewJob(func(ctx context.Context) error {
		time.Sleep(3 * time.Second)
		return nil
	}, asyncjob.WithName("job2"))

	//if err := job1.Execute(context.Background()); err != nil {
	//	log.Println(err)
	//
	//	for {
	//		err := job1.Retry(context.Background())
	//
	//		if err != nil || job1.State() == asyncjob.StateRetryFailed {
	//			break
	//		}
	//	}
	//}

	if err := asyncjob.NewGroup(true, job1, job2).Run(context.Background()); err != nil {
		log.Print(err)
	}

}
