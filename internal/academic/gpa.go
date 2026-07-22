package academic

// CalculateGPA computes GPA using quality points divided by total credit units.
func CalculateGPA(totalQualityPoints float64, totalCreditUnits int) float64 {
	if totalCreditUnits == 0 {
		return 0
	}

	return totalQualityPoints / float64(totalCreditUnits)
}
