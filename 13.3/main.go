package main

import (
	"errors"
	"fmt"
)

type Product struct {
	ID    int
	Name  string
	Price float64
}

type Cart struct {
	Items []CartItem
}

type CartItem struct {
	Product  Product
	Quantity int
}

func (c *Cart) AddProduct(p Product, qty int) {
	c.Items = append(c.Items, CartItem{Product: p, Quantity: qty})
}

func (c *Cart) GetTotal() float64 {
	total := 0.0
	for _, item := range c.Items {
		total += item.Product.Price * float64(item.Quantity)
	}
	return total
}

type PaymentMethod string

const (
	PaymentCard   PaymentMethod = "card"
	PaymentPayPal PaymentMethod = "paypal"
	PaymentCash   PaymentMethod = "cash_on_delivery"
)

type PromoCode struct {
	Code            string
	DiscountPercent float64
}

type Order struct {
	ID            int
	CustomerName  string
	Address       string
	Cart          Cart
	PaymentMethod PaymentMethod
	TotalAmount   float64
	Status        string
	Cancelled     bool
}

type NotificationService struct{}

func (ns *NotificationService) Notify(msg string) {
	fmt.Printf("Notification: %s\n", msg)
}

type OrderProcessor struct {
	NextOrderID int
	Notifier    *NotificationService
}

func NewOrderProcessor() *OrderProcessor {
	return &OrderProcessor{
		NextOrderID: 1,
		Notifier:    &NotificationService{},
	}
}

func (op *OrderProcessor) CreateCart() *Cart {
	return &Cart{}
}

func (op *OrderProcessor) CreateOrder(cart *Cart, name, address string, paymentMethod PaymentMethod) *Order {
	if len(cart.Items) == 0 {
		panic("Cart is empty")
	}
	order := &Order{
		ID:            op.NextOrderID,
		CustomerName:  name,
		Address:       address,
		Cart:          *cart,
		PaymentMethod: paymentMethod,
		Status:        "created",
		Cancelled:     false,
	}
	op.NextOrderID++
	return order
}

func (op *OrderProcessor) Pay(order *Order, promo *PromoCode) error {
	if order.Cancelled {
		return errors.New("order cancelled")
	}

	if !op.simulatePayment(order.PaymentMethod) {
		return errors.New("payment failed")
	}

	total := order.Cart.GetTotal()
	if promo != nil {
		discount := total * (promo.DiscountPercent / 100)
		total -= discount
		op.Notifier.Notify(fmt.Sprintf("Promo code %s applied. Discount: %.2f", promo.Code, discount))
	}

	order.TotalAmount = total
	order.Status = "paid"
	op.Notifier.Notify(fmt.Sprintf("Payment successful. Total: %.2f", total))
	return nil
}

func (op *OrderProcessor) simulatePayment(method PaymentMethod) bool {
	fmt.Printf("Processing payment via %s...\n", method)
	return true
}

func (op *OrderProcessor) ProcessAndShip(order *Order) error {
	if order.Status != "paid" {
		return errors.New("payment not confirmed")
	}
	op.Notifier.Notify("Order is being processed at the warehouse")
	op.Notifier.Notify(fmt.Sprintf("Order #%d shipped to address: %s", order.ID, order.Address))
	order.Status = "shipped"
	return nil
}

func (op *OrderProcessor) CancelOrder(order *Order) {
	if order.Status == "paid" || order.Status == "shipped" {
		fmt.Println("Cannot cancel paid order")
		return
	}
	order.Cancelled = true
	order.Status = "cancelled"
	op.Notifier.Notify("Order cancelled")
}

func main() {
	processor := NewOrderProcessor()

	phone := Product{ID: 1, Name: "Smartphone", Price: 50000}
	charger := Product{ID: 2, Name: "Charger", Price: 1500}

	cart := processor.CreateCart()
	cart.AddProduct(phone, 1)
	cart.AddProduct(charger, 2)
	fmt.Printf("Cart: %.2f RUB\n", cart.GetTotal())

	order := processor.CreateOrder(cart, "Ivan Petrov", "10 Lenin St", PaymentCard)

	promo := &PromoCode{Code: "SAVE10", DiscountPercent: 10}

	err := processor.Pay(order, promo)
	if err != nil {
		fmt.Println("Payment error:", err)
		processor.Pay(order, nil)
	}

	processor.ProcessAndShip(order)

	fmt.Println("\n--- Scenario: cancellation before payment ---")
	cart2 := processor.CreateCart()
	cart2.AddProduct(phone, 1)
	order2 := processor.CreateOrder(cart2, "Maria", "5 Pushkin St", PaymentCash)
	processor.CancelOrder(order2)

	fmt.Println("\n--- Scenario: cancellation attempt after payment ---")
	cart3 := processor.CreateCart()
	cart3.AddProduct(charger, 1)
	order3 := processor.CreateOrder(cart3, "Alexey", "1 Gagarin St", PaymentPayPal)
	processor.Pay(order3, nil)
	processor.CancelOrder(order3)
}
