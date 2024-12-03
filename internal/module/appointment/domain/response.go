package domain

// AppointmentResponse represents the response for appointment operations
type AppointmentResponse struct {
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination represents pagination metadata
type Pagination struct {
	CurrentPage  int   `json:"current_page"`
	TotalItems   int64 `json:"total_items"`
	TotalPages   int   `json:"total_pages"`
	ItemsPerPage int   `json:"items_per_page"`
}

// AvailabilityResponse represents the response for availability check
type AvailabilityResponse struct {
	Available bool   `json:"available"`
	Message   string `json:"message,omitempty"`
}
