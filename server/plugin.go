package main

import (
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/guptarohit/asciigraph"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

//Plugin strucutre
type Plugin struct {
	plugin.MattermostPlugin

	configurationLock sync.RWMutex

	configuration *configuration
}

//
const (
	BASE    = 10
	BITSIZE = 8
	HEIGHT  = 15
	WIDTH   = 60
)

//MessageWillBePosted hook
func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {

	regx := regexp.MustCompile(`(asciiplot|asciigraph)\s(\d+)((,\s*\d+)|(\s*,\s*\d+))*`)
	matches := regx.FindStringSubmatch(post.Message)
	if len(matches) > 0 {
		pointsString := strings.TrimPrefix(matches[0], "asciiplot ")
		pointsString = strings.TrimPrefix(pointsString, "asciigraph ")
		pointsStringArray := strings.Split(pointsString, ",")

		var numbers []float64

		if len(pointsStringArray) <= 1 {
			return nil, ""
		}

		for _, arg := range pointsStringArray {
			arg = strings.Trim(arg, " ")
			if n, err := strconv.ParseFloat(arg, 64); err == nil {
				numbers = append(numbers, n)
			}
		}
		configuration := p.getConfiguration()
		height, err := strconv.ParseInt(configuration.Height, BASE, BITSIZE)
		if err != nil {
			height = HEIGHT
		}
		width, err := strconv.ParseInt(configuration.Width, BASE, BITSIZE)
		if err != nil {
			width = WIDTH
		}
		graph := asciigraph.Plot(numbers, asciigraph.Height(int(height)), asciigraph.Width(int(width)))
		post.Message = strings.Replace(post.Message, matches[0], "```\n"+graph+"\n```", 1)
		return post, ""
	}
	return nil, ""
}
