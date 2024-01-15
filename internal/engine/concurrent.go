package engine

type Concurrent struct {
	Scheduler        Scheduler
	WorkerCount      int
	ItemChan         chan Item
	RequestProcessor Processor
}
type Processor func(Request) (ParseResult, error)

type Scheduler interface {
	Submit(request Request)
	ReadyNotify
	WorkerChan() chan Request
	Run()
}
type ReadyNotify interface {
	WorkerReady(chan Request)
}

func (e *Concurrent) Run(seeds ...Request) {
	out := make(chan ParseResult)
	e.Scheduler.Run()

	for i := 0; i < e.WorkerCount; i++ {
		e.createWorker(e.Scheduler.WorkerChan(), out, e.Scheduler)

	}

	for _, r := range seeds {
		if isDuplicate(r.Url) {
			continue
		}
		e.Scheduler.Submit(r)
	}

	for {
		res := <-out
		for _, item := range res.Items {
			//fmt.Printf("%v -> item %v\n", i, item)
			go func(i Item) {
				e.ItemChan <- i
			}(item)
		}

		for _, r := range res.Requests {
			if isDuplicate(r.Url) {
				continue
			}

			e.Scheduler.Submit(r)
		}
	}

}

func (e *Concurrent) createWorker(in chan Request, out chan ParseResult, ready ReadyNotify) {
	go func() {
		for {
			ready.WorkerReady(in)
			r := <-in
			res, err := Worker(r)
			if err != nil {
				continue
			}
			out <- res
		}
	}()
}

var visitedUrls = make(map[string]bool)

func isDuplicate(url string) bool {
	if visitedUrls[url] {
		return true
	}
	visitedUrls[url] = true
	return false
}
