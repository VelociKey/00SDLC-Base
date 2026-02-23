package main

import "context"
import "dagger/olympusassurance/internal/dagger"

type OlympusAssurance struct{}

func (m *OlympusAssurance) HelloWorld(ctx context.Context) string { return "Hello from OlympusAssurance!" }

func main() { dagger.Serve() }
