// Package snapshot provides point-in-time capture and comparison of driftctl
// scan results.
//
// A Store persists each scan result as a timestamped JSON file in a configurable
// directory. Entries can be loaded individually or listed in chronological order.
//
// The Compare function calculates a Delta between any two entries, exposing the
// change in drifted resource count as well as the sets of newly-drifted and
// resolved resources. This enables regression detection and audit trails across
// multiple CI runs.
//
// Example usage:
//
//	store, _ := snapshot.NewStore("/var/lib/driftctl/snapshots")
//	entry, _ := store.Save(result, "ci-run-42")
//	entries, _ := store.List()
//	delta := snapshot.Compare(entries[0], entries[len(entries)-1])
//	fmt.Println(delta.Summary())
package snapshot
