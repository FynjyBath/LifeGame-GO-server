package service

import (
	"math/rand"
	"time"

	"github.com/fynjybath/LifeGame-GO-server/pkg/life"
)

// для хранения состояния
type LifeService struct {
	currentWorld *life.World
	nextWorld    *life.World
}

func New(height, width int) (*LifeService, error) {
	rand.NewSource(time.Now().UTC().UnixNano())

	currentWorld, err := life.NewWorld(height, width)
	if err != nil {
		return nil, err
	}
	// для упрощения примера хаотично заполним
	currentWorld.RandInit(40)

	newWorld, err := life.NewWorld(height, width)
	if err != nil {
		return nil, err
	}

	ls := LifeService{
		currentWorld: currentWorld,
		nextWorld:    newWorld,
	}

	return &ls, nil
}

// получение очередного состояния игры
func (ls *LifeService) NewState() *life.World {
	life.NextState(ls.currentWorld, ls.nextWorld)

	ls.currentWorld = ls.nextWorld

	return ls.currentWorld
}
