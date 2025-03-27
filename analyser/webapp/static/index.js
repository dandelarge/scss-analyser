const COLORS = {
  selected: "#dd0b77",
  selectedSecondary: "#dd10aa",
  neutral: "#f0c987",
  neutralLink: "#ccc",
  highlightedLink: "#ee0910",
};

const fileDetails = document.getElementById("fileDetails");
const graphContainer = document.getElementById("graphContainer");
const state = {
  selected: null,
  highlightedNodes: [],
};

function colorNode(d) {
  return d.id.indexOf("scss") > -1 ? COLORS.neutral : "green";
}

function traverseLinks(links, id) {
  const children = [];
  links.forEach((link) => {
    if (link.source.id === id) {
      const tmp = traverseLinks(links, link.target.id);
      children.push(...tmp);
    }
  });
  return [id, ...children];
}

function highlightNodes(label) {
  label = label || d3.select(this).attr("id");

  const allNodes = d3.selectAll("circle");
  const allLinks = d3.selectAll("line");

  const linksData = allLinks.data().map((link) => link);

  const nodesToHighlight = traverseLinks(linksData, label);

  const nodes = allNodes.filter((node) => {
    return nodesToHighlight.includes(node.id);
  });

  const links = allLinks.filter((link) => {
    return nodesToHighlight.includes(link.source.id);
  });

  nodes.style("fill", COLORS.selected);
  links.style("stroke", COLORS.highlightedLink);
}

const svg = d3
  .select("#graph")
  .append("svg")
  .attr("width", 3000)
  .attr("height", 3000)
  .append("g")
  .attr("transform", "translate(50,50)");

d3.json("api/d3data").then(runGraph);

function runGraph(data) {
  const link = svg
    .selectAll("line")
    .data(data.links)
    .join("line")
    .style("stroke", COLORS.neutralLink);

  const node = svg
    .selectAll("circle")
    .data(data.nodes)
    .join("circle")
    .attr("r", 8)
    .attr("fill", colorNode);

  const label = svg
    .selectAll("text")
    .data(data.nodes)
    .attr("x", -100)
    .attr("y", -100)
    .join("text")
    .text((d) => d.id)
    .style("fill", "black")
    .style("font-size", 14)
    .style("font-weight", "semibold")
    .style("text-anchor", "left")
    .style("alignment-baseline", "middle");

  d3.forceSimulation(data.nodes)
    .force(
      "link",
      d3
        .forceLink()
        .id((d) => d.id)
        .links(data.links),
    )
    .force("charge", d3.forceManyBody().strength(-200))
    .force("center", d3.forceCenter(1500, 1500))
    .on("end", ticked);

  function ticked() {
    link
      .attr("x1", (d) => d.source.x)
      .attr("y1", (d) => d.source.y)
      .attr("x2", (d) => d.target.x)
      .attr("y2", (d) => d.target.y);

    node
      .attr("cx", (d) => d.x)
      .attr("cy", (d) => d.y)
      .attr("id", (d) => d.id)
      .attr("fill", colorNode)

      .on("mouseover", function () {
        const allNodes = d3.selectAll("circle");
        const allLinks = d3.selectAll("line");
        const allLabels = d3.selectAll("text");

        const selectedLabel = allLabels.filter(
          (d) => d.id === d3.select(this).attr("id"),
        );

        allNodes.style("fill", colorNode);
        allLinks.style("stroke", COLORS.neutralLink);
        allLabels.attr("x", -100).attr("y", -100);

        selectedLabel.attr("x", (d) => d.x + 15).attr("y", (d) => d.y);
        highlightNodes.call(this);
      })
      .on("mouseout", function () {
        const allNodes = d3.selectAll("circle");
        const allLinks = d3.selectAll("line");
        const allLabels = d3.selectAll("text");

        allNodes.style("fill", colorNode);
        allLinks.style("stroke", COLORS.neutralLink);
        allLabels.attr("x", -100).attr("y", -100);

        if (state.selected) {
          highlightNodes(state.selected);
        }
      })
      .on("click", function () {
        fileDetails.innerHTML = "";
        const allNodes = d3.selectAll("circle");
        const allLinks = d3.selectAll("line");
        const allLabels = d3.selectAll("text");

        allNodes.attr("fill", COLORS.neutral);
        allLinks.style("stroke", COLORS.neutralLink);
        allLabels.attr("x", -100).attr("y", -100);

        const list = document.createElement("ul");
        const title = document.createElement("h2");
        title.classList.add("text-4xl");
        const fileName = d3.select(this).attr("id");
        title.innerHTML = fileName;
        state.selected = fileName;
        state.highlightedNodes = [];

        const fileNameList = traverseLinks(data.links, fileName);
        state.highlightedNodes = fileNameList;

        fileNameList
          .reduce(
            (acc, curr) => (acc.includes(curr) ? acc : [...acc, curr]),
            [],
          )
          .forEach((file) => {
            const listItem = document.createElement("li");
            listItem.innerHTML = file;
            list.appendChild(listItem);
          });

        fileDetails.appendChild(title);
        fileDetails.appendChild(list);

        highlightNodes.call(this);
      });

    label.attr("x", -100).attr("y", -100);

    highlightNodes("containers/application.jsx");
    const resetLabel = label.filter(
      (d) => d.id === "containers/application.jsx",
    );
    resetLabel.attr("x", (d) => d.x + 15).attr("y", (d) => d.y);

    graphContainer.scroll(1200, 1200);

    console.log("Ticked");
  }

  const counts = countFiles(data);
  updateCountUI(counts);
}

function countFiles(data) {
  const fileCount = data.nodes.length;
  const scssCount = data.nodes.filter(
    (node) => node.id.indexOf("scss") > -1,
  ).length;
  const jsCount = fileCount - scssCount;
  const importCount = data.links.length;
  const scssImports = data.links.filter(
    (link) => link.source.id.indexOf("scss") > -1,
  ).length;
  const jsImports = importCount - scssImports;

  return {
    total: fileCount,
    scss: scssCount,
    js: jsCount,
    imports: importCount,
    scssImports,
    jsImports,
  };
}

function updateCountUI(counts) {
  const total = document.getElementById("totalFiles");
  const scss = document.getElementById("scssFiles");
  const js = document.getElementById("jsFiles");
  const imports = document.getElementById("imports");
  const cssImportsElement = document.getElementById("scssImports");
  const jsImportsElement = document.getElementById("jsImports");

  total.innerHTML = counts.total;
  scss.innerHTML = counts.scss;
  js.innerHTML = counts.js;
  imports.innerHTML = counts.imports;
  cssImportsElement.innerHTML = counts.scssImports;
  jsImportsElement.innerHTML = counts.jsImports;
}
