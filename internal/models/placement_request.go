package models

import (
	"net"
)

type PlacementRequestTile struct {
	ID    uint    `json:"id"`
	Width uint    `json:"width"`
	Ratio float64 `json:"ratio"`
}

type PlacementRequestContext struct {
	IP        net.IP `json:"ip"`
	UserAgent string `json:"user_agent"`
}

type PlacementRequest struct {
	ID      *string                  `json:"id"`
	Tiles   *[]PlacementRequestTile  `json:"tiles"`
	Context *PlacementRequestContext `json:"context"`
}

func (p *PlacementRequest) Validate() error {
	if p.IsAnyFieldNotExist() {
		return ErrPlacementRequestWrongSchema
	}
	if p.IsTilesEmpty() {
		return ErrPlacementRequestEmptyTiles
	}
	if p.IsAnyFieldEmpty() {
		return ErrPlacementRequestEmptyField
	}
	return nil
}

func (p *PlacementRequest) IsAnyFieldNotExist() bool {
	if p.ID == nil || p.Tiles == nil || p.Context == nil {
		return true
	}
	return false
}

func (p *PlacementRequest) IsTilesEmpty() bool {
	return len(*p.Tiles) == 0
}

func (p *PlacementRequest) IsAnyFieldEmpty() bool {
	if *p.ID == "" {
		return true
	}
	for _, t := range *p.Tiles {
		if t.ID == 0 || t.Width == 0 || t.Ratio == 0 {
			return true
		}
	}
	if string(p.Context.IP) == "" || p.Context.UserAgent == "" {
		return true
	}
	return false
}
