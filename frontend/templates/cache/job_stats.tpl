<?php $_v=&$this->vars; $this->push();$this->load("html_header.tpl");$this->assign($_v);$this->render();$this->pop();?>


<style>
input.form-control::placeholder {
	opacity: .3;
}
</style>

<script src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/2.4.0/Chart.min.js"></script>

<div class="container">
	<div class="row">
		<div class="col-8">
			<h2><?php echo $_v['title'];?>
</h2>
		</div>
		<div class="col-4 tar">
			<a class="btn btn-secondary" href="/">home</a>
			<a class="btn btn-danger" href="/run/<?php echo $_v['job']['jobName'];?>
">Run again</a>
		</div>
	</div>

	<section class="contents">
		<b>Daily stats</b>
		<div class="wrapper">
			<canvas id="chart-daily"></canvas>
		</div>
		<b>Monthly stats</b>
		<div class="wrapper">
			<canvas id="chart-monthly"></canvas>
		</div>
	</section>

	<section class="contents">
		<div class="input-group mb-3">
			<input id="searchTextbox" type="text" class="form-control" placeholder="Search string" aria-label="Search string">
			<div class="input-group-append">
				<button id="searchButton" class="btn btn-outline-secondary" type="button" onclick="getSearch()">Search</button>
			</div>
			<div class="input-group-append">
				<button class="btn btn-outline-secondary" type="button" onclick="reload()">Reload</button>
			</div>
		</div>
		<div class="alert alert-danger" role="alert">
			Use <b>Search</b> with caution, major performance impact can occur!
		</div>

		<table class="table table-striped table-sm">
		<thead>
			<tr>
				<th id="searchTitle" colspan="4"></th>
			</tr>
			<tr>
				<th>ID</th>
				<th>Exit code</th>
				<th>Duration</th>
				<th>Timestamp</th>
			</tr>
		</thead>
		<tbody id="logresults">
		<?php if(!empty($_v['logs']))foreach($_v['logs'] as $_v['log']){?>
			<tr>
				<td><a href="<?php echo $_v['log']['link'];?>
"><?php echo $_v['log']['id'];?>
</a></td>
				<td><?php echo $_v['log']['exitCode'];?>
</td>
				<td><?php echo $_v['log']['duration'];?>
</td>
				<td><?php echo $_v['log']['date'];?>
</td>
			</tr>
		<?php }else{?>

			<tr><td colspan="4">No logs available.</td></tr>
		<?php }?>
		</tbody>
		</table>
	</section>

	<section class="contents">
		<div class="row">
			<div class="col-6">
				<a id="linkPrevious" class="btn btn-primary" href="<?php echo $_v['job']['link'];?>
?pageNumber=<?php echo $_v['job']['pageNumber']-1;?>
" <?php if($_v['job']['pageNumber'] == 0){?>
onclick="return false;"<?php }?>>previous</a>
				<a id="linkNext" class="btn btn-primary" href="<?php echo $_v['job']['link'];?>
?pageNumber=<?php echo $_v['job']['pageNumber']+1;?>
">next</a>
			</div>
			<div class="col-6 tar">
				<a class="btn btn-secondary" href="/">home</a>
			</div>
	</section>
</div>

<style>
#chart-daily,
#chart-monthly {
	height: 224px !important;
}
</style>

<script>
var ctx = document.getElementById('chart-daily').getContext('2d')
new Chart(ctx, {
	type: 'bar',
	data: <?php echo json_encode($_v['daily']);?>
,
	options: {
		tooltips: {
			displayColors: true,
			callbacks: {
				mode: 'x',
			},
		},
		scales: {
			xAxes: [{
				stacked: false,
				gridLines: {
					display: false,
				}
			}],
			yAxes: [{
				stacked: false,
				ticks: {
					beginAtZero: true,
				},
				type: 'linear',
			}]
		},
		// responsive: true,
		maintainAspectRatio: false,
		legend: {
			position: 'bottom',
		},
	}
})

var ctx2 = document.getElementById('chart-monthly').getContext('2d')
new Chart(ctx2, {
	type: 'bar',
	data: <?php echo json_encode($_v['monthly']);?>
,
	options: {
		tooltips: {
			displayColors: true,
			callbacks: {
				mode: 'x',
			},
		},
		scales: {
			xAxes: [{
				stacked: false,
				gridLines: {
					display: false,
				}
			}],
			yAxes: [{
				stacked: false,
				ticks: {
					beginAtZero: true,
				},
				type: 'linear',
			}]
		},
		// responsive: true,
		maintainAspectRatio: false,
		legend: {
			position: 'bottom',
		},
	}
})
</script>

<script>
var searchTextbox = document.getElementById("searchTextbox");
var searchButton = document.getElementById("searchButton");
searchTextbox.addEventListener("keyup", function(event) {
  if (event.keyCode === 13) {
    searchButton.click();
  }
  if (searchTextbox.value.length == 0) {
	console.log("Empty searchTextbox");
	var linkPrevious = document.getElementById("linkPrevious");
	linkPrevious.href = "<?php echo $_v['job']['link'];?>
?pageNumber=<?php echo $_v['job']['pageNumber']-1;?>
";

	var linkNext = document.getElementById("linkNext");
	linkNext.href="<?php echo $_v['job']['link'];?>
?pageNumber=<?php echo $_v['job']['pageNumber']+1;?>
";
  }
});

// search
var searchPage = 0;
function getSearch(pageNumber = 0) {
	var formData = new FormData();
	var searchTextbox = document.getElementById("searchTextbox");
	var searchButton = document.getElementById("searchButton");

	formData.append("jobName", "<?php echo $_v['job']['jobName'];?>
");
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
					ExitCode.innerHTML = '<span class="badge badge-danger">Exit: ' + element.exitCode + '</span>';
				} else {
					ExitCode.innerHTML = '<span class="badge badge-success">OK</span>';
				}
				row.appendChild(ExitCode);

				var Duration = document.createElement("td");
				Duration.innerText = parseFloat(element.duration).toFixed(3) + "s";
				row.appendChild(Duration);

				var Timestamp = document.createElement("td");
				var t = new Date(element.stamp);
				Timestamp.innerText = t.toISOString().slice(0, 19).replace('T', ' ');
				row.appendChild(Timestamp);

				tbody.appendChild(row);
			});

			// make pagination
			var linkPrevious = document.getElementById("linkPrevious");
			linkPrevious.href = "javascript:showSearchPage(" + (searchPage-1) + ")";
			//linkPrevious.addEventListener("click", showSearchPage(searchPage-1));

			var linkNext = document.getElementById("linkNext");
			linkNext.href = "javascript:showSearchPage(" + (searchPage+1) + ")";
			//linkNext.addEventListener("click", showSearchPage(searchPage+1));
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

<?php $this->push();$this->load("html_footer.tpl");$this->assign($_v);$this->render();$this->pop();?>
