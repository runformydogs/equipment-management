let currentTreeData = null;

document.addEventListener('DOMContentLoaded', async () => {
    try {
        const data = await api.getFullTree();
        currentTreeData = data.tree;
        renderTree(currentTreeData);
    } catch (error) {
        console.error('Error loading tree:', error);
        alert('Failed to load network tree');
    }
});

function renderTree(treeData) {
    const d3 = window.d3;
    if (!d3) {
        console.error('D3 not found!');
        return;
    }

    const width = 928;
    const marginTop = 10;
    const marginRight = 10;
    const marginBottom = 10;
    const marginLeft = 40;

    const root = d3.hierarchy({children: treeData}, d => d.children);

    const dx = 10;
    const dy = (width - marginRight - marginLeft) / (1 + root.height);

    const tree = d3.tree().nodeSize([dx, dy]);
    const diagonal = d3.linkHorizontal().x(d => d.y).y(d => d.x);

    d3.select("#tree-container").html("");

    const svg = d3.select("#tree-container")
        .append("svg")
        .attr("width", width)
        .attr("height", dx)
        .attr("viewBox", [-marginLeft, -marginTop, width, dx])
        .attr("style", "max-width: 100%; height: auto; font: 10px sans-serif; user-select: none;");

    const gLink = svg.append("g")
        .attr("fill", "none")
        .attr("stroke", "#555")
        .attr("stroke-opacity", 0.4)
        .attr("stroke-width", 1.5);

    const gNode = svg.append("g")
        .attr("cursor", "pointer")
        .attr("pointer-events", "all");

    root.x0 = dy / 2;
    root.y0 = 0;

    root.descendants().forEach((d, i) => {
        d.id = i;
        d._children = d.children;

        // Сворачиваем все узлы кроме корневых
        if (d.depth > 0 && d.data.type === "node") {
            d.children = null;
        }
    });

    update(root);

    function update(source) {
        const duration = 250;
        const nodes = root.descendants().reverse();
        const links = root.links();

        tree(root);

        let left = root;
        let right = root;
        root.eachBefore(node => {
            if (node.x < left.x) left = node;
            if (node.x > right.x) right = node;
        });

        const height = right.x - left.x + marginTop + marginBottom;

        const transition = svg.transition()
            .duration(duration)
            .attr("height", height)
            .attr("viewBox", [-marginLeft, left.x - marginTop, width, height])
            .tween("resize", window.ResizeObserver ? null : () => () => svg.dispatch("toggle"));

        const node = gNode.selectAll("g")
            .data(nodes, d => d.id);

        const nodeEnter = node.enter().append("g")
            .attr("transform", d => `translate(${source.y0},${source.x0})`)
            .attr("fill-opacity", 0)
            .attr("stroke-opacity", 0)
            .on("click", (event, d) => {
                if (d.data.type === "node") {
                    d.children = d.children ? null : d._children;
                    update(d);
                }
            });

        nodeEnter.append("circle")
            .attr("r", d => d.data.type === "node" ? 4 : 2.5)
            .attr("fill", d => {
                if (d.data.type === "device") return "#999";
                return d._children ? "#555" : "#999";
            })
            .attr("stroke-width", 10);

        nodeEnter.append("text")
            .attr("dy", "0.31em")
            .attr("x", d => d._children ? -8 : 8)
            .attr("text-anchor", d => d._children ? "end" : "start")
            .text(d => d.data.name)
            .attr("stroke-linejoin", "round")
            .attr("stroke-width", 3)
            .attr("stroke", "white")
            .attr("paint-order", "stroke");

        const nodeUpdate = node.merge(nodeEnter).transition(transition)
            .attr("transform", d => `translate(${d.y},${d.x})`)
            .attr("fill-opacity", 1)
            .attr("stroke-opacity", 1);

        const nodeExit = node.exit().transition(transition).remove()
            .attr("transform", d => `translate(${source.y},${source.x})`)
            .attr("fill-opacity", 0)
            .attr("stroke-opacity", 0);

        const link = gLink.selectAll("path")
            .data(links, d => d.target.id);

        const linkEnter = link.enter().append("path")
            .attr("d", d => {
                const o = {x: source.x0, y: source.y0};
                return diagonal({source: o, target: o});
            });

        link.merge(linkEnter).transition(transition)
            .attr("d", diagonal);

        link.exit().transition(transition).remove()
            .attr("d", d => {
                const o = {x: source.x, y: source.y};
                return diagonal({source: o, target: o});
            });

        root.eachBefore(d => {
            d.x0 = d.x;
            d.y0 = d.y;
        });
    }
}

window.renderTree = renderTree;