package mockserver

type PetCategory struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Pet struct {
	Id        int         `json:"id"`
	Name      string      `json:"name"`
	Category  PetCategory `json:"category"`
	PhotoUrls []string    `json:"photoUrls"`
	Status    string      `json:"status"`
}
