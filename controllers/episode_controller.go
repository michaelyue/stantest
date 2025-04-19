package controllers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"stan.com/stantest/models"
)

// deal with the episode data and returns filtered results
func DealwithEpisodes(c echo.Context) error {
	c.Logger().Info("received episode processing request")

	// log the raw request body
	if c.Request() != nil && c.Request().Body != nil {
		rawBody, err := io.ReadAll(c.Request().Body)
		if err != nil {
			c.Logger().Errorf("failed to read request body: %s", err.Error())
		} else {
			c.Logger().Infof("raw request body: %s", string(rawBody))
			// restore the body so it can be read again
			c.Request().Body = io.NopCloser(bytes.NewBuffer(rawBody))
		}
	}

	var request models.EpisodeRequest
	if err := c.Bind(&request); err != nil {
		c.Logger().Errorf("failed to bind request: %s", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not decode request: JSON parsing failed"})
	}

	// validate request data
	if err := validateRequest(request); err != nil {
		c.Logger().Errorf("request validation failed: %s", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not decode request: " + err.Error()})
	}

	episodeCount := len(request.Payload)
	c.Logger().Infof("processing %d episodes", episodeCount)

	// filter episodes based on our criteria
	// DRM enabled (drm: true) and at least one episode (episodeCount > 0).
	var response models.EpisodeResponse
	matchedCount := 0

	for _, episode := range request.Payload {
		c.Logger().Debugf("processing episode: %s", episode.Title)

		if episode.DRM && episode.EpisodeCount > 0 {
			// Validate episode data
			if err := validateEpisode(episode); err != nil {
				c.Logger().Warnf("skipping invalid episode %s: %s", episode.Title, err.Error())
				continue
			}

			response.Response = append(response.Response, models.EpisodeResponseItem{
				Image: episode.Image.ShowImage,
				Slug:  episode.Slug,
				Title: episode.Title,
			})
			matchedCount++
		}
	}

	c.Logger().Infof("processed %d episodes, %d matched criteria", episodeCount, matchedCount)

	// in case no any episodes matched the criteria
	// just return empty error
	if len(response.Response) == 0 {
		c.Logger().Info("no episodes matched the criteria")
		return c.JSON(http.StatusOK, models.EpisodeResponse{Response: []models.EpisodeResponseItem{}})
	}

	return c.JSON(http.StatusOK, response)
}

// payload must be not empty
func validateRequest(request models.EpisodeRequest) error {
	if request.Payload == nil {
		return fmt.Errorf("payload is required")
	}
	return nil
}

// validate an URL
func isValidURL(checkUrl string) bool {
	_, err := url.ParseRequestURI(checkUrl)
	if err != nil {
		return false
	}

	u, err := url.Parse(checkUrl)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	if u.Scheme != "https" && u.Scheme != "http" {
		return false
	}

	return true
}

// validateEpisode checks if an episode has all required fields and valid values
func validateEpisode(episode models.Episode) error {
	if episode.Title == "" {
		return echo.NewHTTPError(400, "title is required")
	}
	if episode.Slug == "" {
		return echo.NewHTTPError(400, "slug is required")
	}
	if episode.Image.ShowImage == "" {
		return echo.NewHTTPError(400, "image.showImage is required")
	}
	if !isValidURL(episode.Image.ShowImage) {
		return echo.NewHTTPError(400, "image.showImage must be a valid URL")
	}
	return nil
}
