package response

type ReadQueue struct {
	Request struct {
		Mbean string `json:"mbean"`
		Type  string `json:"type"`
	} `json:"request"`
	Value struct {
		ProducerFlowControl  bool    `json:"ProducerFlowControl"`
		Options              string  `json:"Options"`
		AlwaysRetroactive    bool    `json:"AlwaysRetroactive"`
		MemoryUsageByteCount int     `json:"MemoryUsageByteCount"`
		AverageBlockedTime   float64 `json:"AverageBlockedTime"`
		MemoryPercentUsage   int     `json:"MemoryPercentUsage"`
		CursorMemoryUsage    int     `json:"CursorMemoryUsage"`
		InFlightCount        int     `json:"InFlightCount"`
		Subscriptions        []struct {
			ObjectName string `json:"objectName"`
		} `json:"Subscriptions"`
		CacheEnabled        bool    `json:"CacheEnabled"`
		ForwardCount        int     `json:"ForwardCount"`
		DLQ                 bool    `json:"DLQ"`
		StoreMessageSize    int     `json:"StoreMessageSize"`
		AverageEnqueueTime  float64 `json:"AverageEnqueueTime"`
		Name                string  `json:"Name"`
		MaxAuditDepth       int     `json:"MaxAuditDepth"`
		BlockedSends        int     `json:"BlockedSends"`
		TotalBlockedTime    int     `json:"TotalBlockedTime"`
		QueueSize           int     `json:"QueueSize"`
		MaxPageSize         int     `json:"MaxPageSize"`
		PrioritizedMessages bool    `json:"PrioritizedMessages"`
		MemoryUsagePortion  float64 `json:"MemoryUsagePortion"`
		Paused              bool    `json:"Paused"`
		EnqueueCount        int     `json:"EnqueueCount"`
		MessageGroups       struct {
		} `json:"MessageGroups"`
		ConsumerCount                  int         `json:"ConsumerCount"`
		AverageMessageSize             int         `json:"AverageMessageSize"`
		CursorFull                     bool        `json:"CursorFull"`
		MaxProducersToAudit            int         `json:"MaxProducersToAudit"`
		ExpiredCount                   int         `json:"ExpiredCount"`
		CursorPercentUsage             int         `json:"CursorPercentUsage"`
		MinEnqueueTime                 int         `json:"MinEnqueueTime"`
		MinMessageSize                 int         `json:"MinMessageSize"`
		MemoryLimit                    int         `json:"MemoryLimit"`
		DispatchCount                  int         `json:"DispatchCount"`
		MaxEnqueueTime                 int         `json:"MaxEnqueueTime"`
		DequeueCount                   int         `json:"DequeueCount"`
		BlockedProducerWarningInterval int         `json:"BlockedProducerWarningInterval"`
		ProducerCount                  int         `json:"ProducerCount"`
		MessageGroupType               string      `json:"MessageGroupType"`
		MaxMessageSize                 int         `json:"MaxMessageSize"`
		UseCache                       bool        `json:"UseCache"`
		SlowConsumerStrategy           interface{} `json:"SlowConsumerStrategy"`
	} `json:"value"`
	Timestamp int `json:"timestamp"`
	ErrorType string `json:"error_type"`
	Error string `json:"error"`
	Status    int `json:"status"`
}
