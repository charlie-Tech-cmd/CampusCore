package academic

// ClassifyDegree returns the degree classification from a CGPA.
func ClassifyDegree(cgpa float64) string {
	switch {
	case cgpa >= 4.50:
		return "First Class"

	case cgpa >= 3.50:
		return "Second Class Upper"

	case cgpa >= 2.40:
		return "Second Class Lower"

	case cgpa >= 1.50:
		return "Third Class"

	case cgpa >= 1.00:
		return "Pass"

	default:
		return "Fail"
	}
}
