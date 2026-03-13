package main

import (
	"context"
	"fmt"
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
	// 1. Initialize the Secure Persistence Bridge
	db.InitDB()

	server := &OlympusServer{}

	// 2. Start gRPC Server
	go func() {
		lis, err := net.Listen("tcp", ":9090")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterOlympusServiceServer(s, server)
		fmt.Println("[*] gRPC Management Server listening on :9090")
		s.Serve(lis)
	}()

	// 3. HTTP Listener (Unified with Database)
	http.HandleFunc("/checkin", func(w http.ResponseWriter, r *http.Request) {
		agentID := r.Header.Get("X-Agent-ID")
		if agentID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// 1. Extract the new Discovery Headers
		osVer := r.Header.Get("X-OS-Version")
		arch := r.Header.Get("X-Arch")

		// 2. Prepare the database record
		agent := models.Agent{
			ID:        agentID,
			Hostname:  r.RemoteAddr,
			OSVersion: osVer, // Save the recon data
			Arch:      arch,  // Save the recon data
			LastSeen:  time.Now(),
		}

		// 3. Upsert into Postgres
		if err := db.DB.Save(&agent).Error; err != nil {
			log.Printf("[-] Failed to save discovery data: %v", err)
			return
		}

		log.Printf("[+] Recon Received from %s: %s (%s)", agentID, osVer, arch)
		w.Write([]byte("OLYMPUS_ACK"))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCheckin(w http.ResponseWriter, r *http.Request) {
	agentID := r.Header.Get("X-Agent-ID")
	if agentID == "" {
		return
	}

	// Prepare the agent data
	agent := models.Agent{
		ID:       agentID,
		Hostname: agentID, // We'll split this later in System Discovery
		LastSeen: time.Now(),
	}

	// GORM "Save" performs an Upsert based on the Primary Key (ID)
	if err := db.DB.Save(&agent).Error; err != nil {
		log.Printf("[-] Failed to save beacon: %v", err)
		return
	}

	log.Printf("[+] Persistence: Updated %s", agentID)
	fmt.Fprintf(w, "OLYMPUS_ACK")
}
