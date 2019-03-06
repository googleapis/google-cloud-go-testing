/*
Copyright 2018 Google LLC

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

/*
These packages contain code that can help you test against the GCP Client
Libraries for Go (https://github.com/GoogleCloudPlatform/google-cloud-go).

We do not recommend using mocks for most testing. Please read
https://testing.googleblog.com/2013/05/testing-on-toilet-dont-overuse-mocks.html.

Note: These packages are in alpha. Some backwards-incompatible changes may
occur.


Embedding Interfaces

All interfaces in this package include an embedToIncludeNewMethods method. This
is intentionally unexported so that any implementor of the interface must
embed the interface in their implementation. Embedding the interface in an
implementation has the effect that any future methods added to the interface
will not cause compile-time errors (the implementation does not implement
the newly-added method), since embedded interfaces provide a default method for
unimplemented methods.

See Example (RecordBuckets) for an example of how to implement interfaces
(including embedding the interface).
*/
package googlecloudgotesting
