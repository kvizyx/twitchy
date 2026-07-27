package helix

type Experimental struct {
	Bits       *ExperimentalBitsService
	Chat       *ExperimentalChatService
	Clips      *ExperimentalClipsService
	GuestStar  *ExperimentalGuestStarService
	Moderation *ExperimentalModerationService
	Users      *ExperimentalUsersService
}

type ExperimentalBitsService struct{ serviceBase }
type ExperimentalChatService struct{ serviceBase }
type ExperimentalClipsService struct{ serviceBase }
type ExperimentalGuestStarService struct{ serviceBase }
type ExperimentalModerationService struct{ serviceBase }
type ExperimentalUsersService struct{ serviceBase }

func initializeExperimentalServices(client *Client) {
	base := serviceBase{client: client}
	client.Experimental = Experimental{
		Bits:       &ExperimentalBitsService{serviceBase: base},
		Chat:       &ExperimentalChatService{serviceBase: base},
		Clips:      &ExperimentalClipsService{serviceBase: base},
		GuestStar:  &ExperimentalGuestStarService{serviceBase: base},
		Moderation: &ExperimentalModerationService{serviceBase: base},
		Users:      &ExperimentalUsersService{serviceBase: base},
	}
}
