package bugzilla

import (
	"fmt"
	"k8s.io/test-infra/prow/bugzilla"
	"math"
	"strconv"
	"sync"
)

type TotBugzillaClient interface {
	bugzilla.Client
	QuickSearch(filters map[string]string, maxLimit uint) ([]*bugzilla.Bug, error)
}

func WrapBugzillaClient(bzClient bugzilla.Client) TotBugzillaClient {
	return &totBugzillaClient{bzClient}
}

type totBugzillaClient struct {
	bugzilla.Client
}

func (t *totBugzillaClient) QuickSearch(filters map[string]string, maxLimit uint) ([]*bugzilla.Bug, error) {
	limit := math.Min(float64(maxLimit), 50)
	maxOffset := int(math.Max(1.0, float64(maxLimit)/limit))
	maxWorkers := int(math.Min(25, float64(maxOffset)))

	doneChannel := make(chan bool)
	offsetChannel := make(chan uint, maxWorkers)
	errorChannel := make(chan error)
	bugChannel := make(chan *bugzilla.Bug)

	offsetProducer := func(doneChannel chan bool, offsetChannel chan<- uint) {
		for offset := 0; offset < maxOffset; offset++ {
			fmt.Printf("Producer: offset %d\n", offset)
			select {
			case <-doneChannel:
				break
			default:
				offsetChannel <- uint(offset)
			}
		}
		close(offsetChannel)
	}

	offsetConsumer := func(offsetChannel <-chan uint, doneChannel chan<- bool, bugChannel chan<- *bugzilla.Bug, errorChannel chan<- error) {
		for {
			fmt.Println("Consumer: waiting for offset")
			offset, ok := <-offsetChannel
			if !ok {
				fmt.Println("Consumer: offset channel closed...bailing")
				return
			}

			fmt.Printf("Consumer: requesting %d/%d\n", offset, limit)

			request := copyMap(filters)
			request["limit"] = strconv.Itoa(int(limit))
			request["offset"] = strconv.Itoa(int(offset))
			bugs, err := t.SearchBugs(request)
			if err != nil {
				fmt.Printf("Consumer: Got error %s\n", err)
				errorChannel <- err
			}

			fmt.Println("Got response, sending bugs")
			if len(bugs) == 0 {
				doneChannel <- true
			}

			for _, bug := range bugs {
				bugChannel <- bug
			}
		}
	}

	bugs := make([]*bugzilla.Bug, 0)

	// start offset producer
	go offsetProducer(doneChannel, offsetChannel)

	// start offset consumers
	var waitGroup sync.WaitGroup
	go func() {
		for workerId := 0; workerId < maxWorkers; workerId++ {
			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()
				offsetConsumer(offsetChannel, doneChannel, bugChannel, errorChannel)
			}()
		}
		waitGroup.Wait()
		close(bugChannel)
	}()

	// collect bugs
	fmt.Println("Collector: starting")
	for {
		select {
		case bug, ok := <-bugChannel:
			if ok {
				bugs = append(bugs, bug)
			} else {
				fmt.Println("Collector: bug channel closed")
				return bugs, nil
			}
		case err := <-errorChannel:
			fmt.Printf("Collector: got error %s\n", err)
			fmt.Printf("Closing offset producer")
			doneChannel <- true
			return nil, err
		}
	}
}

func copyMap(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}
