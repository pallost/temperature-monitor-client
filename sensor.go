package main

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"bytes"
	"encoding/json"
//	"os"
//	"os/signal"
	"io/ioutil"
	"net/url"

	"github.com/d2r2/go-dht"
)

type Measurement struct {
	Temperature float32
	Humidity    float32
	Date        int64
}

func getOutsideTemperature() float32 {
	forecaUrl := "https://www.foreca.com/lv"
	tampereId := "100634963"

	response, err := http.PostForm(forecaUrl, url.Values{
		"id": {tampereId},
	})

	if err != nil {
		log.Println(err)
		return -1
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Println(err)
		return -1
	}

	log.Println(string(body))

	return 1
}

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

	outsideTemperature := getOutsideTemperature()

	temperature, humidity, retried, err :=
		dht.ReadDHTxxWithRetry(sensorType, 4, false, 10)
	if err != nil {
		log.Println(err)
	}
	// print temperature and humidity
	fmt.Printf("Sensor = %v: Temperature = %v*C, Humidity = %v%%, Outside = %v (retried %d times)\n",
		sensorType, temperature, humidity, outsideTemperature, retried)

	// write to GCloud
	newMeasurement := Measurement{temperature, humidity, makeTimestamp()}
	uploadToGCloud(newMeasurement)
}

func main() {
	// calling also once in the beginning
	readAndSend()

/*	interval := time.Duration(30) * time.Minute
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
	}*/
}
