package ffxiv

import (
	"bytes"
	"encoding/binary"
	"errors"
	"unsafe"
)

// These are constants for known message sizes
const (
	MessageHeaderSize                            = int(unsafe.Sizeof(MessageHeader{}))
	MessageMarketTaxRatesSize                    = int(unsafe.Sizeof(MessageMarketTaxRates{}))
	MessageMarketBoardItemListingCountSize       = int(unsafe.Sizeof(MessageMarketBoardItemListingCount{}))
	MessageMarketBoardItemListingSize            = int(unsafe.Sizeof(MessageMarketBoardItemListing{}))
	MessageMarketBoardItemListingHistorySize     = int(unsafe.Sizeof(MessageMarketBoardItemListingHistory{}))
	MessageMarketBoardSearchResultSize           = int(unsafe.Sizeof(MessageMarketBoardSearchResult{}))
	MessageMarketBoardRequestItemListingInfoSize = int(unsafe.Sizeof(MessageMarketBoardRequestItemListingInfo{}))
)

// These are known constants from the game
const (
	townCount = 6
)

// Message is composed of a message header and type-specific data
type Message interface {
	IsMessage()
}

// MessageHeader is the header of an FFXIV message
type MessageHeader struct {
	Length        uint32  `json:"message_length"`
	SourceActorID uint32  `json:"source_actor_id"`
	TargetActorID uint32  `json:"target_actor_id"`
	ControlType   uint16  `json:"control_type"`
	Unknown1      [4]byte `json:"-"`
	Type          uint16  `json:"message_type"`
	Unknown2      [2]byte `json:"-"`
	ServerID      uint16  `json:"server_id"`
	Timestamp     uint32  `json:"timestamp"`
	Unknown3      uint32  `json:"-"`
}

// IsMessage exists for anything embedding a message header
func (MessageHeader) IsMessage() {}

// MessageBytes is a generic message containing a header and raw
// bytes
type MessageBytes struct {
	MessageHeader
	Bytes []byte `json:"bytes"`
}

// MessageMarketTaxRates is a message
type MessageMarketTaxRates struct {
	MessageHeader
	Unknown1 [8]byte           `json:"-"`
	TaxRates [townCount]uint32 `json:"tax_rates"`
	Unknown2 [8]byte           `json:"-"`
}

// MessageMarketBoardItemListingCount is a message
type MessageMarketBoardItemListingCount struct {
	MessageHeader
	ItemCatalogID uint32  `json:"item_catalog_id"`
	Unknown1      [4]byte `json:"-"`
	RequestID     uint16  `json:"request_id"`
	Quantity      uint16  `json:"quantity"`
	Unknown2      [4]byte `json:"-"`
}

// MarketItemListing is part of a message
type MarketItemListing struct {
	ListingID       uint64    `json:"listing_id"`
	RetainerID      uint64    `json:"retainer_id"`
	RetainerOwnerID uint64    `json:"retainer_owner_id"`
	ArtisanID       uint64    `json:"artisan_id"`
	PricePerUnit    uint32    `json:"price_per_unit"`
	TotalTax        uint32    `json:"total_tax"`
	ItemQuantity    uint32    `json:"item_quantity"`
	ItemID          uint32    `json:"item_id"`
	LastReviewTime  uint16    `json:"last_review_time"`
	ContainerID     uint16    `json:"container_id"`
	SlotID          uint32    `json:"slot_id"`
	Durability      uint16    `json:"durability"`
	SpiritBond      uint16    `json:"spirit_bond"`
	MateriaValue    [5]uint16 `json:"materia_value"`
	Unknown1        [6]byte   `json:"-"`
	RetainerName    [32]byte  `json:"retainer_name"`
	PlayerName      [32]byte  `json:"player_name"`
	HQ              bool      `json:"hq"`
	MateriaCount    uint8     `json:"materia_count"`
	OnMannequin     uint8     `json:"on_mannequin"`
	MarketCity      uint8     `json:"market_city"`
	DyeID           uint16    `json:"dye_id"`
	Unknown2        [6]byte   `json:"-"`
}

// MessageMarketBoardItemListing is a message
type MessageMarketBoardItemListing struct {
	MessageHeader
	Listing           [10]MarketItemListing `json:"listing"`
	ListingIndexEnd   uint8                 `json:"listing_index_end"`
	ListingIndexStart uint8                 `json:"listing_index_start"`
	RequestID         uint16                `json:"request_id"`
	Unknown           [4]byte               `json:"-"`
}

// MarketHistoryListing is part of a message
type MarketHistoryListing struct {
	SalePrice     uint32   `json:"sale_price"`
	PurchaseTime  uint32   `json:"purchase_time"`
	Quantity      uint32   `json:"quantity"`
	HQ            uint8    `json:"hq"`
	Unknown       byte     `json:"-"`
	OnMannequin   uint8    `json:"on_mannequin"`
	BuyerName     [33]byte `json:"buyer_name"`
	ItemCatalogID uint32   `json:"item_catalog_id"`
}

// MessageMarketBoardItemListingHistory is a message
type MessageMarketBoardItemListingHistory struct {
	MessageHeader
	ItemCatalogID [2]uint32                `json:"item_catalog_id"`
	Listing       [20]MarketHistoryListing `json:"listing"`
}

// MarketBoardItem is part of a message
type MarketBoardItem struct {
	ItemCatalogID uint32 `json:"item_catalog_id"`
	Quantity      uint16 `json:"quantity"`
	Demand        uint16 `json:"demand"`
}

// MessageMarketBoardSearchResult is a message
type MessageMarketBoardSearchResult struct {
	MessageHeader
	Items          [20]MarketBoardItem `json:"items"`
	ItemIndexEnd   uint32              `json:"item_index_end"`
	Unknown        [4]byte             `json:"-"`
	ItemIndexStart uint32              `json:"item_index_start"`
	RequestID      uint32              `json:"request_id"`
}

// MessageMarketBoardRequestItemListingInfo is a message
type MessageMarketBoardRequestItemListingInfo struct {
	MessageHeader
	CatalogID uint32 `json:"catalog_id"`
	RequestID uint32 `json:"request_id"`
}

// NewMessageHeader returns a new message header from a given byte array
func NewMessageHeader(b []byte) (*MessageHeader, error) {
	if len(b) != MessageHeaderSize {
		return nil, errors.New("invalid message header size")
	}
	m := &MessageHeader{}
	r := bytes.NewReader(b)
	if err := binary.Read(r, binary.LittleEndian, m); err != nil {
		return nil, err
	}
	return m, nil
}

// NewMessage returns a new message from a given byte array based on its type
func NewMessage(b []byte) (Message, error) {
	if len(b) < MessageHeaderSize {
		return nil, errors.New("message length too small")
	}
	mh, err := NewMessageHeader(b[:MessageHeaderSize])
	if err != nil {
		return nil, err
	}
	if uint32(len(b)) != mh.Length {
		return nil, errors.New("message length invalid")
	}
	if mh.ControlType != 3 {
		return mh, nil
	}
	mp, ok := messageParsers[mh.Type]
	if !ok {
		return nil, errors.New("unrecognized message type")
	}
	return mp(b)
}

// NewMessageBytes returns a new message from a given byte array containing a message
// header and the raw bytes for the remaining message.
func NewMessageBytes(b []byte) (Message, error) {
	mh, err := NewMessageHeader(b[:MessageHeaderSize])
	if err != nil {
		return nil, err
	}
	return &MessageBytes{
		MessageHeader: *mh,
		Bytes:         b[MessageHeaderSize:],
	}, nil
}
