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

	go bot.Start()
	go log.Fatal(http.ListenAndServe(":"+p, r))
}

func botAwake(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode("{success: true}")
}
