package constants

const (
	ConstHeaderIssuerId       = "issuer_id"
	ConstTenantTypeAccessRead = "READ"
	ConstTenantHeaderContext  = "TENANT_PER_REQUEST"
	ConstInfoColor            = "\033[1;34m%s\033[0m"
	ConstNoticeColor          = "\033[1;36m%s\033[0m"
	ConstWarningColor         = "\033[1;33m%s\033[0m"
	ConstErrorColor           = "\033[1;31m%s\033[0m"
	ConstDebugColor           = "\033[0;36m%s\033[0m"
)

const (
	EventTypePerson            = "person"
	EventTypeAddress           = "address"
	EventTypePhone             = "phone"
	EventTypeAnalysis          = "creditanalysis"
	EventTypeAccount           = "account"
	EventTypeCard              = "card"
	EventTypeLimit             = "limit"
	EventTypePurchase          = "purchase"
	EventTypePurchaseApproved  = "purchase_approved"
	EventTypePurchaseCanceled  = "purchase_canceled"
	EventTypePurchaseCleared   = "purchase_cleared"
	EventTypePurchaseProcessed = "purchase_processed"
	EventTypeDeclined          = "declined"
	EventTypeStandIn           = "standin"
	EventTypeTransaction       = "transaction"
	EventTypePayment           = "payment"
	EventTypeAdjustment        = "adjustment"
	EventTypeStatement         = "statement"
	EventTypeBankSlip          = "bankslip"
	EventTypeToken             = "token"
)

func EventTypes() []string {
	return []string{EventTypePerson, EventTypeAddress, EventTypePhone, EventTypeAnalysis, EventTypeAccount, EventTypeCard, EventTypeLimit, EventTypePurchase, EventTypeDeclined, EventTypeStandIn, EventTypeTransaction, EventTypePayment, EventTypeAdjustment, EventTypeStatement, EventTypeBankSlip, EventTypeToken, EventTypePurchaseApproved, EventTypePurchaseCanceled, EventTypePurchaseCleared, EventTypePurchaseProcessed}
}
