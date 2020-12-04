package ffxiv

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var messageParsers = map[uint16]MessageParser{
	0x025F: parseMessageMarketBoardItemListing,
	0x038F: parseMessageMarketBoardItemListingCount,
	0x0186: parseMessageMarketBoardItemListingHistory,
	0x0102: parseMessageMarketBoardRequestItemListingInfo,
	0x032C: parseMessageMarketBoardSearchResult,
	0x01F8: parseMessageMarketTaxRates,
}

// MessageParser is a function to parse a message out of a byte slice
type MessageParser func([]byte) (Message, error)

func parseMessageMarketTaxRates(b []byte) (Message, error) {
	if len(b) != MessageMarketTaxRatesSize {
		return nil, errors.New("invalid message size")
	}
	msg := &MessageMarketTaxRates{}
	r := bytes.NewReader(b)
	if err := binary.Read(r, binary.LittleEndian, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func parseMessageMarketBoardItemListingCount(b []byte) (Message, error) {
	if len(b) != MessageMarketBoardItemListingCountSize {
		return nil, errors.New("invalid message size")
	}
	msg := &MessageMarketBoardItemListingCount{}
	r := bytes.NewReader(b)
	if err := binary.Read(r, binary.LittleEndian, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func parseMessageMarketBoardItemListing(b []byte) (Message, error) {
	if len(b) != MessageMarketBoardItemListingSize {
		return nil, errors.New("invalid message size")
	}
	msg := &MessageMarketBoardItemListing{}
	r := bytes.NewReader(b)
	if err := binary.Read(r, binary.LittleEndian, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func parseMessageMarketBoardItemListingHistory(b []byte) (Message, error) {
	if len(b) != MessageMarketBoardItemListingHistorySize {
		return nil, errors.New("invalid message size")
	}
	msg := &MessageMarketBoardItemListingHistory{}
	r := bytes.NewReader(b)
	if err := binary.Read(r, binary.LittleEndian, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func parseMessageMarketBoardSearchResult(b []byte) (Message, error) {
	if len(b) != MessageMarketBoardSearchResultSize {
		return nil, errors.New("invalid message size")
	}
	msg := &MessageMarketBoardSearchResult{}
	r := bytes.NewReader(b)
	if err := binary.Read(r, binary.LittleEndian, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func parseMessageMarketBoardRequestItemListingInfo(b []byte) (Message, error) {
	if len(b) != MessageMarketBoardRequestItemListingInfoSize {
		return nil, errors.New("invalid message size")
	}
	msg := &MessageMarketBoardRequestItemListingInfo{}
	r := bytes.NewReader(b)
	if err := binary.Read(r, binary.LittleEndian, msg); err != nil {
		return nil, err
	}
	return msg, nil
}
