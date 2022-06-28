package service

import (
	"encoding/json"
	"errors"
	"lila_games/internal"
	"lila_games/internal/models"
	"net/http"
	"strconv"
	"strings"
)

// Service represents the game service
type Service struct {
	store internal.Store
}

// New creates new instance of game service
func New(store internal.Store) *Service {
	return &Service{store}
}

// RootHandler handles the service intro
func (l *Service) RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`
	<html>
	<head>
		<title>Game Service</title>
	</head>
	<body>
		<h1>Welcome to game service! </h1>
		<h3>The Game Service is bootstapped to work with a game named "FAMOUS SHOOTER" and supports following 3 modes:</h3>
		<ol>
			<li>Battle Royal</li>
			<li>Team Deathmatch</li>
			<li>Capture the flag</li>
		</ol>
		<h3>Following are the services offered</h3>
		<ul>
			<li>Get game's top modes in a given area 		<b>GET    /top-modes?area=blr</b></li>
			<li>Add new player to the game in a given area 	<b>POST   /player?modeid=1&area=blr</b></li>
			<li>Remove a player from the game 				<b>DELETE /player?playerid=1</b></li>
		</ul>
	</body>
	</html>`))
}

// TopModesHandler handles retrieval of game's top modes
func (l *Service) TopModesHandler(w http.ResponseWriter, r *http.Request) {
	// onyl GET method are allowed
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	gameid := strings.TrimSpace(r.URL.Query().Get("gameid"))
	var gameId int = 1 // default game id
	if len(gameid) > 0 {
		gameId, _ = strconv.Atoi(gameid)
	}
	area := strings.TrimSpace(r.URL.Query().Get("area"))
	if len(area) == 0 {
		http.Error(w, "area param missing", http.StatusBadRequest)
		return
	} else if len(area) != 3 {
		http.Error(w, "invalid area code", http.StatusBadRequest)
		return
	}
	modes, err := l.store.GetTopModeByArea(r.Context(), gameId, area)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(modes)
}

// PlayerHandler handles addition or removal of player
func (l *Service) PlayerHandler(w http.ResponseWriter, r *http.Request) {
	// only POST & DELETE allowed
	switch r.Method {
	case http.MethodPost:
		var (
			gameId int = 1 // default game id
			modeId int
			area   string
		)
		gameid := strings.TrimSpace(r.URL.Query().Get("gameid"))
		if len(gameid) > 0 {
			id, err := strconv.Atoi(gameid)
			if err != nil {
				http.Error(w, "invalid gameid param", http.StatusBadRequest)
				return
			}
			gameId = id
		}
		modeid := strings.TrimSpace(r.URL.Query().Get("modeid"))
		if len(modeid) == 0 {
			http.Error(w, "modeid param missing", http.StatusBadRequest)
			return
		}
		if len(modeid) > 0 {
			id, err := strconv.Atoi(modeid)
			if err != nil {
				http.Error(w, "invalid gameid param", http.StatusBadRequest)
				return
			}
			modeId = id
		}
		area = strings.TrimSpace(r.URL.Query().Get("area"))
		if len(area) == 0 {
			http.Error(w, "area param missing", http.StatusBadRequest)
			return
		} else if len(area) > 3 {
			http.Error(w, "invalid area code", http.StatusBadRequest)
			return
		}
		playerId, err := l.store.AddPlayer(r.Context(), gameId, modeId, area)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			PlayerID int `json:"playerId"`
		}{
			PlayerID: playerId,
		})
		return
	case http.MethodDelete:
		playerid := strings.TrimSpace(r.URL.Query().Get("playerid"))
		if len(playerid) == 0 {
			http.Error(w, "playerid param missing", http.StatusBadRequest)
			return
		}
		playerId, err := strconv.Atoi(playerid)
		if err != nil {
			http.Error(w, "invalid playerid param", http.StatusBadRequest)
			return
		}
		err = l.store.RemovePlayer(r.Context(), playerId)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusInternalServerError)
}

// Close closes database connection
func (l *Service) Close() error {
	return l.store.Close()
}
