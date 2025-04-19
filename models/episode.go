package models

type EpisodeRequest struct {
	Payload []Episode `json:"payload"`
	Skip    int       `json:"skip"`
	Take    int       `json:"take"`
	Total   int       `json:"totalRecords"`
}

type Episode struct {
	Country      string       `json:"country"`
	Description  string       `json:"description"`
	DRM          bool         `json:"drm"`
	EpisodeCount int          `json:"episodeCount"`
	Genre        string       `json:"genre"`
	Image        Image        `json:"image"`
	Language     string       `json:"language"`
	NextEpisode  *NextEpisode `json:"nextEpisode"`
	PrimaryColor string       `json:"primaryColour"`
	Seasons      []Season     `json:"seasons"`
	Slug         string       `json:"slug"`
	Title        string       `json:"title"`
	TVChannel    string       `json:"tvChannel"`
}

type Image struct {
	ShowImage string `json:"showImage"`
}

type NextEpisode struct {
	Channel     string `json:"channel"`
	ChannelLogo string `json:"channelLogo"`
	Date        string `json:"date"`
	HTML        string `json:"html"`
	URL         string `json:"url"`
}

type Season struct {
	Slug string `json:"slug"`
}

type EpisodeResponse struct {
	Response []EpisodeResponseItem `json:"response"`
}

type EpisodeResponseItem struct {
	Image string `json:"image"`
	Slug  string `json:"slug"`
	Title string `json:"title"`
}
