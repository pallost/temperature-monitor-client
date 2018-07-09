package main

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"bytes"
	"encoding/json"
	"os"
	"os/signal"

	"github.com/d2r2/go-dht"
)

type Measurement struct {
	Temperature float32
	Humidity    float32
	Date        int64
}

//func getOutsideTemperature() {
//	const getUrl = "http://data.fmi.fi/fmi-apikey/d95d22a1-1da2-44d4-abfc-4f9ecb889700/wfs?request=getFeature&storedquery_id=fmi::observations::weather::cities::timevaluepair&parameters=Temperature"
//
//	_, err := http.Get(postUrl, "application/json", bytes.NewBuffer(jsonValue))
//	if err != nil {
//		log.Println(err)
//	}
//}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func uploadToGCloud(meas Measurement) {
	const postUrl = "https://guestbook-sample-179506.appspot.com/add"

	jsonValue, _ := json.Marshal(meas)

	_, err := http.Post(postUrl, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Println(err)
	}
}

func readAndSend() {
	// Read DHT11 sensor data from pin 4, retrying 10 times in case of failure.
	// You may enable "boost GPIO performance" parameter, if your device is old
	// as Raspberry PI 1 (this will require root privileges). You can switch off
	// "boost GPIO performance" parameter for old devices, but it may increase
	// retry attempts. Play with this parameter.
	sensorType := dht.DHT22

	temperature, humidity, retried, err :=
		dht.ReadDHTxxWithRetry(sensorType, 4, false, 10)
	if err != nil {
		log.Println(err)
	}
	// print temperature and humidity
	fmt.Printf("Sensor = %v: Temperature = %v*C, Humidity = %v%% (retried %d times)\n",
		sensorType, temperature, humidity, retried)

	// write to GCloud
	newMeasurement := Measurement{temperature, humidity, makeTimestamp()}
	uploadToGCloud(newMeasurement)
}

func main() {
	// calling also once in the beginning
	readAndSend()

	interval := time.Duration(30) * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	for {
		select {
		case <-signalCh:
			fmt.Println("Done!")
			return
		case <-ticker.C:
			readAndSend()
		}
	}
}
