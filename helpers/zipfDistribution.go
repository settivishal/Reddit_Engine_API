package helpers

// ZipfDistribution generates a Zipf distribution
func ZipfDistribution(n int, alpha float64) []int {
	z := make([]int, n)
	for i := 0; i < n; i++ {
		z[i] = int(float64(i+1) * alpha)
	}
	return z
}
