package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {
	// Открытие файла
	file, err := os.Open("config/getVideoContract.json")
	if err != nil {
		panic(err)
	}

	// Запись данных из JSON файла
	var jsonData map[string]interface{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&jsonData)
	if err != nil {
		panic(err)
	}

	body, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}

	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 5 * time.Second

	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		Body:   body,
		URL:    "http://localhost:8666/getVideoContract",
	})

	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Stress test") {
		metrics.Add(res)
	}

	metrics.Close()

	fmt.Printf("Errors        %v\n",
		metrics.Errors)
	fmt.Printf("Success       [1 - true, 2 - false]: [%v]\n",
		metrics.Success)
	fmt.Printf("Requests      [total, rate, throughput]: [%d, %v, %v]\n",
		metrics.Requests, rate, metrics.Throughput)
	fmt.Printf("Duration      [total, attack, wait]: [%s, %s, %s]\n",
		metrics.Duration, "атака", "ожидание")
	fmt.Printf("Latencies     [mean, 50, 95, 99, max]: [%s, %s, %s, %s, %s]\n",
		metrics.Latencies.Mean, metrics.Latencies.P50, metrics.Latencies.P95,
		metrics.Latencies.P99, metrics.Latencies.Max)
	fmt.Printf("Bytes In      [total, mean]: [%d, %v]\n",
		metrics.BytesIn.Total, metrics.BytesIn.Mean)
	fmt.Printf("Bytes Out     [total, mean]: [%d, %v]\n",
		metrics.BytesOut.Total, metrics.BytesOut.Mean)
}
