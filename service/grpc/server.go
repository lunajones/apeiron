package grpc

import (
	"log"
	"net"
	"time"

	"github.com/lunajones/apeiron/service/zone"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	UnimplementedCreatureSyncServer
	UnimplementedNavMeshServiceServer
}

func StartGRPCServer(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("[gRPC] Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	RegisterCreatureSyncServer(s, &grpcServer{})
	RegisterNavMeshServiceServer(s, &grpcServer{})
	reflection.Register(s)

	log.Printf("[gRPC] Server listening on %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("[gRPC] Failed to serve: %v", err)
	}
}

func (s *grpcServer) StreamCreatureUpdates(req *SnapshotStreamRequest, stream CreatureSync_StreamCreatureUpdatesServer) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	previousAlive := make(map[string]bool)

	for {
		select {
		case <-stream.Context().Done():
			log.Printf("[gRPC] Cliente fechou conexão de stream")
			return nil

		case <-ticker.C:
			var batch CreatureSnapshotBatch
			currentAlive := make(map[string]bool)

			creatures := zone.Zones[0].Creatures
			for _, c := range creatures {
				id := c.Handle.String()
				if !c.IsAlive() {
					continue
				}

				pos := c.GetPosition()
				snap := &CreatureSnapshot{
					Id:        id,
					Name:      c.Name,
					Type:      c.GetPrimaryType(),
					X:         float32(pos.X),
					Y:         float32(pos.Y),
					Z:         float32(pos.Z),
					Hp:        float32(c.HP),
					MaxHp:     float32(c.MaxHP),
					Animation: string(c.AnimationState),
					State:     string(c.AIState),
					Timestamp: time.Now().UnixMilli(),
					FaceYaw:   float32(c.GetFacingDirection().YawDeg()),
					TorsoYaw:  float32(c.GetTorsoDirection().YawDeg()),
				}
				batch.Snapshots = append(batch.Snapshots, snap)
				currentAlive[id] = true
			}

			for id := range previousAlive {
				if !currentAlive[id] {
					batch.Despawns = append(batch.Despawns, &CreatureDespawn{Id: id})
				}
			}

			previousAlive = currentAlive

			if err := stream.Send(&batch); err != nil {
				log.Printf("[gRPC] Erro ao enviar snapshot: %v", err)
				return err
			}
		}
	}
}

func (s *grpcServer) StreamNavMeshUpdates(_ *NavMeshStreamRequest, stream NavMeshService_StreamNavMeshUpdatesServer) error {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			log.Printf("[gRPC] Cliente fechou conexão de stream de NavMesh")
			return nil

		case <-ticker.C:
			nav := zone.Zones[0].NavMesh
			if nav == nil {
				continue
			}

			var batch NavMeshSnapshotBatch
			for _, poly := range nav.Polygons {
				snap := &NavMeshSnapshot{
					Id:        int32(poly.ID),
					GridX:     int32(poly.GridX),
					GridZ:     int32(poly.GridZ),
					OffsetX:   float32(poly.OffsetX),
					OffsetZ:   float32(poly.OffsetZ),
					Y:         float32(poly.Y),
					Slope:     float32(poly.Slope),
					AreaType:  poly.AreaType,
					Neighbors: convertNeighbors(poly.Neighbors),
				}
				batch.Polygons = append(batch.Polygons, snap)
			}

			if err := stream.Send(&batch); err != nil {
				log.Printf("[gRPC] Erro ao enviar snapshot do NavMesh: %v", err)
				return err
			}
		}
	}
}

func convertNeighbors(in []int) []int32 {
	out := make([]int32, len(in))
	for i, val := range in {
		out[i] = int32(val)
	}
	return out
}
