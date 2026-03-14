package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/NotKidding/olympus-server/internal/db"
	"github.com/NotKidding/olympus-server/internal/models"
)

func main() {
	// Parse command line flags
	agentID := flag.String("id", "", "The ID of the agent (e.g., archlinux-Nandu)")
	cmd := flag.String("cmd", "", "The command to execute")
	flag.Parse()

	if *agentID == "" || *cmd == "" {
		log.Fatal("Usage: go run cmd/admin/main.go --id <agent_id> --cmd <command>")
	}

	// Connect to DB
	db.InitDB()

	// Create a new pending task
	newTask := models.Task{
		AgentID: *agentID,
		Command: *cmd,
		Status:  "pending",
	}

	if err := db.DB.Create(&newTask).Error; err != nil {
		log.Fatalf("[-] Failed to inject task: %v", err)
	}

	fmt.Printf("[+] Task injected for %s: %s\n", *agentID, *cmd)
}
