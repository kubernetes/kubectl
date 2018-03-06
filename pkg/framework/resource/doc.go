/*
Copyright 2018 The Kubernetes Authors.

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

Package resource implements tools for the discovery and filtering
of resources in an API server.

The starting point should be the creation of a new Parser object
with NewParser(). This can then use the Resources() method to discover
resources in the API server.

Filtering can also be applied by implementing the Filter interface.

See examples below, starting with the Resources one.

*/
package resource
