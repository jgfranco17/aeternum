package db

import (
	"log"
	"os"
	"sync"

	supabase "github.com/supabase-community/supabase-go"
)

var (
	supaClient *supabase.Client
	once       sync.Once
)

// GetSupabaseClient returns a singleton Supabase client
func GetSupabaseClient() *supabase.Client {
	once.Do(func() {
		client, err := supabase.NewClient(
			os.Getenv("AETERNUM_DB_URL"),
			os.Getenv("AETERNUM_DB_KEY"),
			&supabase.ClientOptions{},
		)
		if err != nil {
			log.Fatalf("Failed to initialize Supabase client: %v", err)
		}
		supaClient = client
	})
	return supaClient
}
