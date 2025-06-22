var creatureManager = creature.NewManager()

func main() {
    creatureManager.AddCreature(...)
    for range ticker.C {
        creatureManager.TickAll()
    }
}