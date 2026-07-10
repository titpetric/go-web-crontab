{load html_header.tpl}

<script src="/assets/sparkline.js"></script>
<script>
var sparkOptions = {
	onmousemove: function (event, datapoint) {
		var svg = event.currentTarget;
		var tip = svg.parentNode.querySelector('.spark-tooltip');
		if (!tip || !datapoint) { return; }
		tip.hidden = false;
		tip.textContent = datapoint.name + '  ' + Number(datapoint.value).toFixed(3) + 's';
		tip.style.top = (event.offsetY - 8) + 'px';
		tip.style.left = (event.offsetX + 14) + 'px';
	},
	onmouseout: function (event) {
		var tip = event.currentTarget.parentNode.querySelector('.spark-tooltip');
		if (tip) { tip.hidden = true; }
	}
};
</script>

<div class="container">

	<div class="page-head">
		<div>
			<h1>Dashboard</h1>
			<div class="subtitle">Scheduled and one-time jobs, with recent run history.</div>
		</div>
	</div>

	<div class="stat-grid">
		<div class="stat-card">
			<div class="stat-label">Total jobs</div>
			<div class="stat-value">{summary.total}</div>
			<div class="stat-hint">tracked in crontab</div>
		</div>
		<div class="stat-card">
			<div class="stat-label">Scheduled</div>
			<div class="stat-value">{summary.scheduled}</div>
			<div class="stat-hint">{summary.onetime} one-time</div>
		</div>
		<div class="stat-card is-ok">
			<div class="stat-label">Healthy</div>
			<div class="stat-value">{summary.healthy}</div>
			<div class="stat-hint">last run exit 0</div>
		</div>
		<div class="stat-card is-fail">
			<div class="stat-label">Failing</div>
			<div class="stat-value">{summary.failing}</div>
			<div class="stat-hint">last run non-zero</div>
		</div>
	</div>

	<section class="panel">
		<div class="panel__head">
			<h2 class="panel__title">One-time jobs</h2>
		</div>
		<div class="panel__body">
			<table class="table">
			<thead>
				<tr>
					<th>Script</th>
					<th>Status</th>
					<th>History</th>
					<th class="tar">Duration</th>
					<th class="tar">Last run</th>
				</tr>
			</thead>
			<tbody>
			{foreach $body['jobs'] as $k => $job}
				{if $job['active'] == 1 && strpos($job['name'], "onetime/") !== false}
				<tr>
					<td class="job-name"><a href="/{$job['name']}"><code>{$job['name']}</code></a></td>
					{if !empty($job['lastRun'])}
					<td>{if $job['lastRun']['exitCode'] == 0}<span class="badge badge-ok">OK</span>{else}<span class="badge badge-fail">Exit: {$job['lastRun']['exitCode']}</span>{/if}</td>
					<td>
						<div class="sparkline-box">
							<svg class="job-{$k}" width="240" height="28" stroke-width="2" {if $job['lastRun']['exitCode'] == 0}stroke="#16a34a" fill="rgba(22,163,74,0.12)"{else}stroke="#dc2626" fill="rgba(220,38,38,0.12)"{/if}></svg>
							<span class="spark-tooltip" hidden="true"></span>
						</div>
						<script>Sparkline.draw(document.querySelector(".job-{$k}"), {$job['history']|history}, sparkOptions)</script>
					</td>
					<td class="tar mono">{eval echo sprintf("%.3fs", $job['lastRun']['duration'])}</td>
					<td class="tar mono">{eval echo date("Y-m-d H:i:s", strtotime($job['lastRun']['stamp']))}</td>
					{else}
					<td colspan="4" class="empty-cell">No run history yet</td>
					{/if}
				</tr>
				{/if}
			{else}
				<tr><td colspan="5" class="empty-cell">No jobs found.</td></tr>
			{/foreach}
			</tbody>
			</table>
		</div>
	</section>

	<section class="panel">
		<div class="panel__head">
			<h2 class="panel__title">Scheduled jobs</h2>
		</div>
		<div class="panel__body">
			<table class="table">
			<thead>
				<tr>
					<th>Script</th>
					<th>Status</th>
					<th>History</th>
					<th class="tar">Duration</th>
					<th class="tar">Last run</th>
				</tr>
			</thead>
			<tbody>
			{foreach $body['jobs'] as $k => $job}
				{if $job['active'] == 1 && strpos($job['name'], "onetime/") === false}
				<tr>
					<td class="job-name"><a href="/{$job['name']}"><code>{$job['name']}</code></a></td>
					{if !empty($job['lastRun'])}
					<td>{if $job['lastRun']['exitCode'] == 0}<span class="badge badge-ok">OK</span>{else}<span class="badge badge-fail">Exit: {$job['lastRun']['exitCode']}</span>{/if}</td>
					<td>
						<div class="sparkline-box">
							<svg class="job-{$k}" width="240" height="28" stroke-width="2" {if $job['lastRun']['exitCode'] == 0}stroke="#16a34a" fill="rgba(22,163,74,0.12)"{else}stroke="#dc2626" fill="rgba(220,38,38,0.12)"{/if}></svg>
							<span class="spark-tooltip" hidden="true"></span>
						</div>
						<script>Sparkline.draw(document.querySelector(".job-{$k}"), {$job['history']|history}, sparkOptions)</script>
					</td>
					<td class="tar mono">{eval echo sprintf("%.3fs", $job['lastRun']['duration'])}</td>
					<td class="tar mono">{eval echo date("Y-m-d H:i:s", strtotime($job['lastRun']['stamp']))}</td>
					{else}
					<td colspan="4" class="empty-cell">No run history yet</td>
					{/if}
				</tr>
				{/if}
			{else}
				<tr><td colspan="5" class="empty-cell">No jobs found.</td></tr>
			{/foreach}
			</tbody>
			</table>
		</div>
	</section>

</div>

{load html_footer.tpl}
