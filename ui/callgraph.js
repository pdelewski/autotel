var callgraph = {
    nodes: [
        { data: { id: 'main' } },
        { data: { id: 'FibonacciHelper' } },
        { data: { id: 'Fibonacci' } }
        

    ],
    edges: [
        { data: { id: 'e1', source: 'main', target: 'FibonacciHelper' } },
        { data: { id: 'e2', source: 'FibonacciHelper', target: 'Fibonacci' } }                      
    ]
};