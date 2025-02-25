package models

// CreateDeal represents the request body for creating a deal.
// Fields and examples are based on the Pipedrive API documentation.
type CreateDeal struct {
	Title             string   `json:"title" example:"Test Deal"`
	Value             string   `json:"value" example:"1000"`
	Label             []int    `json:"label" example:"[1,2,3]"`
	Currency          string   `json:"currency" example:"USD"`
	UserID            int      `json:"user_id" example:"123"`
	PersonID          int      `json:"person_id" example:"456"`
	OrgID             int      `json:"org_id" example:"789"`
	PipelineID        int      `json:"pipeline_id" example:"10"`
	StageID           int      `json:"stage_id" example:"20"`
	Status            string   `json:"status" example:"open" enums:"open,won,lost,deleted"`
	OriginID          string   `json:"origin_id" example:"integration_xyz"`
	Channel           int      `json:"channel" example:"1"`
	ChannelID         string   `json:"channel_id" example:"ch_123"`
	AddTime           string   `json:"add_time" example:"2023-08-21 12:34:56"`
	WonTime           string   `json:"won_time" example:"2023-08-22 13:00:00"`
	LostTime          string   `json:"lost_time" example:"2023-08-22 14:00:00"`
	CloseTime         string   `json:"close_time" example:"2023-08-23 15:00:00"`
	ExpectedCloseDate string   `json:"expected_close_date" example:"2023-08-30"`
	Probability       float64  `json:"probability" example:"75.5"`
	LostReason        string   `json:"lost_reason" example:"Price too high"`
	VisibleTo         string   `json:"visible_to" example:"1" enums:"1,3,5,7"`
}

// UpdateDeal represents the request body for updating a deal.
// All fields are optional.
type UpdateDeal struct {
	Title             *string  `json:"title,omitempty" example:"Updated Deal Title"`
	Value             *string  `json:"value,omitempty" example:"1500"`
	Label             []int    `json:"label,omitempty" example:"[1,2]"`
	Currency          *string  `json:"currency,omitempty" example:"USD"`
	UserID            *int     `json:"user_id,omitempty" example:"123"`
	PersonID          *int     `json:"person_id,omitempty" example:"456"`
	OrgID             *int     `json:"org_id,omitempty" example:"789"`
	PipelineID        *int     `json:"pipeline_id,omitempty" example:"10"`
	StageID           *int     `json:"stage_id,omitempty" example:"20"`
	Status            *string  `json:"status,omitempty" example:"open" enums:"open,won,lost,deleted"`
	Channel           *int     `json:"channel,omitempty" example:"1"`
	ChannelID         *string  `json:"channel_id,omitempty" example:"ch_123"`
	WonTime           *string  `json:"won_time,omitempty" example:"2023-08-22 13:00:00"`
	LostTime          *string  `json:"lost_time,omitempty" example:"2023-08-22 14:00:00"`
	CloseTime         *string  `json:"close_time,omitempty" example:"2023-08-23 15:00:00"`
	ExpectedCloseDate *string  `json:"expected_close_date,omitempty" example:"2023-08-30"`
	Probability       *float64 `json:"probability,omitempty" example:"75.5"`
	LostReason        *string  `json:"lost_reason,omitempty" example:"Price too high"`
	VisibleTo         *string  `json:"visible_to,omitempty" example:"1" enums:"1,3,5,7"`
}
