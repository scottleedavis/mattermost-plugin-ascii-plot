package main

import (
	"strings"
	"testing"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMessageWillBePosted(t *testing.T) {

	t.Run("message without asciiplot are unmodified", func(t *testing.T) {
		setupAPI := func() *plugintest.API {
			api := &plugintest.API{}
			api.On("LogInfo", mock.Anything).Maybe()
			return api
		}

		api := setupAPI()
		defer api.AssertExpectations(t)
		p := &Plugin{}
		p.API = api

		post := &model.Post{
			Id:      model.NewId(),
			Message: "This is a test",
		}
		// (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string)
		outPost, output := p.MessageWillBePosted(nil, post)
		assert.Equal(t, output, "")
		assert.Nil(t, outPost)
	})

	t.Run("message with asciiplot have a correct plot", func(t *testing.T) {
		setupAPI := func() *plugintest.API {
			api := &plugintest.API{}
			api.On("LogInfo", mock.Anything).Maybe()
			return api
		}

		api := setupAPI()
		defer api.AssertExpectations(t)
		p := &Plugin{}
		p.API = api
		p.configuration = &configuration{
			Width:  "50",
			Height: "15",
		}

		post := &model.Post{
			Id:      model.NewId(),
			Message: "asciiplot 1,2",
		}
		// (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string)
		outPost, output := p.MessageWillBePosted(nil, post)
		assert.Equal(t, output, "")
		graph := "```" + `
 2.00 ┤                                               ╭─ 
 1.93 ┤                                            ╭──╯  
 1.87 ┤                                        ╭───╯     
 1.80 ┤                                     ╭──╯         
 1.73 ┤                                  ╭──╯            
 1.67 ┤                               ╭──╯               
 1.60 ┤                           ╭───╯                  
 1.53 ┤                        ╭──╯                      
 1.47 ┤                     ╭──╯                         
 1.40 ┤                 ╭───╯                            
 1.33 ┤              ╭──╯                                
 1.27 ┤           ╭──╯                                   
 1.20 ┤        ╭──╯                                      
 1.13 ┤    ╭───╯                                         
 1.07 ┤ ╭──╯                                             
 1.00 ┼─╯                                                
` + "```"
		//fmt.Printf(outPost.Message)
		//fmt.Println()
		//fmt.Printf(graph)
		assert.True(t, strings.Contains(outPost.Message, graph))
		assert.Equal(t, outPost.Message, graph)
	})
}
