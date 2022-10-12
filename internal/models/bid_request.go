package models

import (
	"math"
	"net"
)

type BidRequestImp struct {
	ID        uint `json:"id"`
	MinWidth  uint `json:"minwidth"`
	MinHeight uint `json:"minheight"`
}

type BidRequestContext struct {
	IP        net.IP `json:"ip"`
	UserAgent string `json:"user_agent"`
}

type BidRequest struct {
	ID      string            `json:"id"`
	Imp     []BidRequestImp   `json:"imp"`
	Context BidRequestContext `json:"context"`
}

func (b *BidRequest) LoadFromPlacementRequest(p PlacementRequest) {
	b.ID = *p.ID
	b.Context = BidRequestContext(*p.Context)
	for _, tile := range *p.Tiles {
		b.Imp = append(b.Imp, BidRequestImp{
			ID:        tile.ID,
			MinWidth:  tile.Width,
			MinHeight: uint(math.Floor(float64(tile.Width) * tile.Ratio)),
		})
	}
}
