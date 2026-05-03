package jobs

import (
	"github.com/alireza-akbarzadeh/shopping-platform/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
)

func (c *CronJobs) registerProductJobs() {
	c.addJob(constants.CronLowStockAlert, "low_stock_alert", c.checkLowStock)
	c.addJob(constants.CronSyncProductPrices, "sync_product_prices", c.syncProductPrices) // example
}

func (c *CronJobs) checkLowStock() {
	utils.Log.Info("Checking low stock products...")
	if err := c.svc.Product.CheckLowStockAndAlert(); err != nil {
		utils.Log.WithError(err).Error("Low stock check failed")
	}
}

func (c *CronJobs) syncProductPrices() {
	utils.Log.Info("Syncing product prices with external feed...")
	// external API call, update DB, etc.
}
