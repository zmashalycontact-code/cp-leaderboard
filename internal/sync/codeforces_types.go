package sync

// cfResponse is the top-level envelope returned by every Codeforces API endpoint.
//
//	{ "status": "OK", "result": [...] }
//	{ "status": "FAILED", "comment": "handles: User with handle X not found" }
type cfResponse struct {
	Status  string   `json:"status"`
	Comment string   `json:"comment"` // populated only on FAILED
	Result  []cfUser `json:"result"`
}

// cfUser mirrors the Codeforces `User` object returned by /api/user.info.
// Only the fields we actually persist are decoded; the rest are silently
// ignored by encoding/json.
//
// Reference: https://codeforces.com/apiHelp/objects#User
type cfUser struct {
	Handle    string `json:"handle"`
	Rating    int    `json:"rating"`    // current rating; 0 if unrated
	MaxRating int    `json:"maxRating"` // all-time peak; 0 if unrated
}
