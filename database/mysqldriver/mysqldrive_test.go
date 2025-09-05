package mysqldriver

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func TestBuildDSN(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		port     string
		user     string
		password string
		dbname   string
		expected string
	}{
		{
			name:     "valid DSN with all fields",
			host:     "localhost",
			port:     "3306",
			user:     "root",
			password: "secret",
			dbname:   "myapp",
			expected: "root:secret@tcp(localhost:3306)/myapp?parseTime=true&loc=UTC",
		},
		{
			name:     "empty password",
			host:     "localhost",
			port:     "3306",
			user:     "root",
			password: "",
			dbname:   "myapp",
			expected: "root:@tcp(localhost:3306)/myapp?parseTime=true&loc=UTC",
		},
		{
			name:     "remote database",
			host:     "db.example.com",
			port:     "3306",
			user:     "appuser",
			password: "complexpass123",
			dbname:   "production_db",
			expected: "appuser:complexpass123@tcp(db.example.com:3306)/production_db?parseTime=true&loc=UTC",
		},
		{
			name:     "custom port",
			host:     "localhost",
			port:     "3307",
			user:     "dev",
			password: "devpass",
			dbname:   "testdb",
			expected: "dev:devpass@tcp(localhost:3307)/testdb?parseTime=true&loc=UTC",
		},
		{
			name:     "special characters in password",
			host:     "localhost",
			port:     "3306",
			user:     "user",
			password: "p@ssw0rd!",
			dbname:   "app",
			expected: "user:p@ssw0rd!@tcp(localhost:3306)/app?parseTime=true&loc=UTC",
		},
		{
			name:     "empty fields",
			host:     "",
			port:     "",
			user:     "",
			password: "",
			dbname:   "",
			expected: ":@tcp(:)/?parseTime=true&loc=UTC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildDSN(tt.host, tt.port, tt.user, tt.password, tt.dbname)
			if result != tt.expected {
				t.Errorf("BuildDSN() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestSession(t *testing.T) {
	ctx := context.Background()

	// Start MySQL container
	container, err := mysql.Run(ctx,
		"mysql:8.0",
		mysql.WithDatabase("testdb"),
		mysql.WithUsername("testuser"),
		mysql.WithPassword("testpass"),
	)

	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}
	defer container.Terminate(ctx)

	// Get connection string
	connStr, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Failed to get connection string: %v", err)
	}

	// Connect to database
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()

	// Finally, we test our session!
	session, err := Session(db)
	if err != nil {
		t.Errorf("Session() error = %v", err)
	}
	if session == nil {
		t.Error("Session() returned nil")
	}

}
