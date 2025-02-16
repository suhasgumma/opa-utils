package reportsummary

import (
	"github.com/kubescape/opa-utils/reporthandling/apis"
	helpersv1 "github.com/kubescape/opa-utils/reporthandling/helpers/v1"
)

type ICounters interface {
	Excluded() int
	Passed() int
	Skipped() int
	Failed() int
	All() int
}

// =================================== Counters ============================================

// NumberOfExcluded get the number of excluded resources
func (resourceCounters *ResourceCounters) Set(allLists *helpersv1.AllLists) {
	resourceCounters.ExcludedResources = len(allLists.Excluded())
	resourceCounters.FailedResources = len(allLists.Failed())
	resourceCounters.PassedResources = len(allLists.Passed())
}

// NumberOfExcluded get the number of excluded resources
func (resourceCounters *ResourceCounters) Excluded() int {
	return resourceCounters.ExcludedResources
}

// NumberOfPassed get the number of passed resources
func (resourceCounters *ResourceCounters) Passed() int {
	return resourceCounters.PassedResources
}

// NumberOfSkipped get the number of skipped resources
func (resourceCounters *ResourceCounters) Skipped() int {
	return 0
}

// NumberOfFailed get the number of failed resources
func (resourceCounters *ResourceCounters) Failed() int {
	return resourceCounters.FailedResources
}

// NumberOfAll get the number of all resources
func (resourceCounters *ResourceCounters) All() int {
	return resourceCounters.Excluded() + resourceCounters.Failed() + resourceCounters.Skipped() + resourceCounters.Passed()
}

// =================================== Setters ============================================

// Increase increases the counter based on the status
func (resourceCounters *ResourceCounters) Increase(status apis.IStatus) {
	switch status.Status() {
	case apis.StatusExcluded:
		resourceCounters.ExcludedResources++
	case apis.StatusFailed:
		resourceCounters.FailedResources++
	case apis.StatusPassed:
		resourceCounters.PassedResources++
	}
}
