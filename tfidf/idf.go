package ktfidf

import (
	"log"
	"math"
)

func calculateIDF(token kparser, filter kfilter) map[string]float32 {
	count := 0
	result := make(map[string]int)
	for s, err := token.Next(); err == nil; s, err = token.Next() {
		if token.IsEOS(s) {
			count++
		}
	}
	log.Printf("idf: total %d lines", count)
	token.Reset()
	for s, err := token.Next(); err == nil; s, err = token.Next() {
		sentfilter := make(map[string]int)
		for ; !token.IsEOS(s) && err == nil; s, err = token.Next() {
			if !filter.Has(s) {
				if _, ok := sentfilter[s]; !ok {
					if _, ok := result[s]; ok {
						result[s]++
					} else {
						result[s] = 1
					}
					sentfilter[s] = 1
				}
			}
		}
	}
	resultd := make(map[string]float32)
	for i, v := range result {
		resultd[i] = float32(math.Log(float64(count) / float64(v)))
	}
	log.Printf("idf: total %d words", len(resultd))
	return resultd
}
