/*******************************************************************************
 * Copyright 2017.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package timerapp

import (
	"context"
	"github.com/winc-link/hummingbird/internal/dtos"
	resourceContainer "github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"

	//"github.com/winc-link/hummingbird/internal/edge/sharp/module/timer/db"

	//"github.com/winc-link/hummingbird/internal/pkg/timer/db"
	"github.com/winc-link/hummingbird/internal/pkg/timer/jobrunner"
	"github.com/winc-link/hummingbird/internal/pkg/timer/jobs"

	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"sort"
	"sync"
	"time"
)

type EdgeTimer struct {
	mutex  sync.Mutex
	logger logger.LoggingClient
	db     interfaces.DBClient

	//db     db.TimerDBClient
	// job id map
	jobMap map[string]struct{}
	// 任务链表 有序
	entries []*entry
	// 停止信号
	stop chan struct{}
	// 添加任务channel
	add chan *entry
	// 更新任务
	//update chan *jobs.UpdateJobStu
	// 删除任务 uuid
	rm chan string
	// 启动标志
	running  bool
	location *time.Location
	f        jobrunner.JobRunFunc
}

func NewCronTimer(ctx context.Context,
	f jobrunner.JobRunFunc, dic *di.Container) *EdgeTimer {
	dbClient := resourceContainer.DBClientFrom(dic.Get)
	l := container.LoggingClientFrom(dic.Get)
	et := &EdgeTimer{
		logger:   l,
		db:       dbClient,
		rm:       make(chan string),
		add:      make(chan *entry),
		entries:  nil,
		jobMap:   make(map[string]struct{}),
		stop:     make(chan struct{}),
		running:  false,
		location: time.Local,
		f:        f,
	}
	// restore
	et.restoreJobs()

	go et.run()
	return et
}

func (et *EdgeTimer) restoreJobs() {
	scenes, _, _ := et.db.SceneSearch(0, -1, dtos.SceneSearchQueryRequest{})

	if len(scenes) == 0 {
		return
	}

	for _, scene := range scenes {
		if len(scene.Conditions) > 0 && scene.Status == constants.SceneStart {
			if scene.Conditions[0].ConditionType == "timer" {
				job, err := scene.ToRuntimeJob()
				if err != nil {
					et.logger.Errorf("restore jobs runtime job err %v", err.Error())
					continue
				}
				err = et.AddJobToRunQueue(job)
				if err != nil {
					et.logger.Errorf("restore jobs add job to queue err %v", err.Error())
				}
			}
		}
	}
	return
}

func (et *EdgeTimer) Stop() {
	et.mutex.Lock()
	defer et.mutex.Unlock()
	if et.running {
		close(et.stop)
		et.running = false
	}
}

type (
	// 任务
	entry struct {
		JobID    string
		Schedule *jobs.JobSchedule
		Next     time.Time
		Prev     time.Time
	}
)

func (e entry) Valid() bool { return e.JobID != "" }

type byTime []*entry

func (s byTime) Len() int      { return len(s) }
func (s byTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byTime) Less(i, j int) bool {
	if s[i].Next.IsZero() {
		return false
	}
	if s[j].Next.IsZero() {
		return true
	}
	return s[i].Next.Before(s[j].Next)
}

func (et *EdgeTimer) now() time.Time {
	return time.Now().In(et.location)
}

func (et *EdgeTimer) run() {
	et.mutex.Lock()
	if et.running {
		et.mutex.Unlock()
		return
	}
	et.running = true
	et.mutex.Unlock()
	et.logger.Info("edge timer started...")
	now := et.now()
	for _, entry := range et.entries {
		if next, b := entry.Schedule.Next(now); !b {
			entry.Next = next
		}
	}
	var timer = time.NewTimer(100000 * time.Hour)
	for {
		// Determine the next entry to run.
		sort.Sort(byTime(et.entries))

		if len(et.entries) == 0 || et.entries[0].Next.IsZero() {
			// If there are no entries yet, just sleep - it still handles new entries
			// and stop requests.
			timer.Reset(100000 * time.Hour)
		} else {
			et.logger.Debugf("next wake time: %+v with jobID: %s", et.entries[0].Next, et.entries[0].JobID)
			timer.Reset(et.entries[0].Next.Sub(now))
		}

		select {
		case now = <-timer.C:
			timer.Stop()
			now = now.In(et.location)
			et.logger.Infof("wake now: %+v with jobID: %s", now, et.entries[0].JobID)
			var ()
			for i, e := range et.entries {
				if e.Next.After(now) || e.Next.IsZero() {
					break
				}
				// async call
				go et.f(e.JobID, *et.entries[i].Schedule)

				//times := e.Schedule.ScheduleAdd1()

				if false {
					//finished = append(finished, i)
				} else {
					e.Prev = e.Next
					if next, b := e.Schedule.Next(now); !b {
						e.Next = next
						et.logger.Infof("run now: %+v, entry: jobId: %s, jobName: %s, next: %+v", now, e.JobID, e.Schedule.JobName, e.Next)
					}
					//}
				}
			}

		case newEntry := <-et.add:
			timer.Stop()
			now = et.now()
			if next, b := newEntry.Schedule.Next(now); !b {
				newEntry.Next = next
				et.entries = append(et.entries, newEntry)
				et.logger.Infof("added job now: %+v, next: %+v", now, newEntry.Next)
			}
			et.logger.Infof("added job: %v, now: %+v, next: %+v", newEntry.JobID, now, newEntry.Next)
		case entryID := <-et.rm:
			timer.Stop()
			now = et.now()
			et.removeEntry(entryID)
		case <-et.stop:
			timer.Stop()
			et.logger.Info("tedge timer stopped...")
			return
		}
	}
}

func (et *EdgeTimer) schedule(schedule *jobs.JobSchedule) {
	et.mutex.Lock()
	defer et.mutex.Unlock()
	entry := &entry{
		JobID:    schedule.GetJobId(),
		Schedule: schedule,
	}
	if !et.running {
		et.entries = append(et.entries, entry)
	} else {
		et.add <- entry
	}
}

func (et *EdgeTimer) remove(id string) {
	if et.running {
		et.rm <- id
	} else {
		et.removeEntry(id)
	}
}

func (et *EdgeTimer) removeEntry(id string) {
	var b bool
	et.mutex.Lock()
	defer et.mutex.Unlock()
	for i, e := range et.entries {
		if e.JobID == id {
			et.entries[i], et.entries[len(et.entries)-1] = et.entries[len(et.entries)-1], et.entries[i]
			b = true
			break
		}
	}
	if b {
		et.entries[len(et.entries)-1] = nil
		et.entries = et.entries[:len(et.entries)-1]
		delete(et.jobMap, id)
		et.logger.Debugf("entry length: %d, deleted job id: %s", len(et.entries), id)
	} else {
		et.logger.Warnf("unknown jobs,id: %s", id)
	}
}

func (et *EdgeTimer) DeleteJob(id string) {
	et.remove(id)
}

func (et *EdgeTimer) AddJobToRunQueue(j *jobs.JobSchedule) error {
	if _, ok := et.jobMap[j.JobID]; ok {
		et.logger.Warnf("job is already in map: %s", j.JobID)
		return nil
	}

	if _, err := jobs.ParseStandard(j.TimeData.Expression); err != nil {
		return err
	}

	et.schedule(j)
	et.mutex.Lock()
	defer et.mutex.Unlock()
	et.jobMap[j.JobID] = struct{}{}
	return nil
}
