/*
Copyright 2017 Heptio Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sinks

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// captureStdout runs fn with os.Stdout redirected to a pipe and returns what it printed.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Stdout = old
		// Close the write end even if fn panicked, so the drain
		// goroutine below sees EOF and cannot leak.
		w.Close()
	}()
	// Drain the pipe in the background so fn can write more than the
	// pipe buffer without deadlocking.
	done := make(chan string, 1)
	go func() {
		b, _ := io.ReadAll(r)
		done <- string(b)
	}()
	os.Stdout = w
	fn()
	w.Close()
	return <-done
}

// A created event is serialized to stdout as JSON with verb ADDED and no old_event.
func TestStdoutSinkUpdateEventsAdded(t *testing.T) {
	eNew := &v1.Event{
		ObjectMeta:     metav1.ObjectMeta{Name: "e1", ResourceVersion: "1"},
		InvolvedObject: v1.ObjectReference{Kind: "Pod", Name: "foo"},
	}
	out := captureStdout(t, func() { NewStdoutSink().UpdateEvents(eNew, nil) })

	var data EventData
	if err := json.Unmarshal([]byte(out), &data); err != nil {
		t.Fatalf("stdout output is not valid JSON: %v (%q)", err, out)
	}
	if data.Verb != "ADDED" {
		t.Errorf("verb = %q, want ADDED", data.Verb)
	}
	if data.Event == nil || data.Event.Name != "e1" {
		t.Errorf("event name = %v, want e1", data.Event)
	}
	if data.OldEvent != nil {
		t.Errorf("old_event = %v, want nil", data.OldEvent)
	}
}

// captureStdout must restore os.Stdout even if fn panics, so a failing
// helper call cannot silence the rest of the test process's output.
func TestCaptureStdoutRestoresStdoutOnPanic(t *testing.T) {
	before := os.Stdout
	func() {
		defer func() { _ = recover() }()
		captureStdout(t, func() { panic("boom") })
	}()
	if os.Stdout != before {
		t.Fatal("captureStdout did not restore os.Stdout after a panic in fn")
	}
}

// An updated event is serialized with verb UPDATED and the previous event in old_event.
func TestStdoutSinkUpdateEventsUpdated(t *testing.T) {
	eOld := &v1.Event{
		ObjectMeta: metav1.ObjectMeta{Name: "e1", ResourceVersion: "1"},
	}
	eNew := &v1.Event{
		ObjectMeta: metav1.ObjectMeta{Name: "e1", ResourceVersion: "2"},
	}
	out := captureStdout(t, func() { NewStdoutSink().UpdateEvents(eNew, eOld) })

	var data EventData
	if err := json.Unmarshal([]byte(out), &data); err != nil {
		t.Fatalf("stdout output is not valid JSON: %v (%q)", err, out)
	}
	if data.Verb != "UPDATED" {
		t.Errorf("verb = %q, want UPDATED", data.Verb)
	}
	if data.OldEvent == nil {
		t.Error("old_event missing, want the previous event")
	}
}
