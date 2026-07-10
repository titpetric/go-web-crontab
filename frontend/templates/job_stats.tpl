{load html_header.tpl}

<script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.1/dist/chart.umd.min.js"></script>

<div class="container">

	<div class="page-head">
		<div>
			<h2>{title|escape}</h2>
			<div class="subtitle">Run duration trends and log history</div>
		</div>
		<div class="page-actions">
			<a class="btn btn-secondary" href="/">Dashboard</a>
			<a class="btn btn-danger" href="/run/{job.jobName}">Run again</a>
		</div>
	</div>

	<section class="panel">
		<div class="panel__head">
			<h2 class="panel__title">Duration statistics</h2>
		</div>
		<div class="panel__body--pad">
			<div class="chart-grid">
				<div class="chart-card">
					<h3>Last 24 hours</h3>
					<div class="chart-holder"><canvas id="chart-daily"></canvas></div>
				</div>
				<div class="chart-card">
					<h3>Last 30 days</h3>
					<div class="chart-holder"><canvas id="chart-monthly"></canvas></div>
				</div>
			</div>
		</div>
	</section>

	<section class="panel">
		<div class="panel__head">
			<h2 class="panel__title">Run log</h2>
		</div>
		<div class="panel__body--pad">
			<div class="toolbar">
				<input id="searchTextbox" type="text" class="form-control" placeholder="Search output string" aria-label="Search string">
				<button id="searchButton" class="btn btn-primary" type="button" onclick="getSearch()">Search</button>
				<button class="btn btn-secondary" type="button" onclick="reload()">Reload</button>
			</div>
			<div class="alert alert-warning" role="alert">
				<b>Search</b> scans raw output and can be slow on large log tables — use sparingly.
			</div>
		</div>
		<div class="panel__body">
			<table class="table">
			<thead>
				<tr>
					<th id="searchTitle" colspan="4"></th>
				</tr>
				<tr>
					<th>ID</th>
					<th>Status</th>
					<th class="tar">Duration</th>
					<th class="tar">Timestamp</th>
				</tr>
			</thead>
			<tbody id="logresults">
			{foreach $logs as $log}
				<tr>
					<td><a href="{log.link}">{log.id}</a></td>
					<td>{log.exitCode}</td>
					<td class="tar mono">{log.duration}</td>
					<td class="tar mono">{log.date}</td>
				</tr>
			{else}
				<tr><td colspan="4" class="empty-cell">No logs available.</td></tr>
			{/foreach}
			</tbody>
			</table>
		</div>
		<div class="pager">
			<div class="page-actions">
				<a id="linkPrevious" class="btn btn-secondary" href="{job.link}?pageNumber={$job['pageNumber']-1}" {if $job['pageNumber'] == 0}onclick="return false;"{/if}>&larr; Previous</a>
				<a id="linkNext" class="btn btn-secondary" href="{job.link}?pageNumber={$job['pageNumber']+1}">Next &rarr;</a>
			</div>
		</div>
	</section>
</div>

<script>
var barOptions = {
	maintainAspectRatio: false,
	responsive: true,
	scales: {
		x: { grid: { display: false } },
		y: { beginAtZero: true, grid: { color: 'rgba(148,163,184,0.2)' } }
	},
	plugins: {
		legend: { position: 'bottom', labels: { boxWidth: 12, usePointStyle: true } },
		tooltip: { displayColors: true }
	}
};

new Chart(document.getElementById('chart-daily').getContext('2d'), {
	type: 'bar',
	data: {daily|json_encode},
	options: barOptions
});

new Chart(document.getElementById('chart-monthly').getContext('2d'), {
	type: 'bar',
	data: {monthly|json_encode},
	options: barOptions
});
</script>

<script>
var searchTextbox = document.getElementById("searchTextbox");
var searchButton = document.getElementById("searchButton");
searchTextbox.addEventListener("keyup", function(event) {
  if (event.keyCode === 13) {
    searchButton.click();
  }
  if (searchTextbox.value.length == 0) {
	var linkPrevious = document.getElementById("linkPrevious");
	linkPrevious.href = "{job.link}?pageNumber={$job['pageNumber']-1}";

	var linkNext = document.getElementById("linkNext");
	linkNext.href="{job.link}?pageNumber={$job['pageNumber']+1}";
  }
});

// search
var searchPage = 0;
function getSearch(pageNumber = 0) {
	var formData = new FormData();
	var searchTextbox = document.getElementById("searchTextbox");
	var searchButton = document.getElementById("searchButton");

	formData.append("jobName", "{$job['jobName']}");
	formData.append("searchStr", searchTextbox.value);
	formData.append("pageNumber", pageNumber);

	var xhr = new XMLHttpRequest();
	var url = window.location.origin + "/search";

	searchTextbox.disabled = true;
	searchButton.disabled = true;

	xhr.open("POST", url);
	xhr.send(formData);

	xhr.onreadystatechange = function() {
		var searchTitle = document.getElementById('searchTitle');
		searchTitle.innerHTML = "SEARCH FOR: <i>" + searchTextbox.value + "</i>";

		if (xhr.readyState === 4) { // DONE
			var data = JSON.parse(xhr.response);

			var tbody = document.getElementById('logresults');
			var trows = tbody.getElementsByTagName('tr');
			var rowCount = trows.length;

			// empty table
			for (var x=rowCount-1; x>=0; x--) {
				tbody.removeChild(trows[x]);
			}

			// fill new rows
			data.logs.forEach(element => {
				var row = document.createElement("tr");
				var ID = document.createElement("td");
				ID.innerHTML = '<a href="' + element.link + '/' + element.ID + '">' + element.ID + '</a>';
				row.appendChild(ID);

				var ExitCode = document.createElement("td");
				if (element.exitCode != 0) {
					ExitCode.innerHTML = '<span class="badge badge-fail">Exit: ' + element.exitCode + '</span>';
				} else {
					ExitCode.innerHTML = '<span class="badge badge-ok">OK</span>';
				}
				row.appendChild(ExitCode);

				var Duration = document.createElement("td");
				Duration.className = "tar mono";
				Duration.innerText = parseFloat(element.duration).toFixed(3) + "s";
				row.appendChild(Duration);

				var Timestamp = document.createElement("td");
				Timestamp.className = "tar mono";
				var t = new Date(element.stamp);
				Timestamp.innerText = t.toISOString().slice(0, 19).replace('T', ' ');
				row.appendChild(Timestamp);

				tbody.appendChild(row);
			});

			// make pagination
			var linkPrevious = document.getElementById("linkPrevious");
			linkPrevious.href = "javascript:showSearchPage(" + (searchPage-1) + ")";

			var linkNext = document.getElementById("linkNext");
			linkNext.href = "javascript:showSearchPage(" + (searchPage+1) + ")";
		}
		// enable SEARCH button
		searchTextbox.disabled = false;
		searchButton.disabled = false;
	}
}

function showSearchPage(pageNumber) {
	if (pageNumber <= 0) {
		pageNumber = 0;
	}
	searchPage = pageNumber;
	getSearch(pageNumber);
}

function reload() {
	location.reload();
}
</script>

{load html_footer.tpl}
