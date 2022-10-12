package models

type BidResponseImp struct {
	ID     uint    `json:"id"`
	Width  uint    `json:"width"`
	Height uint    `json:"height"`
	Title  string  `json:"title"`
	URL    string  `json:"URL"`
	Price  float64 `json:"price"`
}

type BidResponse struct {
	ID  string           `json:"id"`
	Imp []BidResponseImp `json:"imp"`
}

func (b *BidResponse) IsEmpty() bool {
	return b.ID == "" || len(b.Imp) == 0
}

func (b *BidResponse) LeaveOnlyMaxProfitImps() {
	if len(b.Imp) == 0 {
		return
	}

	bidResponseImpMap := make(map[uint]BidResponseImp)
	for _, imp := range b.Imp {
		if _, k := bidResponseImpMap[imp.ID]; !k || (k && bidResponseImpMap[imp.ID].Price < imp.Price) {
			bidResponseImpMap[imp.ID] = imp
		}
	}

	b.Imp = []BidResponseImp{}
	for _, imp := range bidResponseImpMap {
		b.Imp = append(b.Imp, imp)
	}
}
