package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func Multiplicate(rtp float64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		RTP := 0.0
		gen := 25000
		for range gen {
			rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

			seqs := genSequence(rnd)
			sum1, sum0 := calculateRTP(rtp, seqs)

			RTP += sum1 / sum0
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{"result": RTP / float64(gen)})
	}
}

func calculateRTP(rtp float64, seqs []float64) (float64, float64) {
	sum1 := 0.0
	sum0 := 0.0

	for i := range seqs {
		sum0 += seqs[i]
	}
	sum1 += sum0 * rtp

	return sum1, sum0
}

func genSequence(rnd *rand.Rand) []float64 {
	min, max := 1.0, 10000.0
	length := int(((max - min) + min))

	sequences := make([]float64, length)
	for i := 0; i < length; i++ {
		sequences[i] = randomFloat(min, max, rnd)
	}
	return sequences
}

func randomFloat(min, max float64, rnd *rand.Rand) float64 {
	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()
	return rnd.Float64()*(max-min) + min
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
