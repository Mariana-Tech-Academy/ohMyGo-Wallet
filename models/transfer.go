package models



type TransferRequest struct {
	RecipientAcc string  `json:"recipient_account" binding:"required,len=10"`  //we bind to the recipient's account number 
	Amount       float64 `json:"amount" binding:"required,gt=0"`
	Description  string  `json:"description" binding:"max=255"`
}


//repo query can handle the strings for account numbers for now changed from uint to string


//the user ID is associated with the transfer request and authenticated route
// figma image requests for recipient not recipient's acc num
// for every transfer, 2 transactions occur.
// 1. debit for the user
// 2. credit for the receiver --- //this will be dealt with seperately in transactions.go 
