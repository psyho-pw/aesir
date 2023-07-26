package slackbot

import (
	"aesir/src/channels"
	"aesir/src/common"
	"github.com/gofiber/fiber/v2"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func TestNewRouter(t *testing.T) {
	type args struct {
		router  fiber.Router
		db      *gorm.DB
		handler Handler
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewRouter(tt.args.router, tt.args.db, tt.args.handler)
		})
	}
}

func TestNewSlackHandler(t *testing.T) {
	type args struct {
		service Service
	}
	tests := []struct {
		name string
		args args
		want Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSlackHandler(tt.args.service); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSlackHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSlackService(t *testing.T) {
	type args struct {
		config         *common.Config
		channelService channels.Service
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSlackService(tt.args.config, tt.args.channelService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSlackService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_slackHandler_EventMux(t *testing.T) {
	type fields struct {
		service Service
	}
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := slackHandler{
				service: tt.fields.service,
			}
			if err := handler.EventMux(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("EventMux() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_slackHandler_FindChannelById(t *testing.T) {
	type fields struct {
		service Service
	}
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := slackHandler{
				service: tt.fields.service,
			}
			if err := handler.FindChannelById(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("FindChannelById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_slackHandler_FindChannels(t *testing.T) {
	type fields struct {
		service Service
	}
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := slackHandler{
				service: tt.fields.service,
			}
			if err := handler.FindChannels(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("FindChannels() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_slackHandler_FindLatestChannelMessage(t *testing.T) {
	type fields struct {
		service Service
	}
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := slackHandler{
				service: tt.fields.service,
			}
			if err := handler.FindLatestChannelMessage(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("FindLatestChannelMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_slackHandler_FindTeam(t *testing.T) {
	type fields struct {
		service Service
	}
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := slackHandler{
				service: tt.fields.service,
			}
			if err := handler.FindTeam(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("FindTeam() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_slackHandler_FindTeamUsers(t *testing.T) {
	type fields struct {
		service Service
	}
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := slackHandler{
				service: tt.fields.service,
			}
			if err := handler.FindTeamUsers(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("FindTeamUsers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_slackHandler_WhoAmI(t *testing.T) {
	type fields struct {
		service Service
	}
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := slackHandler{
				service: tt.fields.service,
			}
			if err := handler.WhoAmI(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("WhoAmI() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_slackService_EventMux(t *testing.T) {
	type fields struct {
		api            *slack.Client
		channelService channels.Service
	}
	type args struct {
		innerEvent slackevents.EventsAPIInnerEvent
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &slackService{
				api:            tt.fields.api,
				channelService: tt.fields.channelService,
			}
			if err := service.EventMux(tt.args.innerEvent); (err != nil) != tt.wantErr {
				t.Errorf("EventMux() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_slackService_FindChannel(t *testing.T) {
	type fields struct {
		api            *slack.Client
		channelService channels.Service
	}
	type args struct {
		channelId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *slack.Channel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &slackService{
				api:            tt.fields.api,
				channelService: tt.fields.channelService,
			}
			got, err := service.FindChannel(tt.args.channelId)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindChannel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindChannel() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_slackService_FindChannels(t *testing.T) {
	type fields struct {
		api            *slack.Client
		channelService channels.Service
	}
	tests := []struct {
		name    string
		fields  fields
		want    []slack.Channel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &slackService{
				api:            tt.fields.api,
				channelService: tt.fields.channelService,
			}
			got, err := service.FindChannels()
			if (err != nil) != tt.wantErr {
				t.Errorf("FindChannels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindChannels() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_slackService_FindJoinedChannels(t *testing.T) {
	type fields struct {
		api            *slack.Client
		channelService channels.Service
	}
	tests := []struct {
		name    string
		fields  fields
		want    []slack.Channel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &slackService{
				api:            tt.fields.api,
				channelService: tt.fields.channelService,
			}
			got, err := service.FindJoinedChannels()
			if (err != nil) != tt.wantErr {
				t.Errorf("FindJoinedChannels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindJoinedChannels() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_slackService_FindLatestChannelMessage(t *testing.T) {
	type fields struct {
		api            *slack.Client
		channelService channels.Service
	}
	type args struct {
		channelId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *slack.Message
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &slackService{
				api:            tt.fields.api,
				channelService: tt.fields.channelService,
			}
			got, err := service.FindLatestChannelMessage(tt.args.channelId)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindLatestChannelMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindLatestChannelMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_slackService_FindTeam(t *testing.T) {
	type fields struct {
		api            *slack.Client
		channelService channels.Service
	}
	tests := []struct {
		name    string
		fields  fields
		want    *slack.TeamInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &slackService{
				api:            tt.fields.api,
				channelService: tt.fields.channelService,
			}
			got, err := service.FindTeam()
			if (err != nil) != tt.wantErr {
				t.Errorf("FindTeam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindTeam() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_slackService_FindTeamUsers(t *testing.T) {
	type fields struct {
		api            *slack.Client
		channelService channels.Service
	}
	type args struct {
		teamId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []slack.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &slackService{
				api:            tt.fields.api,
				channelService: tt.fields.channelService,
			}
			got, err := service.FindTeamUsers(tt.args.teamId)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindTeamUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindTeamUsers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_slackService_WhoAmI(t *testing.T) {
	type fields struct {
		api            *slack.Client
		channelService channels.Service
	}
	tests := []struct {
		name    string
		fields  fields
		want    *WhoAmI
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &slackService{
				api:            tt.fields.api,
				channelService: tt.fields.channelService,
			}
			got, err := service.WhoAmI()
			if (err != nil) != tt.wantErr {
				t.Errorf("WhoAmI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WhoAmI() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_slackService_WithTx(t *testing.T) {
	type fields struct {
		api            *slack.Client
		channelService channels.Service
	}
	type args struct {
		tx *gorm.DB
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Service
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &slackService{
				api:            tt.fields.api,
				channelService: tt.fields.channelService,
			}
			if got := service.WithTx(tt.args.tx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithTx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_slackService_handleMemberJoinEvent(t *testing.T) {
	type fields struct {
		api            *slack.Client
		channelService channels.Service
	}
	type args struct {
		event *slackevents.MemberJoinedChannelEvent
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &slackService{
				api:            tt.fields.api,
				channelService: tt.fields.channelService,
			}
			if err := service.handleMemberJoinEvent(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("handleMemberJoinEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_slackService_handleMessageEvent(t *testing.T) {
	type fields struct {
		api            *slack.Client
		channelService channels.Service
	}
	type args struct {
		event *slackevents.MessageEvent
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &slackService{
				api:            tt.fields.api,
				channelService: tt.fields.channelService,
			}
			if err := service.handleMessageEvent(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("handleMessageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
