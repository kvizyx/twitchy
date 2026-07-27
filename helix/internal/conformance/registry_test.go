package conformance

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

var serviceTypes = map[string]reflect.Type{
	"AdsService":                    reflect.TypeOf((*helix.AdsService)(nil)),
	"AnalyticsService":              reflect.TypeOf((*helix.AnalyticsService)(nil)),
	"BitsService":                   reflect.TypeOf((*helix.BitsService)(nil)),
	"ChannelsService":               reflect.TypeOf((*helix.ChannelsService)(nil)),
	"ChannelPointsService":          reflect.TypeOf((*helix.ChannelPointsService)(nil)),
	"CharityService":                reflect.TypeOf((*helix.CharityService)(nil)),
	"ChatService":                   reflect.TypeOf((*helix.ChatService)(nil)),
	"ClipsService":                  reflect.TypeOf((*helix.ClipsService)(nil)),
	"ConduitsService":               reflect.TypeOf((*helix.ConduitsService)(nil)),
	"CCLsService":                   reflect.TypeOf((*helix.CCLsService)(nil)),
	"EntitlementsService":           reflect.TypeOf((*helix.EntitlementsService)(nil)),
	"ExtensionsService":             reflect.TypeOf((*helix.ExtensionsService)(nil)),
	"EventSubService":               reflect.TypeOf((*helix.EventSubService)(nil)),
	"GamesService":                  reflect.TypeOf((*helix.GamesService)(nil)),
	"GoalsService":                  reflect.TypeOf((*helix.GoalsService)(nil)),
	"HypeTrainService":              reflect.TypeOf((*helix.HypeTrainService)(nil)),
	"ModerationService":             reflect.TypeOf((*helix.ModerationService)(nil)),
	"PollsService":                  reflect.TypeOf((*helix.PollsService)(nil)),
	"PredictionsService":            reflect.TypeOf((*helix.PredictionsService)(nil)),
	"RaidsService":                  reflect.TypeOf((*helix.RaidsService)(nil)),
	"ScheduleService":               reflect.TypeOf((*helix.ScheduleService)(nil)),
	"SearchService":                 reflect.TypeOf((*helix.SearchService)(nil)),
	"StreamsService":                reflect.TypeOf((*helix.StreamsService)(nil)),
	"SubscriptionsService":          reflect.TypeOf((*helix.SubscriptionsService)(nil)),
	"TagsService":                   reflect.TypeOf((*helix.TagsService)(nil)),
	"TeamsService":                  reflect.TypeOf((*helix.TeamsService)(nil)),
	"UsersService":                  reflect.TypeOf((*helix.UsersService)(nil)),
	"VideosService":                 reflect.TypeOf((*helix.VideosService)(nil)),
	"WhispersService":               reflect.TypeOf((*helix.WhispersService)(nil)),
	"ExperimentalBitsService":       reflect.TypeOf((*helix.ExperimentalBitsService)(nil)),
	"ExperimentalChatService":       reflect.TypeOf((*helix.ExperimentalChatService)(nil)),
	"ExperimentalClipsService":      reflect.TypeOf((*helix.ExperimentalClipsService)(nil)),
	"ExperimentalGuestStarService":  reflect.TypeOf((*helix.ExperimentalGuestStarService)(nil)),
	"ExperimentalModerationService": reflect.TypeOf((*helix.ExperimentalModerationService)(nil)),
	"ExperimentalUsersService":      reflect.TypeOf((*helix.ExperimentalUsersService)(nil)),
}

func resolveService(client *helix.Client, selector string) (reflect.Value, error) {
	parts := strings.Split(selector, ".")
	if len(parts) < 2 || parts[0] != "Client" {
		return reflect.Value{}, fmt.Errorf("invalid selector %q", selector)
	}
	value := reflect.ValueOf(client).Elem()
	for _, part := range parts[1:] {
		value = value.FieldByName(part)
		if !value.IsValid() {
			return reflect.Value{}, fmt.Errorf("selector %q is not exported on client", selector)
		}
		if value.Kind() == reflect.Pointer && value.IsNil() {
			return reflect.Value{}, fmt.Errorf("selector %q resolves to nil", selector)
		}
	}
	return value, nil
}

func resolveMethod(operation manifest.Operation, client *helix.Client) (reflect.Value, error) {
	service, err := resolveService(client, operation.Implementation.Selector)
	if err != nil {
		return reflect.Value{}, err
	}
	method := service.MethodByName(operation.Implementation.Method)
	if !method.IsValid() {
		return reflect.Value{}, fmt.Errorf("%s: method %s is missing", operation.Anchor, operation.Implementation.Method)
	}
	if service.Type() != serviceTypes[operation.Implementation.ServiceType] {
		return reflect.Value{}, fmt.Errorf("%s: selector has type %s, want %s", operation.Anchor, service.Type(), operation.Implementation.ServiceType)
	}
	if err := verifyMethodType(operation, method.Type()); err != nil {
		return reflect.Value{}, err
	}
	return method, nil
}

func verifyMethodType(operation manifest.Operation, got reflect.Type) error {
	contextType := reflect.TypeOf((*context.Context)(nil)).Elem()
	if got.NumIn() != 2 || got.In(0) != contextType || got.In(1).Name() != operation.Implementation.RequestType || got.In(1).PkgPath() != "github.com/kvizyx/twitchy/helix" {
		return fmt.Errorf("%s: method signature has wrong inputs: %s", operation.Anchor, got)
	}
	if got.NumOut() != 2 || got.Out(0).Kind() != reflect.Pointer || !strings.Contains(got.Out(0).String(), "Response[") || got.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		return fmt.Errorf("%s: method signature has wrong outputs: %s", operation.Anchor, got)
	}
	return nil
}

func verifyServiceSurface(loaded manifest.Manifest) error {
	expected := make(map[string]struct{}, len(loaded.Operations))
	for _, operation := range loaded.Operations {
		expected[operation.Implementation.ServiceType+"."+operation.Implementation.Method] = struct{}{}
	}
	actual := make(map[string]struct{}, len(expected))
	for serviceName, serviceType := range serviceTypes {
		for index := 0; index < serviceType.NumMethod(); index++ {
			method := serviceType.Method(index)
			if strings.HasSuffix(method.Name, "Pager") {
				continue
			}
			actual[serviceName+"."+method.Name] = struct{}{}
		}
	}
	for mapping := range expected {
		if _, ok := actual[mapping]; !ok {
			return fmt.Errorf("missing compiled method %s", mapping)
		}
	}
	for mapping := range actual {
		if _, ok := expected[mapping]; !ok {
			return fmt.Errorf("extra compiled method %s", mapping)
		}
	}
	if len(actual) != 149 {
		return fmt.Errorf("compiled method count: got %d, want 149", len(actual))
	}
	return nil
}
