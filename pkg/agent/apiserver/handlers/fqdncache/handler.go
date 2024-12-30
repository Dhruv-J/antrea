package fqdncache

import (
	"encoding/json"
	"net/http"

	"k8s.io/klog/v2"

	agentquerier "antrea.io/antrea/pkg/agent/querier"
)

func HandleFunc(aq agentquerier.AgentQuerier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dnsEntryCache := aq.GetFqdnCache()
		if dnsEntryCache == nil {
			return
		}
		if err := json.NewEncoder(w).Encode(dnsEntryCache); err != nil {
			http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
			klog.ErrorS(err, "Failed to encode response")
		}
	}
}
