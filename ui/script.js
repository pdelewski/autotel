$(function() {
    var cy = cytoscape({
        container: document.getElementById('cy'),
        boxSelectionEnabled : true,
        style: [
            {
                selector: 'node',
                css: {
                    width: 50,
                    height: 50,
                    'background-color':'#61bffc',
                    content: 'data(id)',
                    'text-valign' : 'center',
                    'text-halign' : 'center',
                    shape: 'rectangle',
                    
                }
            },
            {
                selector: "edge",
                css: {
                    "curve-style": "bezier",
                    "control-point-step-size": 20,
                    "target-arrow-shape": "triangle"
                }
            },
            {
                selector: ':selected',
                style: {
                  'background-color': 'red',
                  'line-color': 'black',
                  'target-arrow-color': 'black',
                  'source-arrow-color': 'black',
                }
            },
        ],
        elements: callgraph,
        layout: {
            name: 'breadthfirst',
            directed: true,
            padding: 10,
            /* color: "#ffff00",*/
            fit: true
        }
    });
    cy.on('cxttap', "node", function(event) { alert("right click on node");});
    cy.on('cxttap', "edge", function(event) { alert("right click on edge");});
});
