// analysis
package invitro_server

type Analysis struct {
	Id          int    `json:"id"`
	Kind        string `json:"kind"`
	Subtype     string `json:"subtype"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
