package exactonline

import "github.com/mcnijman/go-exactonline/types"

// SubscriptionLine stores SubscriptionLine from exactonline
//
type SubscriptionLine struct {
	Item     types.GUID `json:"Item"`
	FromDate types.Date `json:"FromDate"`
}
