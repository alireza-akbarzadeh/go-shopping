package jobs

import "log"

func (c *CronJobs) registerOrderJobs() {
	c.addJob("0 2 * * *", "update_overdue_orders", c.updateOverdueOrders)
}

func (c *CronJobs) updateOverdueOrders() {
	log.Println("Updating overdue orders...")
	if err := c.svc.Order.UpdateOverdueOrders(); err != nil {
		log.Printf("Overdue orders update failed: %v", err)
	}
}
