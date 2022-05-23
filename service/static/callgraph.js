var callgraph = {
	nodes: [
		 { data: { id: 'FibonacciHelper' } },
		 { data: { id: 'uint64' } },
		 { data: { id: 'uint' } },
		 { data: { id: 'main' } },
		 { data: { id: 'Fibonacci' } },
	],
	edges: [
		 { data: { id: 'e0', source: 'Fibonacci', target: 'uint' } },
		 { data: { id: 'e1', source: 'main', target: 'FibonacciHelper' } },
		 { data: { id: 'e2', source: 'FibonacciHelper', target: 'Fibonacci' } },
		 { data: { id: 'e3', source: 'Fibonacci', target: 'uint64' } },
	]
};