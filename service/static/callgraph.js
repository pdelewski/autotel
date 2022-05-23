var callgraph = {
	nodes: [
		 { data: { id: 'FibonacciHelper' } },
		 { data: { id: 'uint64' } },
		 { data: { id: 'uint' } },
		 { data: { id: 'main' } },
		 { data: { id: 'Fibonacci' } },
	],
	edges: [
		 { data: { id: 'e0', source: 'FibonacciHelper', target: 'Fibonacci' } },
		 { data: { id: 'e1', source: 'Fibonacci', target: 'uint64' } },
		 { data: { id: 'e2', source: 'Fibonacci', target: 'uint' } },
		 { data: { id: 'e3', source: 'main', target: 'FibonacciHelper' } },
	]
};