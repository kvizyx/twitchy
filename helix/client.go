package helix

import (
	"errors"
	"net/http"
	"net/url"
)

var (
	ErrInvalidOption      = errors.New("helix: invalid option")
	ErrConflictingOptions = errors.New("helix: conflicting options")
	ErrInvalidClient      = errors.New("helix: invalid client")
)

const defaultBaseURL = "https://api.twitch.tv/helix"

type Client struct {
	Ads           *AdsService
	Analytics     *AnalyticsService
	Bits          *BitsService
	Channels      *ChannelsService
	ChannelPoints *ChannelPointsService
	Charity       *CharityService
	Chat          *ChatService
	Clips         *ClipsService
	Conduits      *ConduitsService
	CCLs          *CCLsService
	Entitlements  *EntitlementsService
	Extensions    *ExtensionsService
	EventSub      *EventSubService
	Games         *GamesService
	Goals         *GoalsService
	HypeTrain     *HypeTrainService
	Moderation    *ModerationService
	Polls         *PollsService
	Predictions   *PredictionsService
	Raids         *RaidsService
	Schedule      *ScheduleService
	Search        *SearchService
	Streams       *StreamsService
	Subscriptions *SubscriptionsService
	Tags          *TagsService
	Teams         *TeamsService
	Users         *UsersService
	Videos        *VideosService
	Whispers      *WhispersService

	Experimental Experimental

	httpClient      *http.Client
	baseURL         *url.URL
	userAgent       string
	tokenSource     TokenSource
	staticToken     *Credential
	rateLimitPolicy RateLimitPolicy
	executor        *transportExecutor
	valid           bool
}

func New(options ...Option) (*Client, error) {
	baseURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		return nil, err
	}
	client := &Client{
		httpClient: cloneHTTPClient(http.DefaultClient),
		baseURL:    baseURL,
		userAgent:  "twitchy/helix",
	}
	for _, option := range options {
		if option == nil {
			return nil, ErrInvalidOption
		}
		if err := option(client); err != nil {
			return nil, err
		}
	}
	if client.tokenSource != nil && client.staticToken != nil {
		return nil, ErrConflictingOptions
	}
	if client.httpClient == nil || client.baseURL == nil {
		return nil, ErrInvalidClient
	}
	source := client.tokenSource
	if client.staticToken != nil {
		source = NewStaticTokenSource(*client.staticToken)
	}
	client.executor = newTransportExecutor(client.httpClient, source, client.rateLimitPolicy, nil, nil)
	initializeServices(client)
	client.valid = true
	return client, nil
}

func (c *Client) validClient() error {
	if c == nil || !c.valid {
		return ErrInvalidClient
	}
	return nil
}
