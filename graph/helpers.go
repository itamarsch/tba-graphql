package graph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"tba-gql/graph/model"

	"golang.org/x/sync/errgroup"
)

func Fetch[T any](url string) (*T, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-TBA-Auth-Key", "Q7k2X4Erid8KqVXRt82GICzJAC9ZkIpBogQHJ4sIwJACU9u29bHNE6AElJWRKMk4")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	jsonResponse, _ := ioutil.ReadAll(res.Body)

	a := new(T)
	json.Unmarshal(jsonResponse, a)
	return a, nil
}

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

func FetchTeam(number int) (*model.Team, error) {
	key := fmt.Sprintf("frc%v", number)
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

	var events *[]*model.Event
	FetchAsync(&wg, &events, fmt.Sprintf("https://www.thebluealliance.com/api/v3/team/%v/events", key), &eg)
	wg.Wait()
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	team.YearsParticipated = *yearsPaticipated
	team.Robots = *robots
	team.Districts = *districts
	team.Events = *events
	return team, nil
}
