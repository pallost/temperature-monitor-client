package main

import (
    "fmt"
    "log"
    "time"
    "net/http"
    "bytes"
    "encoding/json"

    "github.com/d2r2/go-dht"
)

type Measurement struct {
    Temperature float32
    Humidity float32
    Date int64
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

func main() {
    // Read DHT11 sensor data from pin 4, retrying 10 times in case of failure.
    // You may enable "boost GPIO performance" parameter, if your device is old
    // as Raspberry PI 1 (this will require root privileges). You can switch off
    // "boost GPIO performance" parameter for old devices, but it may increase
    // retry attempts. Play with this parameter.
    sensorType := dht.DHT22
    for {
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

        time.Sleep(time.Duration(30)*time.Minute)
    }
}
