{load html_header.tpl}

<script src="https://cdn.rawgit.com/fnando/sparkline/master/dist/sparkline.js"></script>

<script>
function findClosest(target, tagName) {
  if (target.tagName === tagName) {
    return target;
  }

  while ((target = target.parentNode)) {
    if (target.tagName === tagName) {
      break;
    }
  }

  return target;
}

var options = {
  onmousemove(event, datapoint) {
    var svg = findClosest(event.target, "svg");
    var tooltip = svg.nextElementSibling;

    tooltip.hidden = false;
    tooltip.textContent = datapoint.name + ': ' + (datapoint.value.toFixed(3)) + 's'
    tooltip.style.top = event.offsetY + "px";
    tooltip.style.left = (event.offsetX+20) + "px";
  },

  onmouseout() {
    var svg = findClosest(event.target, "svg");
    var tooltip = svg.nextElementSibling;

    tooltip.hidden = true;
  }
};

</script>

<div class="container">

	<section class="contents">
		<b>Onetime Jobs</b>
		<table class="table table-striped table-sm">
		<thead>
			<tr>
				<th>Script</th>
				<th>History</th>
				<th class="tar">Duration</th>
				<th class="tar">Last timestamp</th>
			</tr>
		</thead>
		<tbody>

		{foreach $body['jobs'] as $k => $job}
			{if $job['active'] == 1 && strpos($job['name'], "onetime/") !== false}
	<tr>
		<td><a href="/{$job['name']}">{$job['name']}</a></td>

				{if !empty($job['lastRun'])}
			<td>

      <div style="position: relative">
        <svg class="job-{$k}" width="300" height="20" stroke-width="2" {if $job['lastRun']['exitCode'] == 0}stroke="blue" fill="rgba(0, 0, 255, .2)"{else}stroke="red" fill="rgba(255, 0, 0, .2)"{/if}></svg>
        <span class="tooltip" hidden="true"></span>
      </div>

				<script>
					sparkline.sparkline(document.querySelector(".job-{$k}"), {$job['history']|history}, options)
				</script>
			</td>
			<td class="tar">{eval echo sprintf("%.3fs", $job['lastRun']['duration'])}</td>
			<td class="tar">{eval echo date("Y-m-d H:i:s", strtotime($job['lastRun']['stamp']))}</td>
				{else}
			<td colspan="4"><i>No last run info</i></td>
				{/if}
	</tr>
			{/if}
		{/foreach}

		</tbody>
		</table>
	</section>

	<section class="contents">
		<b>Scheduled Jobs</b>
		<table class="table table-striped table-sm">
		<thead>
			<tr>
				<th>Script</th>
				<th>History</th>
				<th class="tar">Duration</th>
				<th class="tar">Last timestamp</th>
			</tr>
		</thead>
		<tbody>

		{foreach $body['jobs'] as $k => $job}
			{if $job['active'] == 1 && strpos($job['name'], "onetime/") === false}
	<tr>
		<td><a href="/{$job['name']}">{$job['name']}</a></td>

				{if !empty($job['lastRun'])}
			<td>

      <div style="position: relative">
        <svg class="job-{$k}" width="300" height="20" stroke-width="2" {if $job['lastRun']['exitCode'] == 0}stroke="blue" fill="rgba(0, 0, 255, .2)"{else}stroke="red" fill="rgba(255, 0, 0, .2)"{/if}></svg>
        <span class="tooltip" hidden="true"></span>
      </div>

				<script>
					sparkline.sparkline(document.querySelector(".job-{$k}"), {$job['history']|history}, options)
				</script>
			</td>
			<td class="tar">{eval echo sprintf("%.3fs", $job['lastRun']['duration'])}</td>
			<td class="tar">{eval echo date("Y-m-d H:i:s", strtotime($job['lastRun']['stamp']))}</td>
				{else}
			<td colspan="4"><i>No last run info</i></td>
				{/if}
	</tr>
			{/if}
		{/foreach}

		</tbody>
		</table>
	</section>
</div>

<style>
*[hidden] {
  display: none;
}

.tooltip {
  position: absolute;
  background: rgba(0, 0, 0, .7);
  color: #fff;
  padding: 2px 5px;
  font-size: 12px;
  white-space: nowrap;
  z-index: 9999;
  opacity: 1;
}

.sparkline--cursor {
  stroke: orange;
}

.sparkline--spot {
  fill: red;
  stroke: red;
}
</style>

{load html_footer.tpl}
