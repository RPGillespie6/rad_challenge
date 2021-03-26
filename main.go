/*
    This file implements a simple HTTP API for posting and upvoting messages. The API is as follows:

    - GET /api/messages[?updated_after=<iso time string>]

    Get a list of all messages, including their id, text, number of upvotes, and last updated timestamp.
    Optionally provide an ISO8601 timestamp that only gives you new messages since that timestamp so you aren't pulling duplicate info you already have.

    - POST /api/message

    Create a new message

    - GET /api/vote?id=<uint>[&direction=down]

    Why did I choose GET here? Mainly because vanilla golang seems ill equipped to create APIs with path-based variables like: PUT /api/message/:id/vote

    If I had used a 3rd party http router from GitHub, I might have taken that approach, but I wanted to complete this challenge with vanilla golang first.

    However, there are a few other reasons I chose GET /api/vote:
    a) Ease of testing - I can test GET APIs using browser's URL bar and so I tend to favor GET for that reason
    b) Simplifies backend when you don't have to put in a bunch of conditional switching logic on method and body schema
    c) Concurrency - if two users call GET /api/vote at the same time there is no race and the DB will be updated correctly.
       Not so with a PUT /api/message approach where you provide vote count in body

    Note: Reddit, Hacker News have a dedicated GET /api/vote, while StackOverflow favors something more like POST /api/message/:id/vote

    Also note: anonymous commenting and voting is usually not a good idea. Anonymous voting makes vote manipulation super easy, and
    anonymous commenting almost certainly causes content moderation headaches as the user base grows.
*/
package main

import (
    // "fmt"
    "log"
    "time"
    "strconv"
    "net/http"
    "encoding/json"
)

// Data structures

type Message struct {
    ID            uint     `json:"id"`
    Text          string   `json:"text"`
    Upvotes       int      `json:"upvotes"`      // This should be an int if downvotes or unvotes are allowed
    LastUpdated   string   `json:"last_updated"` // RFC3339 timestamp
    // As this feature grows you probably want fields for who posted the message, who upvoted, etc.
}

// Our in-memory "database" is just a slice of messages and an incrementing id
// Initialize it so we get `[]` instead of `null` when serializing to json
var Messages []Message = make([]Message, 0)

// We don't want first id to be zero or we can't distinguish id 0 from field omitted
var MessageID uint = 1

// Start handlers

func setJsonHeaders(w http.ResponseWriter, status_code int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status_code)
}

func handleMessagePost(w http.ResponseWriter, r *http.Request) {
    var m Message

    err := json.NewDecoder(r.Body).Decode(&m)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if len(m.Text) == 0 {
        http.Error(w, "missing 'text' field", http.StatusBadRequest)
        return
    }

    m = Message{
        ID: MessageID,
        Text: m.Text,
        Upvotes: 0,
        LastUpdated: time.Now().Format(time.RFC3339),
    }

    // Update "database"
    MessageID++
    Messages = append(Messages, m)

    // Send response
    setJsonHeaders(w, http.StatusCreated)
    json.NewEncoder(w).Encode(m)
}

func handleVote(w http.ResponseWriter, r *http.Request) {
    // Grab 'id' query parameter
    query := r.URL.Query()

    id_str := query.Get("id")
    if len(id_str) == 0 {
        http.Error(w, "query parameter 'id' not found", http.StatusBadRequest)
        return
    }

    vote_delta := 1
    vote_dir := query.Get("direction")
    if vote_dir == "down" {
        vote_delta = -1
    }

    // Parse uint64 from id
    id_u64, err := strconv.ParseUint(id_str, 10, 32) // base 10, width 32 bits
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Down convert to uint so we can perform comparison
    ID := uint(id_u64)

    update_index := -1
    for i, m := range Messages {
        if ID == m.ID {
            update_index = i
            break
        }
    }

    if update_index == -1 {
        http.Error(w, "message index not found", http.StatusBadRequest)
        return
    }

    // Update message
    Messages[update_index].LastUpdated = time.Now().Format(time.RFC3339)
    Messages[update_index].Upvotes += vote_delta

    // Send response
    setJsonHeaders(w, http.StatusOK)
    json.NewEncoder(w).Encode(Messages[update_index])
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        handleMessagePost(w, r)
    } else {
        http.NotFound(w, r)
    }
}

func getMessages(w http.ResponseWriter, r *http.Request) {
    // Grab 'id' query parameter
    query := r.URL.Query()
    updated_after := query.Get("updated_after")
    if len(updated_after) == 0 {
        // If no updated_after parameter provided, send everything
        setJsonHeaders(w, http.StatusOK)
        json.NewEncoder(w).Encode(Messages)
        return;
    }

    timestamp, err := time.Parse(time.RFC3339, updated_after)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Only send messages updated after timestamp
    UpdatedMessages := make([]Message, 0)

    for _, m := range Messages {
        m_last_updated, _ := time.Parse(time.RFC3339, m.LastUpdated)
        if m_last_updated.After(timestamp) {
            UpdatedMessages = append(UpdatedMessages, m)
        }
    }

    setJsonHeaders(w, http.StatusOK)
    json.NewEncoder(w).Encode(UpdatedMessages)
}

func main() {
    http.HandleFunc("/api/vote", handleVote)
    http.HandleFunc("/api/message", handleMessage)
    http.HandleFunc("/api/messages", getMessages)

    // If it's not one of the APIs, serve static file(s)
    fs := http.FileServer(http.Dir("static/"))
    http.Handle("/", fs)

    log.Println("Listening on :3000...")
    log.Fatal(http.ListenAndServe(":3000", nil))
}