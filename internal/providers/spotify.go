package providers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"tmix/internal/browser"
	"tmix/internal/config"
	"tmix/internal/player"

	"github.com/zmb3/spotify/v2"
	spot "github.com/zmb3/spotify/v2"
	spotauth "github.com/zmb3/spotify/v2/auth"
)

var (
	port        = "3200"
	redirectUrl = fmt.Sprintf("http://localhost:%s/callback", port)
	ch          = make(chan *spot.Client)
	state       = "ssdafkelonnvnljnsda"
	auth        *spotauth.Authenticator
)

type Spotify struct {
	AbstractMusicProvider
	client *spot.Client
	player *player.SpotifyPlayer
	cache  *config.TokenCache
	config config.SpotifyConfig
}

func NewSpotify(cfg config.SpotifyConfig) *Spotify {
	auth = spotauth.New(
		spotauth.WithRedirectURL(redirectUrl),
		spotauth.WithScopes(
			spotauth.ScopeUserReadCurrentlyPlaying,
			spotauth.ScopeUserReadPlaybackState,
			spotauth.ScopeUserModifyPlaybackState,
			spotauth.ScopePlaylistReadPrivate,
			spotauth.ScopePlaylistReadCollaborative,
		),
		spotauth.WithClientID(cfg.ClientId),
		spotauth.WithClientSecret(cfg.ClientSecret),
	)
	return &Spotify{
		config: cfg,
	}
}

func (s *Spotify) Name() string { return "Spotify" }

func (s *Spotify) Player() player.Player {
	if s.player == nil {
		s.player = player.New(s.client)
	}
	return s.player
}

func (s *Spotify) FetchPlaylists() []Playlist {
	pl, err := s.client.CurrentUsersPlaylists(context.Background())
	if err != nil {
		log.Fatal("Failed to get user playlists: ", err)
	}

	var playlists []Playlist
	for _, playlist := range pl.Playlists {
		playlists = append(playlists, Playlist{
			ProviderId:  playlist.ID.String(),
			Name:        playlist.Name,
			Description: playlist.Description,
		})
	}

	return playlists
}

func (s *Spotify) Queue() []Song {

	queue, err := s.client.GetQueue(context.Background())
	if err != nil {
		log.Fatal("Unable to get the current user's queue: ", err)
	}

	var songs []Song
	for _, song := range queue.Items {
		songs = append(songs, songFromSpotifySong(song))
	}

	return songs

}

func (s *Spotify) Login() {
	if s.cache.Enabled {
		tok := s.cache.FetchToken()
		if tok.Valid() {
			s.client = spot.New(auth.Client(context.Background(), tok))
			return
		}
	}
	s.newLogin()
}

func (s *Spotify) newLogin() {

	shutdownDone := &sync.WaitGroup{}
	shutdownDone.Add(1)

	shutDownSignal := make(chan int)

	go func() { handleCallback(shutDownSignal, shutdownDone) }()

	// start login
	go func() {
		url := auth.AuthURL(state)
		browser.Open(url)
		s.client = <-ch
		s.loggedIn = true
		tok, err := s.client.Token()
		if err != nil {
			log.Fatalf("Unable to retrieve token: %s", err)
		}
		s.cache.StoreToken(tok)
		shutDownSignal <- 1
	}()

	shutdownDone.Wait()
}

func handleCallback(i chan int, done *sync.WaitGroup) {
	http.HandleFunc("/callback", completeAuth)
	server := &http.Server{Addr: ":" + port, Handler: http.DefaultServeMux}

	go server.ListenAndServe()
	log.Print("Started callback listener. Waiting to shutdown...")
	<-i
	server.Shutdown(context.Background())
	log.Print("Succesfully shutdown callback server")
	done.Done()
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatalf("failed to get token: %v", err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s", st, state)
	}

	client := spot.New(auth.Client(r.Context(), tok))

	w.Header().Set("content-type", "text/html")
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}

func songFromSpotifySong(s spotify.FullTrack) Song {
	var artists []string
	for _, artist := range s.Artists {
		artists = append(artists, artist.Name)
	}
	return Song{
		ProviderId: s.ID.String(),
		Name:       s.Name,
		Artist:     strings.Join(artists, ","),
	}
}
