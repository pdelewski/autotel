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
    cy.on('tap', "node", function(event) { 
      var typeIds = cy.elements('node:selected');
      alert(typeIds[1].id());
    });
    
    //cy.on('cxttap', "edge", function(event) { alert("right click on edge");});
    var n1 = cy.$('#main').successors().nodes().size();
    cy.$('#main').select();
    var s = cy.$('#main').successors().nodes().select();
    do {
        n1 = n1.successors().nodes().size();
        s = s.successors().nodes().select();
    } while (n1 > 0)
    //alert(cy.$('#Fibonacci').successors().nodes().size());
});
