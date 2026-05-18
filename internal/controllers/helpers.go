package controllers

func computeDiscountPercent(original, discounted float64) *int {
	if discounted <= 0 || discounted >= original {
		return nil
	}
	percent := int(((original - discounted) / original) * 100)
	return &percent
}
