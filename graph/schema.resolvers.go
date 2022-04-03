package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"sync"
	"tba-gql/graph/generated"
	"tba-gql/graph/model"

	"golang.org/x/sync/errgroup"
)

func FetchAsync[T any](wg *sync.WaitGroup, t **T, url string, eg *errgroup.Group) {
	wg.Add(1)
	go func(t **T) {
		var err error
		*t, err = Fetch[T](url)
		eg.Go(func() error {
			return err
		})
		wg.Done()
	}(t)
}

func (r *queryResolver) TeamByKey(ctx context.Context, key string) (*model.Team, error) {
	wg := sync.WaitGroup{}
	eg := errgroup.Group{}
	var team *model.Team
	FetchAsync(&wg, &team, fmt.Sprintf("https://www.thebluealliance.com/api/v3/team/%v", key), &eg)

	var yearsPaticipated *[]int
	FetchAsync(&wg, &yearsPaticipated, fmt.Sprintf("https://www.thebluealliance.com/api/v3/team/%v/years_participated", key), &eg)

	var robots *[]*model.Robot
	FetchAsync(&wg, &robots, fmt.Sprintf("https://www.thebluealliance.com/api/v3/team/%v/robots", key), &eg)

	var districts *[]*model.District
	FetchAsync(&wg, &districts, fmt.Sprintf("https://www.thebluealliance.com/api/v3/team/%v/districts", key), &eg)
	wg.Wait()
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	team.YearsParticipated = *yearsPaticipated
	team.Robots = *robots
	team.Districts = *districts
	return team, nil
}

func (r *queryResolver) TeamByPageNum(ctx context.Context, pageNum int) ([]*model.Team, error) {
	keys, err := Fetch[[]string](fmt.Sprintf("https://www.thebluealliance.com/api/v3/teams/%v/keys", pageNum))
	wg := sync.WaitGroup{}

	var teams []*model.Team = make([]*model.Team, len(*keys))

	for i, v := range *keys {
		wg.Add(1)
		go func(i int, key string) {
			team, err := r.TeamByKey(ctx, key)
			if err != nil {
				panic(err)
			}
			teams[i] = team
			wg.Done()
		}(i, v)
	}

	wg.Wait()
	return teams, err
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
