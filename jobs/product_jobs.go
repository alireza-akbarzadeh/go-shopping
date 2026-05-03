package jobs

import "log"

func (c *CronJobs) registerProductJobs() {
	c.addJob("0 9 * * *", "low_stock_alert", c.checkLowStock)
	c.addJob("0 1 * * *", "sync_product_prices", c.syncProductPrices) // example
}

func (c *CronJobs) checkLowStock() {
	log.Println("Checking low stock products...")
	if err := c.svc.Product.CheckLowStockAndAlert(); err != nil {
		log.Printf("Low stock check failed: %v", err)
	}
}

func (c *CronJobs) syncProductPrices() {
	log.Println("Syncing product prices with external feed...")
	// external API call, update DB, etc.
}
