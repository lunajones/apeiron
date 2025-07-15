package grpc

import (
	"log"
	"time"

	"net"

	"github.com/lunajones/apeiron/service/zone"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	UnimplementedCreatureSyncServer
}

func StartGRPCServer(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("[gRPC] Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	RegisterCreatureSyncServer(s, &grpcServer{})
	reflection.Register(s)

	log.Printf("[gRPC] Server listening on %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("[gRPC] Failed to serve: %v", err)
	}
}

// StreamCreatureUpdates envia snapshots periódicos
func (s *grpcServer) StreamCreatureUpdates(req *pb.SnapshotStreamRequest, stream CreatureSync_StreamCreatureUpdatesServer) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			log.Printf("[gRPC] Cliente fechou conexão de stream")
			return nil

		case <-ticker.C:
			var batch CreatureSnapshotBatch

			creatures := zone.Zones[0].Creatures // TEMP: só da primeira zona
			for _, c := range creatures {
				if !c.IsAlive() {
					continue
				}

				pos := c.GetPosition()

				snap := &CreatureSnapshot{
					Id:        c.Handle.String(),
					Name:      c.Name,
					Type:      c.GetPrimaryType(), // Ex: "Wolf", "Soldier"
					X:         float32(pos.X),
					Y:         float32(pos.Y),
					Z:         float32(pos.Z),
					Hp:        float32(c.HP),
					MaxHp:     float32(c.MaxHP),
					Animation: string(c.AnimationState),
					State:     string(c.AIState),
					Timestamp: time.Now().UnixMilli(),
				}
				batch.Snapshots = append(batch.Snapshots, snap)
			}

			if err := stream.Send(&batch); err != nil {
				log.Printf("[gRPC] Erro ao enviar snapshot: %v", err)
				return err
			}
		}
	}
}
