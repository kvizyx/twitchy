package helix

type serviceBase struct {
	client *Client
}

type AdsService struct{ serviceBase }
type AnalyticsService struct{ serviceBase }
type BitsService struct{ serviceBase }
type ChannelsService struct{ serviceBase }
type ChannelPointsService struct{ serviceBase }
type CharityService struct{ serviceBase }
type ChatService struct{ serviceBase }
type ClipsService struct{ serviceBase }
type ConduitsService struct{ serviceBase }
type CCLsService struct{ serviceBase }
type EntitlementsService struct{ serviceBase }
type ExtensionsService struct{ serviceBase }
type EventSubService struct{ serviceBase }
type GamesService struct{ serviceBase }
type GoalsService struct{ serviceBase }
type HypeTrainService struct{ serviceBase }
type ModerationService struct{ serviceBase }
type PollsService struct{ serviceBase }
type PredictionsService struct{ serviceBase }
type RaidsService struct{ serviceBase }
type ScheduleService struct{ serviceBase }
type SearchService struct{ serviceBase }
type StreamsService struct{ serviceBase }
type SubscriptionsService struct{ serviceBase }
type TagsService struct{ serviceBase }
type TeamsService struct{ serviceBase }
type UsersService struct{ serviceBase }
type VideosService struct{ serviceBase }
type WhispersService struct{ serviceBase }

func initializeServices(client *Client) {
	base := serviceBase{client: client}
	client.Ads = &AdsService{serviceBase: base}
	client.Analytics = &AnalyticsService{serviceBase: base}
	client.Bits = &BitsService{serviceBase: base}
	client.Channels = &ChannelsService{serviceBase: base}
	client.ChannelPoints = &ChannelPointsService{serviceBase: base}
	client.Charity = &CharityService{serviceBase: base}
	client.Chat = &ChatService{serviceBase: base}
	client.Clips = &ClipsService{serviceBase: base}
	client.Conduits = &ConduitsService{serviceBase: base}
	client.CCLs = &CCLsService{serviceBase: base}
	client.Entitlements = &EntitlementsService{serviceBase: base}
	client.Extensions = &ExtensionsService{serviceBase: base}
	client.EventSub = &EventSubService{serviceBase: base}
	client.Games = &GamesService{serviceBase: base}
	client.Goals = &GoalsService{serviceBase: base}
	client.HypeTrain = &HypeTrainService{serviceBase: base}
	client.Moderation = &ModerationService{serviceBase: base}
	client.Polls = &PollsService{serviceBase: base}
	client.Predictions = &PredictionsService{serviceBase: base}
	client.Raids = &RaidsService{serviceBase: base}
	client.Schedule = &ScheduleService{serviceBase: base}
	client.Search = &SearchService{serviceBase: base}
	client.Streams = &StreamsService{serviceBase: base}
	client.Subscriptions = &SubscriptionsService{serviceBase: base}
	client.Tags = &TagsService{serviceBase: base}
	client.Teams = &TeamsService{serviceBase: base}
	client.Users = &UsersService{serviceBase: base}
	client.Videos = &VideosService{serviceBase: base}
	client.Whispers = &WhispersService{serviceBase: base}
	initializeExperimentalServices(client)
}
