// Copyright 2025
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package record

import (
	"strconv"
	"sync"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

var (
	initOnce        sync.Once
	defaultRecorder record.EventRecorder
)

func init() {
	defaultRecorder = new(record.FakeRecorder)
}

// InitFromRecorder initializes the global default recorder. It can only be called once.
// Subsequent calls are considered noops.
func InitFromRecorder(recorder record.EventRecorder) {
	initOnce.Do(func() {
		defaultRecorder = recorder
	})
}

// Event constructs an event from the given information and puts it in the queue for sending.
func Event(object runtime.Object, generation int64, reason, message string) {
	defaultRecorder.AnnotatedEventf(object, getEventsAnnotations(generation), corev1.EventTypeNormal, title(reason), message)
}

// Eventf is just like Event, but with Sprintf for the message field.
func Eventf(object runtime.Object, generation int64, reason, message string, args ...any) {
	defaultRecorder.AnnotatedEventf(object, getEventsAnnotations(generation), corev1.EventTypeNormal, title(reason), message, args...)
}

// Warn constructs a warning event from the given information and puts it in the queue for sending.
func Warn(object runtime.Object, generation int64, reason, message string) {
	defaultRecorder.AnnotatedEventf(object, getEventsAnnotations(generation), corev1.EventTypeWarning, title(reason), message)
}

// Warnf is just like Warn, but with Sprintf for the message field.
func Warnf(object runtime.Object, generation int64, reason, message string, args ...any) {
	defaultRecorder.AnnotatedEventf(object, getEventsAnnotations(generation), corev1.EventTypeWarning, title(reason), message, args...)
}

func getEventsAnnotations(generation int64) map[string]string {
	if generation == 0 {
		return nil
	}
	return map[string]string{
		"generation": strconv.FormatInt(generation, 10),
	}
}

// title returns a copy of the string s with all Unicode letters that begin words
// mapped to their Unicode title case.
func title(source string) string {
	return cases.Title(language.Und, cases.NoLower).String(source)
}
