/*
 * Copyright (C) 2019 Zilliqa
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package workpool

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"time"
)

type WorkerPool struct {
	maxWorkers int64
	sem        *semaphore.Weighted
	takes      map[string]Task
	ids        []string
}

func NewWorkPool(max int64) *WorkerPool {
	m := semaphore.NewWeighted(max)
	ts := make(map[string]Task)
	ids := make([]string, 0)
	return &WorkerPool{
		maxWorkers: max,
		sem:        m,
		takes:      ts,
		ids:        ids,
	}
}

func (wp *WorkerPool) AddTask(task Task) {
	wp.ids = append(wp.ids, task.UUID())
	wp.takes[task.UUID()] = task
}

func (wp *WorkerPool) Top() Task {
	if len(wp.ids) == 0 {
		return nil
	}

	id := wp.ids[0]
	t := wp.takes[id]

	delete(wp.takes, id)
	wp.ids = wp.ids[1:]
	return t

}

func (wp *WorkerPool) Poll(ctx context.Context, quit <-chan struct{}) {
	for {
		select {
		case <-quit:
			fmt.Println("quit now..")
			break
		default:
			task := wp.Top()
			if task == nil {
				time.Sleep(time.Second * 3)
			} else {
				if err := wp.sem.Acquire(ctx, 1); err != nil {
					break
				}
				go func() {
					defer wp.sem.Release(1)
					task.Run()
				}()
			}
		}

	}
}
