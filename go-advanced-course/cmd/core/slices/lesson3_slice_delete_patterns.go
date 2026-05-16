package main

import "fmt"

/*
LESSON 3: Slice Deletion Patterns

Goal:
Remove an element from a slice.

There are TWO important patterns:

1) Stable Delete
   - preserves order
   - shifts elements left
   - O(n)

2) Unstable Delete
   - does NOT preserve order
   - swaps with last element
   - O(1)

Choosing the right one can drastically
improve performance in large datasets.
*/

type Job struct {
	ID   int
	Name string
}

/*
STABALE DELETE

Keeps order intact.

Example:
[1 2 3 4 5] remove index 2

Result:
[1 2 4 5]

Internally
*/
func stableDelete(s []Job, index int) []Job {
	return append(s[:index], s[index+1:]...)
}

/*
UNSTABLE DELETE

Swap element with the last element
and shrink slice.

Example:

[1 2 3 4 5] remove index 1

Step 1 swap with last
[1 5 3 4 2]

Step 2 shrink
[1 5 3 4]

Cost: O(1)
*/
func unstableDelete(s []Job, index int) []Job {
	last := len(s) - 1

	// move last element into deleted spot
	s[index] = s[last]

	// Shrink slice
	return s[:last]
}

func printJobs(label string, jobs []Job) {

	fmt.Println(label)

	for _, j := range jobs {
		fmt.Printf("JobID=%d Name=%s\n", j.ID, j.Name)
	}

	fmt.Println()
}

func Slice_Delete_Patters() {
	jobs := []Job{
		{1, "email"},
		{2, "image-process"},
		{3, "analytics"},
		{4, "report"},
		{5, "backup"},
	}

	printJobs("Original Jobs", jobs)

	/*
		STABLE DELETE
	*/
	stable := make([]Job, len(jobs))
	copy(stable, jobs)

	stable = stableDelete(stable, 1)

	printJobs("After Stable Delete (remove index 1)", stable)

	/*
		UNSTABLE DELETE
	*/
	unstable := make([]Job, len(jobs))
	copy(unstable, jobs)

	unstable = unstableDelete(unstable, 1)

	printJobs("After Unstable Delete (remove index 1)", unstable)
}
