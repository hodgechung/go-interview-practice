package main

import (
	"fmt"
	"sync"
)

type Result struct {
	Node int
	Seq  []int
}

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.
func worker(graph map[int][]int, tasks <-chan int, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range tasks {
		var result []int
		visited := map[int]bool{t: true}
		queue := []int{t}
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			result = append(result, cur)
			for _, nbr := range graph[cur] {
				if !visited[nbr] {
					queue = append(queue, nbr)
					visited[nbr] = true
				}
			}
		}

		results <- Result{
			t,
			result,
		}
	}
}

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	ans := make(map[int][]int)
	if numWorkers == 0 {
		return ans
	}
	tasks := make(chan int, numWorkers)
	results := make(chan Result, numWorkers)

	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for range numWorkers {
		go worker(graph, tasks, results, &wg)
	}

	go func() {
		for _, q := range queries {
			tasks <- q
		}
		close(tasks)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		ans[r.Node] = r.Seq
	}

	return ans
}

func main() {
	graph := map[int][]int{
		0: {1, 2},
		1: {2, 3},
		2: {3},
		3: {4},
		4: {},
	}
	queries := []int{0, 1, 2}
	numWorkers := 2

	results := ConcurrentBFSQueries(graph, queries, numWorkers)
	fmt.Printf("%+v", results)
	/*
	   Possible output:
	   results[0] = [0 1 2 3 4]
	   results[1] = [1 2 3 4]
	   results[2] = [2 3 4]
	*/
}
