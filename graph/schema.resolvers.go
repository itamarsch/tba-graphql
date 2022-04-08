package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"tba-gql/graph/generated"
	"tba-gql/graph/model"

	"golang.org/x/sync/errgroup"
)

func (r *queryResolver) TeamByNumber(ctx context.Context, number int) (*model.Team, error) {
	return FetchTeam(number)
}

func (r *queryResolver) TeamByPageNum(ctx context.Context, pageNum int) ([]*model.Team, error) {
	keys, err := Fetch[[]string](fmt.Sprintf("https://www.thebluealliance.com/api/v3/teams/%v/keys", pageNum))
	wg := sync.WaitGroup{}
	eg := errgroup.Group{}
	var teams []*model.Team = make([]*model.Team, len(*keys))
	eg.Go(func() error { return err })
	for i, v := range *keys {
		wg.Add(1)
		go func(i int, key string) {
			fmt.Printf("Starting %v\n", i)
			num, convErr := strconv.Atoi(strings.TrimPrefix(key, "frc"))
			eg.Go(func() error { return convErr })
			if convErr != nil {
				return
			}
			team, err := FetchTeam(num)
			eg.Go(func() error { return err })
			teams[i] = team
			wg.Done()
			fmt.Printf("Finished %v\n", i)
		}(i, v)
	}

	wg.Wait()
	if e := eg.Wait(); e != nil {
		return nil, e
	}

	return teams, nil
}

func (r *teamResolver) Events(ctx context.Context, obj *model.Team, where *model.EventComparisonExp) ([]*model.Event, error) {
	if where != nil {
		if where.Key != nil {
			return Filter(obj.Events, func(t *model.Event) bool { return *t.Key == *where.Key.Eq }), nil
		}
	}
	return obj.Events, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Team returns generated.TeamResolver implementation.
func (r *Resolver) Team() generated.TeamResolver { return &teamResolver{r} }

type queryResolver struct{ *Resolver }
type teamResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func Filter[T any](arr []T, cond func(T) bool) []T {
	result := []T{}
	for i := range arr {
		if cond(arr[i]) {
			result = append(result, arr[i])
		}
	}
	return result
}
