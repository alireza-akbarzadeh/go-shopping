package jobs

import (
	"github.com/alireza-akbarzadeh/shopping-platform/internal/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/internal/utils"
)

// Cart-related cron jobs
func (c *CronJobs) registerCartJobs() {
	c.addJob(constants.CronAbandonedCartCleanup, "abandoned_cart_cleanup", c.cleanAbandonedCarts)
}

func (c *CronJobs) cleanAbandonedCarts() {
	utils.Log.Info("Initiating cleanup of abandoned carts...")
	if err := c.svc.Cart.CleanAbandonedCarts(); err != nil {
		utils.Log.WithError(err).Error("Error cleaning abandoned carts")
	}
}
