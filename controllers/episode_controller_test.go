package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"stan.com/stantest/models"
)

func TestDealwithEpisodes(t *testing.T) {
	// create an echo instance for testing
	e := echo.New()

	// make up all testing cases
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedCount  int
		expectedError  string
	}{
		{
			name: "Valid request with matching episodes",
			requestBody: `{
				"payload": [
					{
						"country": "UK",
						"description": "What's life like when you have enough children to field your own football team?",
						"drm": true,
						"episodeCount": 3,
						"genre": "Reality",
						"image": {"showImage": "http://catchup.ninemsn.com.au/img/jump-in/shows/16KidsandCounting1280.jpg"},
						"language": "English",
						"nextEpisode": null,
						"primaryColour": "#ff7800",
						"slug": "show/16kidsandcounting",
						"title": "16 Kids and Counting",
						"tvChannel": "GEM"
					}
				]
			}`,
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name: "Valid request with multiple matching episodes",
			requestBody: `{
				"payload": [
					{
						"country": "UK",
						"description": "What's life like when you have enough children to field your own football team?",
						"drm": true,
						"episodeCount": 3,
						"genre": "Reality",
						"image": {"showImage": "http://catchup.ninemsn.com.au/img/jump-in/shows/16KidsandCounting1280.jpg"},
						"language": "English",
						"nextEpisode": null,
						"primaryColour": "#ff7800",
						"slug": "show/16kidsandcounting",
						"title": "16 Kids and Counting",
						"tvChannel": "GEM"
					},
					{
						"country": "USA",
						"description": "The Taste puts 16 culinary competitors in the kitchen, where four of the World's most notable culinary masters of the food world judges their creations based on a blind taste. Join judges Anthony Bourdain, Nigella Lawson, Ludovic Lefebvre and Brian Malarkey in this pressure-packed contest where a single spoonful can catapult a contender to the top or send them packing.",
						"drm": true,
						"episodeCount": 2,
						"genre": "Reality",
						"image": {"showImage": "http://catchup.ninemsn.com.au/img/jump-in/shows/TheTaste1280.jpg"},
						"language": "English",
						"primaryColour": "#df0000",
						"slug": "show/thetaste",
						"title": "The Taste (Le Goût)",
						"tvChannel": "GEM"
					}
				]
			}`,
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name: "Valid request with mixed matching and non-matching episodes",
			requestBody: `{
				"payload": [
					{
						"country": "UK",
						"description": "What's life like when you have enough children to field your own football team?",
						"drm": true,
						"episodeCount": 3,
						"genre": "Reality",
						"image": {"showImage": "http://catchup.ninemsn.com.au/img/jump-in/shows/16KidsandCounting1280.jpg"},
						"language": "English",
						"nextEpisode": null,
						"primaryColour": "#ff7800",
						"slug": "show/16kidsandcounting",
						"title": "16 Kids and Counting",
						"tvChannel": "GEM"
					},
					{
						"country": "USA",
						"description": "The Taste puts 16 culinary competitors in the kitchen, where four of the World's most notable culinary masters of the food world judges their creations based on a blind taste. Join judges Anthony Bourdain, Nigella Lawson, Ludovic Lefebvre and Brian Malarkey in this pressure-packed contest where a single spoonful can catapult a contender to the top or send them packing.",
						"drm": true,
						"episodeCount": 2,
						"genre": "Reality",
						"image": {"showImage": "http://catchup.ninemsn.com.au/img/jump-in/shows/TheTaste1280.jpg"},
						"language": "English",
						"primaryColour": "#df0000",
						"slug": "show/thetaste",
						"title": "The Taste (Le Goût)",
						"tvChannel": "GEM"
					},
					{
						"country": "UK",
						"description": "The series follows the adventures of International Rescue, an organisation created to help those in grave danger using technically advanced equipment and machinery. The series focuses on the head of the organisation, ex-astronaut Jeff Tracy, and his five sons who piloted the \"Thunderbird\" machines.",
						"drm": true,
						"episodeCount": 0,
						"genre": "Action",
						"image": {"showImage": "http://catchup.ninemsn.com.au/img/jump-in/shows/Thunderbirds_1280.jpg"},
						"language": "English",
						"primaryColour": "#0084da",
						"slug": "show/thunderbirds",
						"title": "Thunderbirds",
						"tvChannel": "Channel 9"
					}
				]
			}`,
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name: "Valid request with no matching episodes",
			requestBody: `{
				"payload": [
					{
						"country": "UK",
						"description": "The series follows the adventures of International Rescue, an organisation created to help those in grave danger using technically advanced equipment and machinery. The series focuses on the head of the organisation, ex-astronaut Jeff Tracy, and his five sons who piloted the \"Thunderbird\" machines.",
						"drm": true,
						"episodeCount": 0,
						"genre": "Action",
						"image": {"showImage": "http://catchup.ninemsn.com.au/img/jump-in/shows/Thunderbirds_1280.jpg"},
						"language": "English",
						"primaryColour": "#0084da",
						"slug": "show/thunderbirds",
						"title": "Thunderbirds",
						"tvChannel": "Channel 9"
					}
				]
			}`,
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{invalid json}`,
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
			expectedError:  "Could not decode request: JSON parsing failed",
		},
		{
			name: "Empty payload array",
			requestBody: `{
				"payload": []
			}`,
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "Missing required fields",
			requestBody: `{
				"payload": [
					{
						"country": "UK",
						"description": "The series follows the adventures of International Rescue, an organisation created to help those in grave danger using technically advanced equipment and machinery. The series focuses on the head of the organisation, ex-astronaut Jeff Tracy, and his five sons who piloted the \"Thunderbird\" machines.",
						"drm": true,
						"episodeCount": 0,
						"genre": "Action",
						"image": {"showImage": "http://catchup.ninemsn.com.au/img/jump-in/shows/Thunderbirds_1280.jpg"},
						"language": "English",
						"primaryColour": "#0084da",
						"slug": "",
						"title": "",
						"tvChannel": "Channel 9"
					}
				]
			}`,
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "Null payload",
			requestBody: `{
				"payload": null
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
			expectedError:  "Could not decode request: payload is required",
		},
		{
			name: "Missing payload field",
			requestBody: `{
				"skip": 0,
				"take": 10
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
			expectedError:  "Could not decode request: payload is required",
		},
		{
			name: "Episode with nextEpisode data",
			requestBody: `{
				"payload": [
					{
						"country": " USA",
						"description": "The Taste puts 16 culinary competitors in the kitchen, where four of the World's most notable culinary masters of the food world judges their creations based on a blind taste. Join judges Anthony Bourdain, Nigella Lawson, Ludovic Lefebvre and Brian Malarkey in this pressure-packed contest where a single spoonful can catapult a contender to the top or send them packing.",
						"drm": true,
						"episodeCount": 2,
						"genre": "Reality",
						"image": {
							"showImage": "http://catchup.ninemsn.com.au/img/jump-in/shows/TheTaste1280.jpg"
						},
						"language": "English",
						"nextEpisode": {
							"channel": null,
							"channelLogo": "http://catchup.ninemsn.com.au/img/player/logo_go.gif",
							"date": null,
							"html": "<br><span class=\"visit\">Visit the Official Website</span></span>",
							"url": "http://go.ninemsn.com.au/"
						},
						"primaryColour": "#df0000",
						"seasons": [
							{
								"slug": "show/thetaste/season/1"
							}
						],
						"slug": "show/thetaste",
						"title": "The Taste (Le Goût)",
						"tvChannel": "GEM"
					}
				]
			}`,
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name: "Episode with seasons data",
			requestBody: `{
				"payload": [
					{
						"country": "UK",
						"description": "The series follows the adventures of International Rescue, an organisation created to help those in grave danger using technically advanced equipment and machinery. The series focuses on the head of the organisation, ex-astronaut Jeff Tracy, and his five sons who piloted the \"Thunderbird\" machines.",
						"drm": true,
						"episodeCount": 24,
						"genre": "Action",
						"image": {
							"showImage": "http://catchup.ninemsn.com.au/img/jump-in/shows/Thunderbirds_1280.jpg"
						},
						"language": "English",
						"nextEpisode": null,
						"primaryColour": "#0084da",
						"seasons": [
							{
								"slug": "show/thunderbirds/season/1"
							},
							{
								"slug": "show/thunderbirds/season/3"
							},
							{
								"slug": "show/thunderbirds/season/4"
							},
							{
								"slug": "show/thunderbirds/season/5"
							},
							{
								"slug": "show/thunderbirds/season/6"
							},
							{
								"slug": "show/thunderbirds/season/8"
							}
						],
						"slug": "show/thunderbirds",
						"title": "Thunderbirds",
						"tvChannel": "Channel 9"
					}
				]
			}`,
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// simulate a request
			req := httptest.NewRequest(http.MethodPost, "/api/v1/episodes", bytes.NewBufferString(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := DealwithEpisodes(c)

			// check returned status code
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK {
				var response models.EpisodeResponse
				err = json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(response.Response))

				// verify response format for matching episodes
				if tt.expectedCount > 0 {
					for _, item := range response.Response {
						assert.NotEmpty(t, item.Image)
						assert.NotEmpty(t, item.Slug)
						assert.NotEmpty(t, item.Title)
					}
				}
			} else if tt.expectedError != "" {
				var errorResponse map[string]string
				err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
				assert.NoError(t, err)
				assert.Contains(t, errorResponse["error"], tt.expectedError)
			}
		})
	}
}

func TestValidateRequest(t *testing.T) {
	tests := []struct {
		name    string
		request models.EpisodeRequest
		wantErr bool
	}{
		{
			name:    "Valid request with empty payload",
			request: models.EpisodeRequest{Payload: []models.Episode{}},
			wantErr: false,
		},
		{
			name:    "Nil payload",
			request: models.EpisodeRequest{Payload: nil},
			wantErr: true,
		},
		{
			name: "Valid request with non-empty payload",
			request: models.EpisodeRequest{
				Payload: []models.Episode{
					{
						Title: "Thunderbirds",
						Slug:  "show/thunderbirds",
						Image: models.Image{ShowImage: "http://catchup.ninemsn.com.au/img/jump-in/shows/Thunderbirds_1280.jpg"},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRequest(tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateEpisode(t *testing.T) {
	tests := []struct {
		name    string
		episode models.Episode
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid episode",
			episode: models.Episode{
				Title: "Thunderbirds",
				Slug:  "show/thunderbirds",
				Image: models.Image{ShowImage: "http://catchup.ninemsn.com.au/img/jump-in/shows/Thunderbirds_1280.jpg"},
			},
			wantErr: false,
		},
		{
			name: "Missing title",
			episode: models.Episode{
				Slug:  "show/thunderbirds",
				Image: models.Image{ShowImage: "http://catchup.ninemsn.com.au/img/jump-in/shows/Thunderbirds_1280.jpg"},
			},
			wantErr: true,
			errMsg:  "title is required",
		},
		{
			name: "Missing slug",
			episode: models.Episode{
				Title: "Thunderbirds",
				Image: models.Image{ShowImage: "http://catchup.ninemsn.com.au/img/jump-in/shows/Thunderbirds_1280.jpg"},
			},
			wantErr: true,
			errMsg:  "slug is required",
		},
		{
			name: "Missing image",
			episode: models.Episode{
				Title: "Thunderbirds",
				Slug:  "show/thunderbirds",
			},
			wantErr: true,
			errMsg:  "image.showImage is required",
		},
		{
			name: "Empty title",
			episode: models.Episode{
				Title: "",
				Slug:  "show/thunderbirds",
				Image: models.Image{ShowImage: "http://catchup.ninemsn.com.au/img/jump-in/shows/Thunderbirds_1280.jpg"},
			},
			wantErr: true,
			errMsg:  "title is required",
		},
		{
			name: "Empty slug",
			episode: models.Episode{
				Title: "Thunderbirds",
				Slug:  "",
				Image: models.Image{ShowImage: "http://catchup.ninemsn.com.au/img/jump-in/shows/Thunderbirds_1280.jpg"},
			},
			wantErr: true,
			errMsg:  "slug is required",
		},
		{
			name: "Empty image URL",
			episode: models.Episode{
				Title: "Thunderbirds",
				Slug:  "show/thunderbirds",
				Image: models.Image{ShowImage: ""},
			},
			wantErr: true,
			errMsg:  "image.showImage is required",
		},
		{
			name: "Invalid image URL - missing scheme",
			episode: models.Episode{
				Title: "Thunderbirds",
				Slug:  "show/thunderbirds",
				Image: models.Image{ShowImage: "catchup.ninemsn.com.au/img/jump-in/shows/Thunderbirds_1280.jpg"},
			},
			wantErr: true,
			errMsg:  "image.showImage must be a valid URL",
		},
		{
			name: "Invalid image URL - invalid scheme",
			episode: models.Episode{
				Title: "Thunderbirds",
				Slug:  "show/thunderbirds",
				Image: models.Image{ShowImage: "ftp://catchup.ninemsn.com.au/img/jump-in/shows/Thunderbirds_1280.jpg"},
			},
			wantErr: true,
			errMsg:  "image.showImage must be a valid URL",
		},
		{
			name: "Invalid image URL - missing host",
			episode: models.Episode{
				Title: "Thunderbirds",
				Slug:  "show/thunderbirds",
				Image: models.Image{ShowImage: "http://"},
			},
			wantErr: true,
			errMsg:  "image.showImage must be a valid URL",
		},
		{
			name: "Invalid image URL - malformed URL",
			episode: models.Episode{
				Title: "Thunderbirds",
				Slug:  "show/thunderbirds",
				Image: models.Image{ShowImage: "htp://catchup.ninemsn.com.au/img/jump-in/shows/Thunderbirds_1280.jpg"},
			},
			wantErr: true,
			errMsg:  "image.showImage must be a valid URL",
		},
		{
			name: "Valid image URL with query parameters",
			episode: models.Episode{
				Title: "Thunderbirds",
				Slug:  "show/thunderbirds",
				Image: models.Image{ShowImage: "http://catchup.ninemsn.com.au/img/jump-in/shows/Thunderbirds_1280.jpg?width=400"},
			},
			wantErr: false,
		},
		{
			name: "Valid image URL with port",
			episode: models.Episode{
				Title: "Thunderbirds",
				Slug:  "show/thunderbirds",
				Image: models.Image{ShowImage: "http://catchup.ninemsn.com.au:8080/img/jump-in/shows/Thunderbirds_1280.jpg"},
			},
			wantErr: false,
		},
		{
			name: "Valid image URL with path",
			episode: models.Episode{
				Title: "Thunderbirds",
				Slug:  "show/thunderbirds",
				Image: models.Image{ShowImage: "http://catchup.ninemsn.com.au/img/jump-in/shows/Thunderbirds_1280.jpg"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEpisode(tt.episode)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
