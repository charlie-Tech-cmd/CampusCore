package academic

// CalculateQualityPoints returns the total quality points
// accumulated from all courses.
func CalculateQualityPoints(gradePoint float64, creditUnits int) float64 {
	return gradePoint * float64(creditUnits)
}
