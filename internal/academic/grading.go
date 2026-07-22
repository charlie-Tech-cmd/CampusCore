package academic

// GradePoint converts a letter grade to grade points.
func GradePoint(grade string) float64 {
	switch grade {
	case "A":
		return 5.0
	case "B":
		return 4.0
	case "C":
		return 3.0
	case "D":
		return 2.0
	case "E":
		return 1.0
	default:
		return 0.0
	}
}
