package entity

import (
	"container/heap"
	"sync"
)

type Book struct {
	Order               []*Order
	Transactions        []*Transaction
	OrdersChannel       chan *Order
	OrdersChannelOutput chan *Order
	WaitGroup           *sync.WaitGroup
}

func NewBook(orderChannel chan *Order, orderChannelOutput chan *Order, waitGroup *sync.WaitGroup) *Book {
	return &Book{
		Order:               []*Order{},
		Transactions:        []*Transaction{},
		OrdersChannel:       orderChannel,
		OrdersChannelOutput: orderChannelOutput,
		WaitGroup:           waitGroup,
	}
}

func (b *Book) Trade() {
	buyOrders := NewOrderQueue()
	sellOrders := NewOrderQueue()

	heap.Init(buyOrders)
	heap.Init(sellOrders)

	for order := range b.OrdersChannel {
		if order.OrderType == "BUY" {
			buyOrders.Push(order)

			if sellOrders.Len() > 0 && sellOrders.Orders[0].Price <= order.Price {
				sellOrder := sellOrders.Pop().(*Order)

				if sellOrder.PendingShares > 0 {
					transaction := NewTransaction(sellOrder, order, sellOrder.Shares, sellOrder.Price)
					b.AddTransaction(transaction, b.WaitGroup)
					sellOrder.Transactions = append(sellOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)
					b.OrdersChannelOutput <- sellOrder
					b.OrdersChannelOutput <- order

					if sellOrder.PendingShares > 0 {
						sellOrders.Push(sellOrders)
					}
				}
			}
		} else if order.OrderType == "SELL" {
			sellOrders.Push(order)

			if buyOrders.Len() > 0 && buyOrders.Orders[0].Price >= order.Price {
				buyOrder := buyOrders.Pop().(*Order)

				if buyOrder.PendingShares > 0 {
					transaction := NewTransaction(order, buyOrder, buyOrder.Shares, buyOrder.Price)
					b.AddTransaction(transaction, b.WaitGroup)
					buyOrder.Transactions = append(buyOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)
					b.OrdersChannelOutput <- buyOrder
					b.OrdersChannelOutput <- order

					if buyOrder.PendingShares > 0 {
						buyOrders.Push(buyOrder)
					}
				}
			}
		}
	}
}

func (b *Book) AddTransaction(transaction *Transaction, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	sellingShares := transaction.SellingOrder.PendingShares
	buyingShares := transaction.BuyingOrder.PendingShares
	minShares := sellingShares

	if buyingShares < minShares {
		minShares = buyingShares
	}

	transaction.SellingOrder.Investor.UpdateAssertPosition(transaction.SellingOrder.Assert.ID, -minShares)

	// TODO - Método para simplificar açoes
	transaction.SellingOrder.PendingShares -= minShares

	transaction.BuyingOrder.Investor.UpdateAssertPosition(transaction.BuyingOrder.Assert.ID, minShares)

	// TODO - Método para simplificar açoes
	transaction.BuyingOrder.PendingShares -= minShares

	// TODO - Calculate transaction total
	transaction.Total = float64(transaction.Shares) * transaction.Price

	// TODO - Método para simplificar açoes
	if transaction.BuyingOrder.PendingShares == 0 {
		transaction.BuyingOrder.Status = "CLOSED"
	}

	// TODO - Método para simplificar açoes
	if transaction.SellingOrder.PendingShares == 0 {
		transaction.SellingOrder.Status = "CLOSED"
	}

	b.Transactions = append(b.Transactions, transaction)
}
