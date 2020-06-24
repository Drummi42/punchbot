package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/drummi42/punchbot/config"
	"math/rand"
	"strings"
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
	Punisher      string
	MentionPlayer string
}

var (
	duels []Duel
	bot   *discordgo.Session
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
}

func Close() {
	err := bot.Close()

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

	if strings.Contains(m.Content, "duel") {
		msg := InitDuel(m)
		_, _ = s.ChannelMessageSend(m.ChannelID, msg)
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == config.BotPrefix+" ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}

// command duel <action> [@player]
func InitDuel(m *discordgo.MessageCreate) string {
	u := m.Author.Username
	mentions := m.Mentions

	if strings.Contains(m.Content, "list") {
		var list string
		for _, duel := range duels {
			if duel.MentionPlayer == u {
				list += duel.Punisher + " "
				// для меншинов нужен userID <@userId>
			}
		}
		if list == "" {
			return u + ", тебя не вызывали на дуэль "
		}
		return u + ", вот список игроков, которые ожидают твоего ответа: " + list
	}
	if strings.Contains(m.Content, "revoke") {
		// отмена дуэли
		for duelIndex, duel := range duels {
			if duel.Punisher == u {
				mention := duel.MentionPlayer
				DeleteDuel(duelIndex)
				// для меншинов нужен userID <@userId>
				return u + " отменил дуэль с " + mention + "!"
			}
		}
		return u + ", ты никого не вызывал на дуэль!"
	}

	// дальше уже команды требующие упоминания игрока
	if len(mentions) == 0 {
		return "Syntax command is: duel <?action> @mentionPlayer"
	}
	mentionU := mentions[0].Username
	if u == mentionU {
		return "Решил поиграть сам с собой?"
	}

	if strings.Contains(m.Content, "accept") {
		// принятие дуэли
		for duelIndex, duel := range duels {
			if duel.MentionPlayer == u {
				winner := RunDuel(duelIndex)
				// для меншинов нужен userID <@userId>
				return u + " принял бой с " + duel.Punisher + "! Победитель: " + winner
			}
		}
		return u + ", " + mentionU + " не вызывал тебя на дуэль!"
	} else {
		// вызыв на дуэль
		for _, duel := range duels {
			if duel.Punisher == u {
				return u + ", ты уже вызвал на дуэль " + duel.MentionPlayer + ". Дождись ответа!"
			}
		}
		duels = append(duels, Duel{u, mentionU})
		return u + " вызвал на дуэль " + mentionU + "!"
	}
}

func RunDuel(duelIndex int) string {
	var winner string

	rand1 := rand.Intn(100)
	rand2 := rand.Intn(100)
	if rand1 > rand2 {
		winner = duels[duelIndex].Punisher
	} else {
		winner = duels[duelIndex].MentionPlayer
	}
	DeleteDuel(duelIndex)
	return winner
}

func DeleteDuel(duelIndex int) {
	duels = append(duels[:duelIndex], duels[duelIndex+1:]...)
}
