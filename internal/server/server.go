package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/storskegg/storskegg.org/internal/serveSPA"
	"github.com/storskegg/storskegg.org/internal/status"
)

type Server interface {
	Serve() error
}

type server struct {
	cmd    *cobra.Command
	config *Config

	router *mux.Router
	srv    *http.Server
}

func New(config *Config, cmd *cobra.Command) Server {
	srv := &server{
		cmd:    cmd,
		config: config,
	}

	srv.router = mux.NewRouter()
	srv.router.Use(handlers.RecoveryHandler()) // Handle panics with grace
	srv.router.Use(handlers.CompressHandler)   // A little GZip goes a long way

	srv.router.HandleFunc("/sk/status", status.HandlerStatus)

	spa := serveSPA.SpaHandler{StaticPath: "dist", IndexPath: "index.html"}
	srv.router.PathPrefix("/").Handler(spa)

	originValidator := handlers.AllowedOriginValidator(func(origin string) bool {
		u, err := url.ParseRequestURI(origin)
		if err != nil {
			log.Println("ERR parsing origin: ", origin)
			return false
		}

		if u.Hostname() != "localhost" {
			log.Println("REJECTED: Origin hostname: ", u.Hostname())
			return false
		}

		switch u.Port() {
		case "1234", "3001":
			log.Println("ACCEPTED: Origin port: ", u.Port())
			return true
		default:
			log.Println("REJECTED: Origin port: ", u.Port())
			return false
		}
	})

	srv.srv = &http.Server{
		Handler: handlers.CORS(
			originValidator,
			handlers.AllowCredentials(),
			handlers.IgnoreOptions(),
			handlers.AllowedMethods([]string{
				http.MethodGet,
				http.MethodHead,
				http.MethodPost,
				http.MethodPut,
				http.MethodDelete,
			}),
		)(srv.router),
		Addr: ":3001",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return srv
}

func (s *server) Serve() error {
	go func() {
		log.Printf("Listening on %s", s.config.Addr)
		openIfNotOpen(fmt.Sprintf("http://127.0.0.1:%s", s.config.Addr))
		if err := s.srv.ListenAndServe(); err != nil {
			log.Print(err)
		}
	}()

	chanSig := make(chan os.Signal, 1)

	// We'll always attempt graceful shutdowns when quit via SIGINT (Ctrl+C), KILL, QUIT or TERM
	signal.Notify(chanSig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-chanSig

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.

	err := s.srv.Shutdown(ctx)
	if err != nil {
		log.Panic(err)
	}

	return nil
}

//////////////////////////////////////////////////////////////////////
// PRIVATE METHODS (Mostly for dev purposes)

func openIfNotOpen(url string) {
	_, err := os.Stat("./.browser.lock")
	if err == nil {
		log.Print("Browser lock detected")
		return
	}
	_, err = os.Create("./.browser.lock")
	if err != nil {
		log.Fatal(err)
		return
	}
	time.Sleep(100 * time.Millisecond)
	err = open(url)
	if err != nil {
		log.Fatal(err)
	}
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
