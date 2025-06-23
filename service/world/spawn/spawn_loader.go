package spawn

import (
	"encoding/json"
	"os"

	"github.com/lunajones/apeiron/lib/position"
)

type SpawnDefinition struct {
	TemplateID     int               `json:"TemplateID"`
	Position       position.Position `json:"Position"`
	Radius         float64           `json:"Radius"`
	Count          int               `json:"Count"`
	RespawnTimeSec int               `json:"RespawnTimeSec"`
}

func LoadSpawnsForZone(filePath string) ([]SpawnDefinition, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var spawns []SpawnDefinition
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&spawns); err != nil {
		return nil, err
	}

	return spawns, nil
}
