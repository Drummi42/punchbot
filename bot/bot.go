package bot

import (
	"fmt"
	"github.com/drummi42/punchbot/config"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// // Book Struct (Model)
// type Book struct {
// 	ID     string  `json:"id"`
// 	Isbn   string  `json:"isbn"`
// 	Title  string  `json:"title"`
// 	Author *Author `json:"author"`
// }

// // Author Struct
// type Author struct {
// 	ID        string `json:"id"`
// 	Firstname string `json:"firstname"`
// 	Lastname  string `json:"lastname"`
// }

// var books []Book

// //Get all books
// func getBooks(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(books)
// }

// func getBook(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)

// 	for _, item := range books {
// 		if item.ID == params["id"] {
// 			json.NewEncoder(w).Encode(item)
// 			return
// 		}
// 	}

// 	json.NewEncoder(w).Encode(&Book{})
// }
// func createBook(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var book Book
// 	_ = json.NewDecoder(r.Body).Decode(&book)
// 	book.ID = strconv.Itoa(rand.Intn(1000000)) // Mock ID
// 	books = append(books, book)

// 	json.NewEncoder(w).Encode(book)
// }
// func updateBook(w http.ResponseWriter, r *http.Request) {

// }
// func deleteBook(w http.ResponseWriter, r *http.Request) {

// }

// Duel model
type Duel struct {
	Player1 string
	Player2 string
}

var player1, player2 string = "", ""
var (
	duels []Duel
)

func Start() {
	// //Init Router
	// r := mux.NewRouter()

	// // Mock
	// books = append(books, Book{ID: "1", Isbn: "3232323", Title: "Book 1", Author: &Author{ID: "1", Firstname: "Jon", Lastname: "Doip"}})
	// books = append(books, Book{ID: "2", Isbn: "34545453", Title: "Book 2", Author: &Author{ID: "2", Firstname: "Von", Lastname: "Joip"}})

	// // Route Handlers / Endpoints
	// r.HandleFunc("/api/books", getBooks).Methods("GET")
	// r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	// r.HandleFunc("/api/books", createBook).Methods("POST")
	// r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	// r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// log.Fatal(http.ListenAndServe(":8000", r))

	bot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("error creating Discord session", err)
		return
	}

	bot.AddHandler(messageHandler)

	err = bot.Open()
	if err != nil {
		fmt.Println("error openening Discord connection", err)
		return
	}

	fmt.Println("Bot is running.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	err = bot.Close()

	if err != nil {
		fmt.Println("error closing Discord connection", err)
		return
	}
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	if !strings.HasPrefix(m.Content, config.BotPrefix) {
		return
	}

	if m.Content == "!fight" {
		p := m.Author.Username
		if player1 == "" {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Need one more dude!")
			player1 = p
			return
		}

		if player1 == p {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Another!")
		}
		if player2 == "" && player1 != p {
			player2 = p
			_, _ = s.ChannelMessageSend(m.ChannelID, "Fight: "+player1+" vs "+player2)

			rand1 := rand.Intn(100)
			rand2 := rand.Intn(100)

			if rand1 > rand2 {
				_, _ = s.ChannelMessageSend(m.ChannelID, player1+" WIN!")
			} else {
				_, _ = s.ChannelMessageSend(m.ChannelID, player2+" WIN!")
			}
			player1 = ""
			player2 = ""

			return
		}

	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}

func duel(p1 string, p2 string) {
	newD := Duel{p1, p2}
	duels = append(duels, newD)
}
