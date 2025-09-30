package handlers

import (
	"banana-auction/rest/handlers/auction"
	"banana-auction/rest/handlers/bid"
	"banana-auction/rest/handlers/lot"
	"banana-auction/rest/handlers/user"
)

var (
	CreateLotHandler     = lot.Create
	UpdateLotHandler     = lot.Update
	DeleteLotHandler     = lot.DeleteLot
	ListLotHandler       = lot.List
	CreateAuctionHandler = auction.Create
	ListBidsHandler      = auction.List
	CreateBidHandler     = bid.Create
	SignupHandler        = user.Signup
	LoginHandler         = user.Login
)
