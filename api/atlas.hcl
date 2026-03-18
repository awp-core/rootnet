env "local" {
  src = "file://internal/db/migrations/schema.hcl"
  url = "postgres://postgres:postgres@localhost:5432/cortexia?sslmode=disable"
  dev = "docker://postgres/16/dev"
}
