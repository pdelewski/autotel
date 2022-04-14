// Copyright SumoLogic, Przemyslaw Delewskis
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

package main

import (
	"context"
	"fmt"

	otel "go.opentelemetry.io/otel"
	"sumologic.com/autotel/rtlib"
)

func main() {
	ts := rtlib.NewTracingState()
	defer func() {
		if err := ts.Tp.Shutdown(context.Background()); err != nil {
			ts.Logger.Fatal(err)
		}
	}()

	otel.SetTracerProvider(ts.Tp)

	{
		ctx := context.Background()
		_, span := otel.Tracer("fib").Start(ctx, "before Fib")
		fmt.Println(Fibonacci(10))
		span.End()
	}

}
