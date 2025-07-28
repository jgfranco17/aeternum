package db

import (
	"log"
	"os"
	"sync"

	"github.com/jgfranco17/aeternum/api/environment"
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
			os.Getenv(environment.ENV_KEY_DB_URL),
			os.Getenv(environment.ENV_KEY_DB_KEY),
			&supabase.ClientOptions{},
		)
		if err != nil {
			log.Fatalf("Failed to initialize Supabase client: %v", err)
		}
		supaClient = client
	})
	return supaClient
}
