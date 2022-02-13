package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/jfyne/live"
)

func WithTemplateRenderer() live.HandlerConfig {
	return func(h live.Handler) error {
		h.HandleRender(func(ctx context.Context, data interface{}) (io.Reader, error) {
			t, err := template.ParseFiles("root.html", "buttons/view.html")
			if err != nil {
				log.Fatal(err)
			}
			var buf bytes.Buffer
			if err := t.Execute(&buf, data); err != nil {
				return nil, err
			}
			return &buf, nil
		})
		return nil
	}
}

const (
	inc = "inc"
	dec = "dec"
)

type counter struct {
	Value int
}

func newCounter(s live.Socket) *counter {
	c, ok := s.Assigns().(*counter)
	if !ok {
		return &counter{}
	}
	return c
}

func main() {
	h := live.NewHandler(WithTemplateRenderer())

	// Set the mount function for this handler.
	h.HandleMount(func(ctx context.Context, s live.Socket) (interface{}, error) {
		s.Upload("file")
		// This will initialise the counter if needed.
		return newCounter(s), nil
	})

	// Client side events.

	// Increment event. Each click will increment the count by one.
	h.HandleEvent(inc, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		// Get this sockets counter struct.
		c := newCounter(s)

		// Increment the value by one.
		c.Value += 1

		if err := s.Broadcast("newmessage", c.Value); err != nil {
			return c, fmt.Errorf("failed broadcasting new messaage: %w", err)
		}

		// Set the counter struct back to the socket data.
		return c, nil
	})

	// Decrement event. Each click will increment the count by one.
	h.HandleEvent(dec, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		// Get this sockets counter struct.
		c := newCounter(s)

		// Decrement the value by one.
		c.Value -= 1

		if err := s.Broadcast("newmessage", c.Value); err != nil {
			return c, fmt.Errorf("failed broadcasting new messaage: %w", err)
		}

		// Set the counter struct back to the socket data.
		return c, nil
	})

	h.HandleEvent("update", func(ctx context.Context, s live.Socket, p live.Params) (interface{}, error) {
		c := newCounter(s)

		s.UploadConsume("file", func(path string) string {
			dest := filepath.Join("public/uploads", filepath.Base(path))
			// fileutils.CopyFile(spath, dest)

			// dest = Path.join("priv/static/uploads", Path.basename(path))
			// File.cp!(path, dest)
			// Routes.static_path(socket, "/uploads/#{Path.basename(dest)}")

			return "/uploads/" + filepath.Base(dest)
		})

		// fmt.Printf("s: %v\n", s)

		// fmt.Printf("p: %v\n", p["myFile"])

		return c, nil
	})

	// h.HandleEvent("allow_upload", func(ctx context.Context, s live.Socket, p live.Params) (interface{}, error) {
	// 	c := newCounter(s)

	// 	// s.uploads[field]
	// 	fmt.Println(p)

	// 	return c, nil
	// })

	h.HandleSelf("newmessage", func(ctx context.Context, s live.Socket, data interface{}) (interface{}, error) {
		c := newCounter(s)

		// Here we don't append to messages as we don't want to use
		// loads of memory. `live-update="append"` handles the appending
		// of messages in the DOM.
		// m.Messages = []Message{NewMessage(data)}
		q, ok := data.(int)
		if !ok {
			q = 0
		}
		c.Value = q
		return c, nil
	})

	// Run the server.
	http.Handle("/", live.NewHttpHandler(live.NewCookieStore("session-name", []byte("weak-secret")), h))
	// http.HandleFunc("/live.js", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./vendor/web/browser/auto.js")
	// })
	// http.HandleFunc("/auto.js.map", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./vendor/web/browser/auto.js.map")
	// })
	// http.Handle("/live.js", http.FileServer(http.Dir("./vendor/web/browser")))
	// http.Handle("/auto.js.map", http.FileServer(http.Dir("./vendor/web/browser")))
	// http.Handle("/upload", media.Media{})

	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./public"))))

	fmt.Println("starting on :8080")
	http.ListenAndServe(":8080", nil)
}
