package controller

import (
	"order-validation-v2/internal/controller/models"
	"order-validation-v2/internal/entity"
	"sync"
)

func (c *Controller) saveSubmission(submission models.Submission, ch chan<- string, wg *sync.WaitGroup) {
	image := models.DecodeSubmissionPayload(submission)
	id, err := c.submissions.NewSubmission(submission.Message, image, submission.TaskID)
	if err != nil {
		c.logger.ErrorLogger.Println("Error while adding submission: ", err.Error())
		wg.Done()
		return
	}
	ch <- id
	wg.Done()
}

func (c *Controller) updateTaskStatus(task *entity.Task, wg *sync.WaitGroup) {
	task.Status = 1
	err := c.task.UpdateTask(task)
	if err != nil {
		c.logger.ErrorLogger.Println("Error while updating task: ", err.Error())
		wg.Done()
		return
	}
	wg.Done()

}

func (c *Controller) deletePrerequisite(prerequisiteTaskID string, wg *sync.WaitGroup) {
	affectedTasks, err := c.task.RemovePrerequisite(prerequisiteTaskID)
	if err != nil {
		c.logger.ErrorLogger.Println("Error while deleting"+prerequisiteTaskID+"From prerequisite: ", err.Error())
		wg.Done()
		return
	}
	var wg2 sync.WaitGroup
	for _, task := range affectedTasks {
		wg2.Add(1)
		go c.updateAffectedTasks(task, &wg2)
	}
	wg2.Wait()
	wg.Done()
}

func (c *Controller) updateAffectedTasks(affectedTask *entity.Task, wg *sync.WaitGroup) {
	affectedTask.ReducePrerequisite()
	if affectedTask.NumOfPrerequisite == 0 {
		affectedTask.Allow()
	}
	err := c.task.UpdateTask(affectedTask)
	if err != nil {
		c.logger.ErrorLogger.Println("Can't Update affected task : ", err.Error())
		return
	}
	wg.Done()
}

func (c *Controller) assignPrerequiste(task *entity.Task, assignedID map[string]string, wg *sync.WaitGroup) {
	for count, prerequiste := range task.Prerequisites {
		task.Prerequisites[count] = assignedID[prerequiste]
	}
	c.task.SaveTask(task)
	wg.Done()
}
