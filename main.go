// Encoding: utf-8
// Cloud Foundry Java Buildpack
// Copyright (c) 2015-2017 the original author or authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry/java-buildpack-memory-calculator/flags"
	"github.com/cloudfoundry/java-buildpack-memory-calculator/memory"
)

const (
	exec_name = "java-buildpack-memory-calculator"
)

func main() {
	// validateFlags() will exit on error
	jreType := flags.ValidateJreType() 

	if jreType == "IBM" {
		memSize, heapRatio := flags.ValidateFlagsForIBM()
		fmt.Fprint(os.Stdout, "-Xmx", memSize.Scale(heapRatio).String())
	} else {
		memSize, numThreads, numLoadedClasses, poolType, rawVmOptions := flags.ValidateFlags()
		// default the number of threads if it was not supplied
		if numThreads == 0 {
			numThreads = 50
		}

		vmOptions, err := memory.NewVmOptions(rawVmOptions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Problem with vmOptions: %s\n", err)
			os.Exit(1)
		}

		allocator, err := memory.NewAllocator(poolType, vmOptions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot allocate JVM memory configuration: %s\n", err)
			os.Exit(1)
		}

		if err = allocator.Calculate(numLoadedClasses, numThreads, memSize); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot calculate JVM memory configuration: %s\n", err)
			os.Exit(1)
		}

		// Print outputs to standard output for consumption by the caller
		fmt.Fprint(os.Stdout, allocator.String())
	}

	
}
