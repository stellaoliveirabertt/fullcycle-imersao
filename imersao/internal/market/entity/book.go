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

	// Toda vez que formos receber uma ordem ele irá executar o for
	for order := range b.OrdersChannel {
		if order.OrderType == "BUY" {
			// Toda vez que a ordem for uma ordem tipo de compra ele vai adicionar na fila
			buyOrders.Push(order)
			// verifica se existe alguma ordem de venda e se o preço da ordem de venda é menor ou igual ao preço da ordem de compra quer dizer que aqui existe uma negociacao
			if sellOrders.Len() > 0 && sellOrders.Orders[0].Price <= order.Price {
				// vou remover essa ordem da minha fila para ela nao estar mais em negociação
				sellOrder := sellOrders.Pop().(*Order)

				// se a quantidade de ações da ordem de venda for maior que 0
				// quer dizer que ela ainda tem algo a negociar, se for menor que zero quer dizer que ela ja foi liquidada
				if sellOrder.PendingShares > 0 {
					// Crio uma transacao para fazer uma nova negociacao
					transaction := NewTransaction(sellOrder, order, sellOrder.Shares, sellOrder.Price)
					// adiciono a transacao na minha lista de transacoes para fazer calculos etc
					b.AddTransaction(transaction, b.WaitGroup)
					// uma vez que adicionamos ela na lista de transacoes, ela vai ser adicionada na lista de transacoes do book
					sellOrder.Transactions = append(sellOrder.Transactions, transaction) // aqui ela vai ser jogada para o kafka
					order.Transactions = append(order.Transactions, transaction)
					// os dados vao ser pegos em outra tred para serem processados
					b.OrdersChannelOutput <- sellOrder
					b.OrdersChannelOutput <- order

					// se a transacao nao terminou completamente ela volta pra fila
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

	// Preciso saber quem tem menos para subtrair de quem tem menos e adicionar as shares
	minShares := sellingShares

	if buyingShares < minShares {
		minShares = buyingShares
	}

	// Atualiza as shares do vendedor
	transaction.SellingOrder.Investor.UpdateAssertPosition(transaction.SellingOrder.Assert.ID, -minShares)
	// A quantidade de shares pendentes das ordens que estão sendo vendidas diminui
	transaction.SellingOrder.PendingShares -= minShares

	// Atualiza as shares do comprador
	transaction.BuyingOrder.Investor.UpdateAssertPosition(transaction.BuyingOrder.Assert.ID, minShares)
	transaction.BuyingOrder.PendingShares -= minShares

	transaction.Total = float64(transaction.Shares) * transaction.Price

	if transaction.BuyingOrder.PendingShares == 0 {
		transaction.BuyingOrder.Status = "CLOSED"
	}
	if transaction.SellingOrder.PendingShares == 0 {
		transaction.SellingOrder.Status = "CLOSED"
	}

	// Adiciona a transacao na lista de transacoes do book
	b.Transactions = append(b.Transactions, transaction)
}
