# Go backend for DIY 12V UPS

## Structure

```
├── build              // Build backend
├── Dockerfile
├── run                // Run interactivelly
├── rund               // Run as a daemon
├── src
│   ├── app
│   │   └── app.go
│   ├── config
│   │   └── config.go
│   ├── db
│   │   └── db.go      // Database implementation 
│   ├── handler
│   │   ├── helper.go
│   │   └── ups.go
│   ├── main.go
│   ├── model
│   │   └── volt.go
│   ├── monitor
│   │   └── monitor.go // Battery monitor
│   ├── reader
│   │   └── reader.go  // Store voltage to database
│   ├── server
│   │   └── server.go  // HTTP server
│   └── writer
│       └── writer.go
└── system_shutdown    // Host power off script on UPS low battery
```

## Build

```bash
docker build -t goupsbe .
```

## Run

```bash
docker run --rm -v $(pwd):/app/db -v /var/run/shutdown_signal:/shutdown_signal -p 3000:3000 --env LISTEN_PORT=3000 --env SERIAL_DEVICE=/dev/ttyUSB0 --device=/dev/ttyUSB0 goupsbe
```

## API

Voltage without frag (5022 -> 5.022V)

#### /volt
* `GET` : Get voltage

/volt example:
```JSON
{
  "Vcc": 5022,   // The controller voltage
  "Vin": 12453,  // Input from power supplier
  "Vcin": 12306, // Input on DC-DC boost converter
  "Vout": 12393, // Output from DC-DC boost converter (UPS voltage)
  "Vb1": 7509,   // Battery output voltage
  "Vb2": 3840,   // Battery output voltage: middle point
  "Created": "2020-04-26T12:33:27.062857348Z"
}
```

#### /hist/:page_number/:page_size
* `GET` : Get voltage history

/hist/1/1 example:
```JSON
{
  "Volt": [
    {
      "Vcc": 5022,
      "Vin": 12559,
      "Vcin": 12490,
      "Vout": 12409,
      "Vb1": 7384,
      "Vb2": 3687,
      "Created": "2020-04-26T12:45:00Z"
    }
  ],
  "PageNumber": 1,
  "PageSize": 1,
  "Records": 440
}
```
