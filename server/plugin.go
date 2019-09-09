package main

import (
	"strconv"
	"strings"
	"sync"

	"github.com/guptarohit/asciigraph"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	configurationLock sync.RWMutex

	configuration *configuration
}

//MessageWillBePosted hook
func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {

	if strings.HasPrefix(post.Message, "asciiplot ") || strings.HasPrefix(post.Message, "asciigraph ") {
		pointsString := strings.TrimPrefix(post.Message, "asciiplot ")
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
		height, err := strconv.ParseInt(configuration.Height, 10, 8)
		if err != nil {
			height = 15
		}
		width, err := strconv.ParseInt(configuration.Width, 10, 8)
		if err != nil {
			width = 60
		}
		graph := asciigraph.Plot(numbers, asciigraph.Height(int(height)), asciigraph.Width(int(width)))
		post.Message = "```\n" + graph + "\n```"
		return post, ""

	}
	return nil, ""
}
