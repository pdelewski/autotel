var callgraph = {
	nodes: [
		 { data: { id: 'main' } },
		 { data: { id: 'Fibonacci' } },
		 { data: { id: 'FibonacciHelper' } },
		 { data: { id: 'uint64' } },
		 { data: { id: 'uint' } },
	],
	edges: [
		 { data: { id: 'e0', source: 'Fibonacci', target: 'uint64' } },
		 { data: { id: 'e1', source: 'Fibonacci', target: 'uint' } },
		 { data: { id: 'e2', source: 'main', target: 'FibonacciHelper' } },
		 { data: { id: 'e3', source: 'FibonacciHelper', target: 'Fibonacci' } },
	]
};