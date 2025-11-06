package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/arjunmoola/go-arxiv/client"
)

type App struct {
	arxivClient *client.Client
}
