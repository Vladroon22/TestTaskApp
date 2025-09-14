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
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

		seqs := genSequence(rnd)
		mults := genMults(rtp, rnd, seqs)
		trans := transformed(mults, seqs)

		sum1, sum0 := calculateRTP(seqs, trans)

		RTP := sum1 / sum0
		writeJSON(w, http.StatusOK, map[string]interface{}{"result": RTP})
	}
}

func calculateRTP(seqs, trans []float64) (float64, float64) {
	sum1 := 0.0
	sum0 := 0.0

	for i := range seqs {
		sum0 += seqs[i]
		sum1 += trans[i]
	}

	return sum1, sum0
}

func transformed(mults, seqs []float64) []float64 {
	trans := make([]float64, len(seqs))
	for i, num := range seqs {
		mult := mults[i]

		Mult := (mult - 1.0) / (10000.0 - 1.0)
		switch {
		case num < 10:
			trans[i] = num * Mult * 0.5
		case num < 100:
			trans[i] = num * Mult * 0.7
		case num < 1000:
			trans[i] = num * Mult * 0.9
		default:
			trans[i] = num * Mult
		}
	}
	return trans
}

func genMults(rtp float64, rnd *rand.Rand, seqs []float64) []float64 {
	mults := make([]float64, len(seqs))
	for i := range seqs {
		mults[i] = genMultiplier(rtp, rnd)
	}
	return mults
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

func genMultiplier(rtp float64, rnd *rand.Rand) float64 {
	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	min := 1.0
	max := 10000.0

	base := rnd.Float64()*(max-min) + min

	// Корректируем множитель на основе RTP
	// Чем выше RTP, тем выше множители в среднем
	adjusted := base * (0.5 + rtp)

	if adjusted < min {
		return min
	} else if adjusted > max {
		return max
	}

	return adjusted
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
