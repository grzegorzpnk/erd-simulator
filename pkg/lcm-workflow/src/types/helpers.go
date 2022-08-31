package types

type NotifyReason string
type CellId string
type ApiHandler string
type ClientListenerUri string

// There should be a fixed list of Notification Reasons
const (
	CELLCHANGED NotifyReason = "CELL_ID_CHANGED"
)

// CellChangedInfo is sent to the subscriber when current cell ID changed in the LOCATION_REPORT notification
type CellChangedInfo struct {
	Reason NotifyReason `json:"notify-reason"`
	Cell   CellId       `json:"new-cell-id"`
}

// UnsubscribeBody TODO it can be used to unsubscribe based on Endpoint and EventTypes
type UnsubscribeBody struct {
	NotifyEndpoint ClientListenerUri `json:"event-notify-uri"`
	EventTypes     []AmfEventType    `json:"amf-events"`
}
