package jobs

import (
	"github.com/alireza-akbarzadeh/shopping-platform/internal/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/internal/utils"
)

func (c *CronJobs) registerOrderJobs() {
	c.addJob(constants.CronUpdateOverdueOrders, "update_overdue_orders", c.updateOverdueOrders)
}

func (c *CronJobs) updateOverdueOrders() {
	utils.Log.Info("Updating overdue orders...")
	if err := c.svc.Order.UpdateOverdueOrders(); err != nil {
		utils.Log.WithError(err).Error("Overdue orders update failed")
	}
}
