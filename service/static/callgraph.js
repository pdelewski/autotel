var callgraph = {
	nodes: [
		 { data: { id: 'uint64' } },
		 { data: { id: 'uint' } },
		 { data: { id: 'main' } },
		 { data: { id: 'Fibonacci' } },
		 { data: { id: 'FibonacciHelper' } },
	],
	edges: [
		 { data: { id: 'e0', source: 'main', target: 'FibonacciHelper' } },
		 { data: { id: 'e1', source: 'FibonacciHelper', target: 'Fibonacci' } },
		 { data: { id: 'e2', source: 'Fibonacci', target: 'uint64' } },
		 { data: { id: 'e3', source: 'Fibonacci', target: 'uint' } },
	]
};