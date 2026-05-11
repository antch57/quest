package jambase

type Event struct {
	Name     string
	DoorTime string
	Venue    string
	Address  string
}

type apiResponse struct {
	Success    bool          `json:"success"`
	Pagination apiPagination `json:"pagination"`
	Events     []apiEvent    `json:"events"`
}

type apiPagination struct {
	Page         int    `json:"page"`
	PerPage      int    `json:"perPage"`
	TotalItems   int    `json:"totalItems"`
	TotalPages   int    `json:"totalPages"`
	NextPage     string `json:"nextPage"`
	PreviousPage string `json:"previousPage"`
}

type apiEvent struct {
	Type        string         `json:"@type"`
	Name        string         `json:"name"`
	Identifier  string         `json:"identifier"`
	URL         string         `json:"url"`
	DoorTime    string         `json:"doorTime"`
	StartDate   string         `json:"startDate"`
	EndDate     string         `json:"endDate"`
	EventStatus string         `json:"eventStatus"`
	Location    apiVenue       `json:"location"`
	Performers  []apiPerformer `json:"performer"`
}

type apiVenue struct {
	Name    string     `json:"name"`
	Address apiAddress `json:"address"`
}

type apiPerformer struct {
	Type              string   `json:"@type"`
	Name              string   `json:"name"`
	Identifier        string   `json:"identifier"`
	URL               string   `json:"url"`
	Image             string   `json:"image"`
	DatePublished     string   `json:"datePublished"`
	DateModified      string   `json:"dateModified"`
	BandOrMusician    string   `json:"x-bandOrMusician"`
	NumUpcomingEvents int      `json:"x-numUpcomingEvents"`
	Genre             []string `json:"genre"`
	PerformanceDate   string   `json:"x-performanceDate"`
	PerformanceRank   int      `json:"x-performanceRank"`
	IsHeadliner       bool     `json:"x-isHeadliner"`
	DateIsConfirmed   bool     `json:"x-dateIsConfirmed"`
}

type apiAddress struct {
	StreetAddress   string           `json:"streetAddress"`
	AddressLocality string           `json:"addressLocality"`
	PostalCode      string           `json:"postalCode"`
	AddressRegion   apiNamedLocation `json:"addressRegion"`
	AddressCountry  apiNamedLocation `json:"addressCountry"`
}

type apiNamedLocation struct {
	Type          string `json:"@type"`
	Identifier    string `json:"identifier"`
	Name          string `json:"name"`
	AlternateName string `json:"alternateName"`
}

type apiCity struct {
	Success bool          `json:"success"`
	Cities  []apiCityInfo `json:"cities"`
}

type apiCityInfo struct {
	Type       string              `json:"@type"`
	Name       string              `json:"name"`
	Identifier string              `json:"identifier"`
	Metro      apiContainedInPlace `json:"containedInPlace"`
}

type apiContainedInPlace struct {
	Type       string `json:"@type"`
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
}
