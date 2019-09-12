package main

import (
	"io/ioutil"
	"os"
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
		outPost, output := p.MessageWillBePosted(nil, post)
		assert.Equal(t, output, "")
		graph := "\n```" + `
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
` + "```\n"
		assert.True(t, strings.Contains(outPost.Message, graph))
		assert.Equal(t, outPost.Message, graph)

		post = &model.Post{
			Id:      model.NewId(),
			Message: "asciigraph 1,2",
		}
		outPost, output = p.MessageWillBePosted(nil, post)
		assert.Equal(t, output, "")
	})

	t.Run("message with asciiplot located in not at beginning have a correct plot", func(t *testing.T) {
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
			Message: "Let's look at the figures.\nasciiplot 1,2",
		}
		outPost, output := p.MessageWillBePosted(nil, post)
		assert.Equal(t, output, "")
		graph := "Let's look at the figures.\n\n```" + `
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
` + "```\n"
		assert.True(t, strings.Contains(outPost.Message, graph))
		assert.Equal(t, outPost.Message, graph)
	})

	t.Run("message with multiple asciiplots", func(t *testing.T) {
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
			Message: "Let's look at the figures.\nasciiplot 1,2\nasciiplot 1,2",
		}
		outPost, output := p.MessageWillBePosted(nil, post)
		assert.Equal(t, output, "")
		file, _ := os.Open("../assets/multiple_test.txt")
		graph, _ := ioutil.ReadAll(file)

		assert.True(t, strings.Contains(outPost.Message, string(graph)))
		assert.Equal(t, outPost.Message, string(graph))
	})
}
