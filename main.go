package main

import (
	"encoding/json"
	"fmt"
	"github.com/drummi42/punchbot/bot"
	"github.com/drummi42/punchbot/config"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {

	// //Init Router
	r := mux.NewRouter()
	p := os.Getenv("PORT")
	// Route Handlers / Endpoints
	r.HandleFunc("/heroku/awake", botAwake).Methods("GET")

	err := config.ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()
	log.Fatal(http.ListenAndServe(":"+p, r))

	//sc := make(chan os.Signal, 1)
	//signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	//<-sc
	//
	//// Cleanly close down the Discord session.
	//bot.Close()
}

func botAwake(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode("{success: true}")
}
