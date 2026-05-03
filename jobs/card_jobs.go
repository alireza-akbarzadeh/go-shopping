package jobs

import "log"

// Cart-related cron jobs
func (c *CronJobs) registerCartJobs() {
	c.addJob("@every 30m", "abandoned_cart_cleanup", c.cleanAbandonedCarts)
}

func (c *CronJobs) cleanAbandonedCarts() {
	log.Println("Initiating cleanup of abandoned carts...")
	if err := c.svc.Cart.CleanAbandonedCarts(); err != nil {
		log.Printf("Error cleaning abandoned carts: %v", err)
	}
}
