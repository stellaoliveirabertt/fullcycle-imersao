package entity

type Investor struct {
	ID             string
	Name           string
	AssertPosition []*InvestorAssertPosition
}

type InvestorAssertPosition struct {
	AssertID string
	Shares   int
}

func NewInvestor(id string) *Investor {
	return &Investor{
		ID:             id,
		AssertPosition: []*InvestorAssertPosition{},
	}
}

func (i *Investor) AddAssertPosition(assertPosition *InvestorAssertPosition) {
	i.AssertPosition = append(i.AssertPosition, assertPosition)
}

func (i *Investor) UpdateAssertPosition(assertID string, qtdShares int) {
	assertPosition := i.GetAssertPosition(assertID)
	if assertPosition == nil {
		i.AssertPosition = append(i.AssertPosition, NewInvestorAssertPosition(assertID, qtdShares))
	} else {
		assertPosition.Shares += qtdShares
	}
}

func NewInvestorAssertPosition(assertID string, shares int) *InvestorAssertPosition {
	return &InvestorAssertPosition{
		AssertID: assertID,
		Shares:   shares,
	}
}

func (i *Investor) GetAssertPosition(assertID string) *InvestorAssertPosition {
	for _, assertPosition := range i.AssertPosition {
		if assertPosition.AssertID == assertID {
			return assertPosition
		}
	}
	return nil
}
