package jambase

// Event is the normalized show record returned by SearchShows.
type Event struct {
	Name     string
	Date     string
	DoorTime string
	Venue    string
	Address  string
	Timezone string
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
	EventStatus string         `json:"eventStatus"`
	StartDate   string         `json:"startDate"`
	EndDate     string         `json:"endDate"`
	DoorTime    string         `json:"doorTime"`
	Location    apiVenue       `json:"location"`
	Performers  []apiPerformer `json:"performer"`
}

type apiVenue struct {
	Name       string     `json:"name"`
	Identifier string     `json:"identifier"`
	Address    apiAddress `json:"address"`
}

type apiPerformer struct {
	Type              string   `json:"@type"`
	Name              string   `json:"name"`
	Identifier        string   `json:"identifier"`
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
	Timezone        string           `json:"x-timezone"`
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
