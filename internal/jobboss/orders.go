package jobboss

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// Order models a JobBOSS order payload.
type Order struct {
	UniqueID            int     `json:"uniqueID,omitempty"`
	OrderNumber         string  `json:"orderNumber,omitempty"`
	CustomerCode        string  `json:"customerCode,omitempty"`
	CustomerDescription string  `json:"customerDescription,omitempty"`
	Status              string  `json:"status,omitempty"`
	PONumber            string  `json:"PONumber,omitempty"`
	SalesID             string  `json:"salesID,omitempty"`
	TermsCode           string  `json:"termsCode,omitempty"`
	Territory           string  `json:"territory,omitempty"`
	CurrencyCode        string  `json:"currencyCode,omitempty"`
	DateEntered         string  `json:"dateEntered,omitempty"`
	LastModDate         string  `json:"lastModDate,omitempty"`
	OrderTotal          float64 `json:"orderTotal,omitempty"`
	OrderTotalForeign   float64 `json:"orderTotalForeign,omitempty"`
}

// OrderCreate is a minimal create payload for /api/v1/orders.
type OrderCreate struct {
	OrderNumber     string `json:"orderNumber,omitempty"`
	CustomerCode    string `json:"customerCode,omitempty"`
	PONumber        string `json:"PONumber,omitempty"`
	SalesID         string `json:"salesID,omitempty"`
	TermsCode       string `json:"termsCode,omitempty"`
	Territory       string `json:"territory,omitempty"`
	CurrencyCode    string `json:"currencyCode,omitempty"`
	DateEntered     string `json:"dateEntered,omitempty"`
	MainDueDate     string `json:"mainDueDate,omitempty"`
	NotesToCustomer string `json:"notesToCustomer,omitempty"`
}

// OrderUpdate is a partial update payload for /api/v1/orders/{orderNumber}.
type OrderUpdate struct {
	PONumber        *string `json:"PONumber,omitempty"`
	SalesID         *string `json:"salesID,omitempty"`
	TermsCode       *string `json:"termsCode,omitempty"`
	Territory       *string `json:"territory,omitempty"`
	CurrencyCode    *string `json:"currencyCode,omitempty"`
	MainDueDate     *string `json:"mainDueDate,omitempty"`
	NotesToCustomer *string `json:"notesToCustomer,omitempty"`
}

// GetOrders fetches orders from /api/v1/orders.
func (c *Client) GetOrders(ctx context.Context, query url.Values) (ListResponse[Order], error) {
	return GetJSON[ListResponse[Order]](ctx, c, "/api/v1/orders", query)
}

// GetAllOrderLineItems fetches order line items from /api/v1/order-line-items.
func (c *Client) GetAllOrderLineItems(ctx context.Context, query url.Values) (ListResponse[OrderLineItem], error) {
	return GetJSON[ListResponse[OrderLineItem]](ctx, c, "/api/v1/order-line-items", query)
}

// GetAllOrderRoutings fetches order routings from /api/v1/order-routings.
func (c *Client) GetAllOrderRoutings(ctx context.Context, query url.Values) (ListResponse[OrderRouting], error) {
	return GetJSON[ListResponse[OrderRouting]](ctx, c, "/api/v1/order-routings", query)
}

// GetOrderByNumber fetches a single order by order number.
func (c *Client) GetOrderByNumber(ctx context.Context, orderNumber string, query url.Values) (SingleResponse[Order], error) {
	orderNumber = strings.TrimSpace(orderNumber)
	if orderNumber == "" {
		return SingleResponse[Order]{}, fmt.Errorf("order number is required")
	}
	endpoint := "/api/v1/orders/" + url.PathEscape(orderNumber)
	return GetJSON[SingleResponse[Order]](ctx, c, endpoint, query)
}

// CreateOrder creates an order using /api/v1/orders.
func (c *Client) CreateOrder(ctx context.Context, payload OrderCreate) (SingleResponse[Order], error) {
	return PostJSON[OrderCreate, SingleResponse[Order]](ctx, c, "/api/v1/orders", payload)
}

// UpdateOrder updates an order using PATCH /api/v1/orders/{orderNumber}.
func (c *Client) UpdateOrder(ctx context.Context, orderNumber string, payload OrderUpdate) (UpdateResult, error) {
	orderNumber = strings.TrimSpace(orderNumber)
	if orderNumber == "" {
		return UpdateResult{}, fmt.Errorf("order number is required")
	}
	endpoint := "/api/v1/orders/" + url.PathEscape(orderNumber)
	return PatchNoContentJSON[OrderUpdate](ctx, c, endpoint, payload)
}

// OrderLineItem models an Order Line Item payload.
type OrderLineItem struct {
	UniqueID                  int     `json:"uniqueID,omitempty"`
	OrderNumber               string  `json:"orderNumber,omitempty"`
	ItemNumber                int     `json:"itemNumber,omitempty"`
	PartNumber                string  `json:"partNumber,omitempty"`
	PartDescription           string  `json:"partDescription,omitempty"`
	ProductCode               string  `json:"productCode,omitempty"`
	WorkCode                  string  `json:"workCode,omitempty"`
	Status                    string  `json:"status,omitempty"`
	QuantityOrdered           float64 `json:"quantityOrdered,omitempty"`
	QuantityToMake            float64 `json:"quantityToMake,omitempty"`
	QuantityShippedToCustomer float64 `json:"quantityShippedToCustomer,omitempty"`
	DueDate                   string  `json:"dueDate,omitempty"`
	Revision                  string  `json:"revision,omitempty"`
	UnitPrice                 float64 `json:"unitPrice,omitempty"`
	UnitPriceForeign          float64 `json:"unitPriceForeign,omitempty"`
}

// OrderLineItemCreate is a minimal payload for POST /api/v1/orders/{orderNumber}/order-line-items.
type OrderLineItemCreate struct {
	PartNumber      string  `json:"partNumber,omitempty"`
	PartDescription string  `json:"partDescription,omitempty"`
	ProductCode     string  `json:"productCode,omitempty"`
	WorkCode        string  `json:"workCode,omitempty"`
	QuantityOrdered float64 `json:"quantityOrdered,omitempty"`
	DueDate         string  `json:"dueDate,omitempty"`
	UnitPrice       float64 `json:"unitPrice,omitempty"`
	Revision        string  `json:"revision,omitempty"`
	QuoteNumber     string  `json:"quoteNumber,omitempty"`
	QuoteItemNumber int     `json:"quoteItemNumber,omitempty"`
}

// OrderRouting models an Order Routing payload.
type OrderRouting struct {
	UniqueID         int     `json:"uniqueID,omitempty"`
	OrderNumber      string  `json:"orderNumber,omitempty"`
	ItemNumber       int     `json:"itemNumber,omitempty"`
	StepNumber       int     `json:"stepNumber,omitempty"`
	OperationCode    string  `json:"operationCode,omitempty"`
	Description      string  `json:"description,omitempty"`
	DepartmentNumber string  `json:"departmentNumber,omitempty"`
	WorkCenter       string  `json:"workCenter,omitempty"`
	VendorCode       string  `json:"vendorCode,omitempty"`
	Status           string  `json:"status,omitempty"`
	Total            float64 `json:"total,omitempty"`
	TimeUnit         string  `json:"timeUnit,omitempty"`
}

// OrderRoutingCreate is a minimal payload for POST order routings.
type OrderRoutingCreate struct {
	StepNumber       int     `json:"stepNumber,omitempty"`
	OperationCode    string  `json:"operationCode,omitempty"`
	Description      string  `json:"description,omitempty"`
	DepartmentNumber string  `json:"departmentNumber,omitempty"`
	WorkCenter       string  `json:"workCenter,omitempty"`
	VendorCode       string  `json:"vendorCode,omitempty"`
	TimeUnit         string  `json:"timeUnit,omitempty"`
	Total            float64 `json:"total,omitempty"`
}

// Release models an Order Release payload.
type Release struct {
	UniqueID               int     `json:"uniqueID,omitempty"`
	OrderNumber            string  `json:"orderNumber,omitempty"`
	ItemNumber             int     `json:"itemNumber,omitempty"`
	PartNumber             string  `json:"partNumber,omitempty"`
	PartDescription        string  `json:"partDescription,omitempty"`
	Quantity               float64 `json:"quantity,omitempty"`
	DueDate                string  `json:"dueDate,omitempty"`
	DateComplete           string  `json:"dateComplete,omitempty"`
	DeliveryTicketNumber   string  `json:"deliveryTicketNumber,omitempty"`
	DestinationJobNumber   string  `json:"destinationJobNumber,omitempty"`
	ManufacturingJobNumber string  `json:"manufacturingJobNumber,omitempty"`
	BinLocation            string  `json:"binLocation,omitempty"`
	Comments               string  `json:"comments,omitempty"`
	LotNumber              string  `json:"lotNumber,omitempty"`
}

// ReleaseCreate is a minimal payload for POST releases.
type ReleaseCreate struct {
	Quantity     float64 `json:"quantity,omitempty"`
	DueDate      string  `json:"dueDate,omitempty"`
	BinLocation  string  `json:"binLocation,omitempty"`
	Comments     string  `json:"comments,omitempty"`
	LotNumber    string  `json:"lotNumber,omitempty"`
	DeliveryType int     `json:"deliveryType,omitempty"`
}

// CreateOrderLineItem creates a line item for an order.
func (c *Client) CreateOrderLineItem(ctx context.Context, orderNumber string, payload OrderLineItemCreate) (SingleResponse[OrderLineItem], error) {
	orderNumber = strings.TrimSpace(orderNumber)
	if orderNumber == "" {
		return SingleResponse[OrderLineItem]{}, fmt.Errorf("order number is required")
	}
	endpoint := "/api/v1/orders/" + url.PathEscape(orderNumber) + "/order-line-items"
	return PostJSON[OrderLineItemCreate, SingleResponse[OrderLineItem]](ctx, c, endpoint, payload)
}

// GetOrderLineItem fetches a line item by order number and item number.
func (c *Client) GetOrderLineItem(ctx context.Context, orderNumber, itemNumber string, query url.Values) (SingleResponse[OrderLineItem], error) {
	orderNumber = strings.TrimSpace(orderNumber)
	if orderNumber == "" {
		return SingleResponse[OrderLineItem]{}, fmt.Errorf("order number is required")
	}
	itemNumber = strings.TrimSpace(itemNumber)
	if itemNumber == "" {
		return SingleResponse[OrderLineItem]{}, fmt.Errorf("item number is required")
	}
	endpoint := "/api/v1/orders/" + url.PathEscape(orderNumber) + "/order-line-items/" + url.PathEscape(itemNumber)
	return GetJSON[SingleResponse[OrderLineItem]](ctx, c, endpoint, query)
}

// CreateOrderRouting creates a routing for an order line item.
func (c *Client) CreateOrderRouting(ctx context.Context, orderNumber, itemNumber string, payload OrderRoutingCreate) (SingleResponse[OrderRouting], error) {
	orderNumber = strings.TrimSpace(orderNumber)
	if orderNumber == "" {
		return SingleResponse[OrderRouting]{}, fmt.Errorf("order number is required")
	}
	itemNumber = strings.TrimSpace(itemNumber)
	if itemNumber == "" {
		return SingleResponse[OrderRouting]{}, fmt.Errorf("item number is required")
	}
	endpoint := "/api/v1/orders/" + url.PathEscape(orderNumber) + "/order-line-items/" + url.PathEscape(itemNumber) + "/order-routings"
	return PostJSON[OrderRoutingCreate, SingleResponse[OrderRouting]](ctx, c, endpoint, payload)
}

// GetOrderRouting fetches a routing by order number, item number, and step number.
func (c *Client) GetOrderRouting(ctx context.Context, orderNumber, itemNumber, stepNumber string, query url.Values) (SingleResponse[OrderRouting], error) {
	orderNumber = strings.TrimSpace(orderNumber)
	if orderNumber == "" {
		return SingleResponse[OrderRouting]{}, fmt.Errorf("order number is required")
	}
	itemNumber = strings.TrimSpace(itemNumber)
	if itemNumber == "" {
		return SingleResponse[OrderRouting]{}, fmt.Errorf("item number is required")
	}
	stepNumber = strings.TrimSpace(stepNumber)
	if stepNumber == "" {
		return SingleResponse[OrderRouting]{}, fmt.Errorf("step number is required")
	}
	endpoint := "/api/v1/orders/" + url.PathEscape(orderNumber) + "/order-line-items/" + url.PathEscape(itemNumber) + "/order-routings/" + url.PathEscape(stepNumber)
	return GetJSON[SingleResponse[OrderRouting]](ctx, c, endpoint, query)
}

// CreateOrderRelease creates a release for an order line item.
func (c *Client) CreateOrderRelease(ctx context.Context, orderNumber, itemNumber string, payload ReleaseCreate) (SingleResponse[Release], error) {
	orderNumber = strings.TrimSpace(orderNumber)
	if orderNumber == "" {
		return SingleResponse[Release]{}, fmt.Errorf("order number is required")
	}
	itemNumber = strings.TrimSpace(itemNumber)
	if itemNumber == "" {
		return SingleResponse[Release]{}, fmt.Errorf("item number is required")
	}
	endpoint := "/api/v1/orders/" + url.PathEscape(orderNumber) + "/order-line-items/" + url.PathEscape(itemNumber) + "/releases"
	return PostJSON[ReleaseCreate, SingleResponse[Release]](ctx, c, endpoint, payload)
}

// GetOrderRelease fetches a release by order number, item number, and unique ID.
func (c *Client) GetOrderRelease(ctx context.Context, orderNumber, itemNumber, uniqueID string, query url.Values) (SingleResponse[Release], error) {
	orderNumber = strings.TrimSpace(orderNumber)
	if orderNumber == "" {
		return SingleResponse[Release]{}, fmt.Errorf("order number is required")
	}
	itemNumber = strings.TrimSpace(itemNumber)
	if itemNumber == "" {
		return SingleResponse[Release]{}, fmt.Errorf("item number is required")
	}
	uniqueID = strings.TrimSpace(uniqueID)
	if uniqueID == "" {
		return SingleResponse[Release]{}, fmt.Errorf("unique id is required")
	}
	endpoint := "/api/v1/orders/" + url.PathEscape(orderNumber) + "/order-line-items/" + url.PathEscape(itemNumber) + "/releases/" + url.PathEscape(uniqueID)
	return GetJSON[SingleResponse[Release]](ctx, c, endpoint, query)
}
