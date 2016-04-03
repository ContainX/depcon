package marathon

import (
	"fmt"
	"strings"
)

func (c *MarathonClient) ListTasks() ([]*Task, error) {
	tasks := new(Tasks)
	resp := c.http.HttpGet(c.marathonUrl(API_TASKS), &tasks)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return tasks.Tasks, nil
}

func (c *MarathonClient) KillAppTasks(id string, host string, scale bool) ([]*Task, error) {
	tasks := new(Tasks)

	url := c.marathonUrl(API_APPS, id, PathTasks)
	if host != "" || scale {
		if host == "" {
			url = fmt.Sprintf("%s?scale=%v", url, scale)
		} else {
			url = fmt.Sprintf("%s?host=%s&scale=%v", url, host, scale)
		}
	}
	resp := c.http.HttpDelete(url, nil, tasks)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return tasks.Tasks, nil
}

func (c *MarathonClient) KillAppTask(taskId string, scale bool) (*Task, error) {
	task := new(Task)
	app := taskId[0:strings.LastIndex(taskId, ".")]
	url := c.marathonUrl(API_APPS, app, PathTasks, taskId)

	if scale {
		url = fmt.Sprintf("%s?scale=%v", url, scale)
	}
	resp := c.http.HttpDelete(url, nil, task)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return task, nil
}

func (c *MarathonClient) KillTasksAndScale(ids ...string) error {
	tasks := new(KillTasksScale)
	tasks.ids = ids

	url := c.marathonUrl(API_TASKS_DELETE)
	url = fmt.Sprintf("%s?scale=true", url)

	resp := c.http.HttpPost(url, tasks, &Tasks{})
	if resp.Error != nil {
		return resp.Error
	}
	return nil
}

func (c *MarathonClient) GetTasks(id string) ([]*Task, error) {
	tasks := new(Tasks)
	resp := c.http.HttpGet(c.marathonUrl(API_APPS, id, PathTasks), &tasks)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return tasks.Tasks, nil
}

func (c *MarathonClient) ListQueue() (*Queue, error) {
	q := new(Queue)
	resp := c.http.HttpGet(c.marathonUrl(API_QUEUE), &q)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return q, nil
}
