package main

import (
  "net/http"
  // "time"
  "log"
  "github.com/gorilla/websocket"
  "github.com/rakyll/portmidi"
)

var upgrader = websocket.Upgrader{}

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
    println("whoops")
    return
  } else {
    // go func() {
      events := midi.Listen()
      println(events)
      for event := range events {
        println(event.Timestamp)
      }
    // }()
  }

  // in, err := portmidi.NewInputStream(deviceID, 1024)
  // if err != nil {
  //   println("error: %s", err)
  //   return
  // }
  // defer in.Close()

  // events, err := in.Read(1024)
  // if err != nil {
  //   // log.Fatal(err)
  //   println("ewps: %s", err)
  //   return
  // }
  // println(events)

  // ch := in.Listen()
  // event := <-ch

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "index.html")
  })

  http.HandleFunc("/v1/ws", func(w http.ResponseWriter, r *http.Request) {
    var conn, _ = upgrader.Upgrade(w, r, nil)
    go func(conn *websocket.Conn) {
      for {
        mType, msg, err := conn.ReadMessage()
        if err != nil {
          return
        }

        conn.WriteMessage(mType, msg)
      }
    }(conn)
  })

  // http.HandleFunc("/v2/ws", func(w http.ResponseWriter, r *http.Request) {
  //   var conn, _ = upgrader.Upgrade(w, r, nil)
  //   go func(conn *websocket.Conn) {
  //     for {
  //       _, msg, err := conn.ReadMessage()
  //       if err != nil {
  //         return
  //       }
  //       println(string(msg))
  //     }
  //   }(conn)
  // })

  // http.HandleFunc("/v3/ws", func(w http.ResponseWriter, r *http.Request) {
  //   var conn, _ = upgrader.Upgrade(w, r, nil)
  //   go func(conn *websocket.Conn) {
  //     ch := time.Tick(5 *time.Second)

  //     for range ch {
  //       conn.WriteJSON(myGsm{
  //         Username: "bmb",
  //         FirstName: "Ben",
  //         LastName: "Benjamin",
  //       })
  //     }

  //   }(conn)
  // })

  // http.HandleFunc("/v4/ws", func(w http.ResponseWriter, r *http.Request) {
  //   var conn, _ = upgrader.Upgrade(w, r, nil)
  //   go func(conn *websocket.Conn) {
  //     for {
  //       _, _, err := conn.ReadMessage()
  //       if err != nil {
  //         conn.Close()
  //       }
  //     }
  //   }(conn)

  //   go func(conn *websocket.Conn) {
  //     ch := time.Tick(5 *time.Second)

  //     for range ch {
  //       conn.WriteJSON(myGsm{
  //         Username: "bmb",
  //         FirstName: "Ben",
  //         LastName: "Benjamin",
  //       })
  //     }

  //   }(conn)
  // })

  http.ListenAndServe(":7000", nil)
}

type myGsm struct {
  Username string `json:"username"`
  FirstName string `json:"firstName"`
  LastName string `json:"lastName"`
}
