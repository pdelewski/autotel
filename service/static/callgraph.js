var callgraph = {
	nodes: [
		 { data: { id: 'Println' } },
		 { data: { id: 'uint64' } },
		 { data: { id: 'SumoAutoInstrument' } },
		 { data: { id: 'main' } },
		 { data: { id: 'FibonacciHelper' } },
		 { data: { id: 'Errorf' } },
		 { data: { id: 'uint' } },
		 { data: { id: 'Fibonacci' } },
		 { data: { id: 'foo' } },
	],
	edges: [
		 { data: { id: 'e0', source: 'Fibonacci', target: 'uint64' } },
		 { data: { id: 'e1', source: 'Fibonacci', target: 'Errorf' } },
		 { data: { id: 'e2', source: 'main', target: 'FibonacciHelper' } },
		 { data: { id: 'e3', source: 'main', target: 'SumoAutoInstrument' } },
		 { data: { id: 'e4', source: 'Fibonacci', target: 'uint' } },
		 { data: { id: 'e5', source: 'foo', target: 'Println' } },
		 { data: { id: 'e6', source: 'main', target: 'Println' } },
		 { data: { id: 'e7', source: 'FibonacciHelper', target: 'foo' } },
		 { data: { id: 'e8', source: 'FibonacciHelper', target: 'Fibonacci' } },
	]
};