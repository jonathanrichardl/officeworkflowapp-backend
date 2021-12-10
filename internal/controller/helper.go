package controller

import (
	"order-validation-v2/internal/controller/models"
	"order-validation-v2/internal/entity"
	"sync"
)

func (c *Controller) saveSubmission(submission models.Submission, wg *sync.WaitGroup) {
	image := models.DecodeSubmissionPayload(submission)
	_, err := c.submissions.NewSubmission(submission.Message, image, submission.TaskID)
	if err != nil {
		panic(err)
	}
	wg.Done()
}

func (c *Controller) updateTaskStatus(taskID string, wg *sync.WaitGroup, status uint8) {
	task, err := c.task.Get(taskID)
	if err != nil {
		panic(err)
	}
	task.SetStatus(status)
	err = c.task.UpdateTask(task)
	if err != nil {
		panic(err)
	}
	wg.Done()

}

func (c *Controller) processReviewForm(taskID string, wg *sync.WaitGroup, approved bool, forwardTo []string) {
	task, err := c.task.Get(taskID)
	var wg2 sync.WaitGroup
	if err != nil {
		panic(err)
	}
	task.NumOfReviewer += uint8(len(forwardTo))
	if approved {
		task.ReduceNumOfReviewer()
		if task.NumOfReviewer == 0 || task.NumOfReviewer == 255 {
			wg2.Add(1)
			go c.updateTaskStatus(taskID, &wg2, 2)
		}
		wg2.Add(1)
		go c.forward(taskID, forwardTo, &wg2)
	} else {
		go c.updateTaskStatus(taskID, &wg2, 1)
	}
	wg2.Wait()
	wg.Done()

}

func (c *Controller) forward(taskID string, adminIDs []string, wg *sync.WaitGroup) {
	for _, id := range adminIDs {
		err := c.task.AddReviewer(taskID, id)
		if err != nil {
			panic(err)
		}
	}
	wg.Done()

}
func (c *Controller) updateRequirementStatus(requirementID int, wg *sync.WaitGroup, status int8) {
	req, err := c.requirements.GetRequirementbyID(requirementID)
	if err != nil {
		panic(err)
	}
	req.SetStatus(status)
	c.requirements.UpdateRequirement(req)
	wg.Done()
}

func (c *Controller) deletePrerequisite(prerequisiteTaskID string, wg *sync.WaitGroup) {
	affectedTasks, err := c.task.RemovePrerequisite(prerequisiteTaskID)
	var wg2 sync.WaitGroup
	if err != nil {
		panic(err)
	}
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
		panic(err)
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
