package models

type PlacementResponseImp struct {
	ID     uint   `json:"id"`
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
	Title  string `json:"title"`
	URL    string `json:"url"`
}

type PlacementResponse struct {
	ID  string                 `json:"id"`
	Imp []PlacementResponseImp `json:"tiles"`
}

func (p *PlacementResponse) IsEmpty() bool {
	return p.ID == "" || len(p.Imp) == 0
}

func (p *PlacementResponse) LoadFromPlacementRequestAndBidResponse(req PlacementRequest, res BidResponse) {
	p.ID = *req.ID
	p.Imp = []PlacementResponseImp{}
	if len(res.Imp) == 0 {
		return
	}

	bidResponseImpMap := make(map[uint]BidResponseImp)
	for _, imp := range res.Imp {
		if _, k := bidResponseImpMap[imp.ID]; !k || (k && bidResponseImpMap[imp.ID].Price < imp.Price) {
			bidResponseImpMap[imp.ID] = imp
		}
	}

	for _, tile := range *req.Tiles {
		if _, k := bidResponseImpMap[tile.ID]; k {
			imp := bidResponseImpMap[tile.ID]
			p.Imp = append(p.Imp, PlacementResponseImp{
				ID:     imp.ID,
				Width:  imp.Width,
				Height: imp.Height,
				Title:  imp.Title,
				URL:    imp.URL,
			})
		}
	}
}
