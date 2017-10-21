package main

import (
  "net/http"
  "time"
  "log"
  "github.com/gorilla/websocket"
  "github.com/rakyll/portmidi"
)

var upgrader = websocket.Upgrader{}

func main() {
  portmidi.Initialize()
  deviceID := portmidi.DefaultOutputDeviceID()
  out, err := portmidi.NewOutputStream(deviceID, 1024, 0)
  if err != nil {
      log.Fatal(err)
  }

  // note on events to play C major chord
  out.WriteShort(0x90, 60, 100)
  out.WriteShort(0x90, 64, 100)
  out.WriteShort(0x90, 67, 100)

  time.Sleep(2 * time.Second)

  // note off events
  out.WriteShort(0x80, 60, 100)
  out.WriteShort(0x80, 64, 100)
  out.WriteShort(0x80, 67, 100)

  out.Close()
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
