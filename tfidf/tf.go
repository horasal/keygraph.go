package ktfidf

import (
	"log"
	"math"
)

func calculateTF(token kparser, filter kfilter) map[string]int {
	result := make(map[string]int)
	for s, err := token.Next(); err == nil; s, err = token.Next() {
		if !filter.Has(s) && !token.IsEOS(s) {
			if _, ok := result[s]; ok {
				result[s]++
			} else {
				result[s] = 1
			}
		}
	}
	log.Printf("tf: total %d words", len(result))
	return result
}

func calculatePoisson(token kparser, filter kfilter) map[string]float32 {
	result := make(map[string]float32)
	for s, err := token.Next(); err == nil; s, err = token.Next() {
		if !filter.Has(s) && !token.IsEOS(s) {
			if _, ok := result[s]; ok {
				result[s] += 1
			} else {
				result[s] = 1
			}
		}
	}
	for i, v := range result {
		result[i] = float32(math.Log(1.0 - math.Exp(float64(-v/float32(len(result))))))
	}
	log.Printf("tf-Poisson: total %d words", len(result))
	return result
}
