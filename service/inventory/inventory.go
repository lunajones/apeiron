package inventory

import (
    "log"
    "sync"
)

var mu sync.Mutex

func Init() {
    log.Println("Inventory service initialized")
}

func TransferItem(from, to, itemID string, qty int) error {
    mu.Lock()
    defer mu.Unlock()
    log.Printf("TransferItem: %d of %s from %s to %s", qty, itemID, from, to)
    return nil
}
