package main

import (
  "net/http"
  "log"
  "github.com/gorilla/websocket"
  "github.com/rakyll/portmidi"
)

var upgrader = websocket.Upgrader{}
var events = []portmidi.Event{}

func main() {
  err := portmidi.Initialize()
  if err != nil {
    log.Fatal(err)
  }
  defer portmidi.Terminate()

  // deviceID := portmidi.DefaultInputDeviceID(0)
  // var deviceID = *portmidi.DeviceID{0}
  // println("info:", portmidi.Info(1).IsInputAvailable)
  midi, err := portmidi.NewInputStream(1, 1024)
  if err != nil {
    log.Fatal("whoops couldn't forge midi stream:", err)
    return
  }

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "index.html")
  })

  http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
    var conn, _ = upgrader.Upgrade(w, r, nil)
    go func(conn *websocket.Conn) {
      events := midi.Listen()
      for event := range events {
        // if (event.Status < 208 || event.Status > 223) {
        //   // don't log Mono Key Pressure events...
          conn.WriteJSON(midiGsm{
            Status: event.Status,
            Data1: event.Data1,
            Data2: event.Data2,
          })
        // }
      }
    }(conn)
  })

  http.ListenAndServe(":7000", nil)
}

type midiGsm struct {
  Status int64 `json:"status"`
  Data1 int64 `json:"data1"`
  Data2 int64 `json:"data2"`
}
