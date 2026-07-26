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

// Draw all sparklines once the DOM and the Sparkline library are ready. Data
// travels in each svg's data-history attribute so rendering never depends on
// inline script execution order.
function drawSparklines() {
	if (typeof Sparkline === 'undefined') { return; }
	var nodes = document.querySelectorAll('svg.sparkline[data-history]');
	for (var i = 0; i < nodes.length; i++) {
		var data;
		try { data = JSON.parse(nodes[i].getAttribute('data-history')); }
		catch (e) { data = []; }
		Sparkline.draw(nodes[i], data, sparkOptions);
	}
}

if (document.readyState === 'loading') {
	document.addEventListener('DOMContentLoaded', drawSparklines);
} else {
	drawSparklines();
}
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
					<th>Duration</th>
					<th>Last run</th>
				</tr>
			</thead>
			<tbody>
			{foreach $body['jobs'] as $k => $job}
				{if $job['active'] == 1 && strpos($job['name'], "onetime/") !== false}
				<tr>
					<td class="job-name">
						<div class="job-name__row">
							<a href="/{$job['name']}"><code>{$job['name']}</code></a>
							<button class="icon-btn" type="button" title="Edit description" data-name="{$job['name']|escape}" data-desc="{$job['description']|escape}" onclick="openDescModal(this)" aria-label="Edit description">
								<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"></path><path d="M16.5 3.5a2.1 2.1 0 0 1 3 3L7 19l-4 1 1-4Z"></path></svg>
							</button>
						</div>
						{if !empty($job['description'])}<div class="job-desc">{$job['description']|escape}</div>{/if}
					</td>
					{if !empty($job['lastRun'])}
					<td>{if $job['lastRun']['exitCode'] == 0}<span class="badge badge-ok">OK</span>{else}<span class="badge badge-fail">Exit: {$job['lastRun']['exitCode']}</span>{/if}</td>
					<td>
						<div class="sparkline-box">
							<svg class="sparkline" width="240" height="28" stroke-width="2" {if $job['lastRun']['exitCode'] == 0}stroke="#16a34a" fill="rgba(22,163,74,0.12)"{else}stroke="#dc2626" fill="rgba(220,38,38,0.12)"{/if} data-history='{$job['history']|history}'></svg>
							<span class="spark-tooltip" hidden="true"></span>
						</div>
					</td>
					<td class="mono">{eval echo sprintf("%.3fs", $job['lastRun']['duration'])}</td>
					<td class="mono">{eval echo date("Y-m-d H:i:s", strtotime($job['lastRun']['stamp']))}</td>
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
					<th>Duration</th>
					<th>Last run</th>
				</tr>
			</thead>
			<tbody>
			{foreach $body['jobs'] as $k => $job}
				{if $job['active'] == 1 && strpos($job['name'], "onetime/") === false}
				<tr>
					<td class="job-name">
						<div class="job-name__row">
							<a href="/{$job['name']}"><code>{$job['name']}</code></a>
							<button class="icon-btn" type="button" title="Edit description" data-name="{$job['name']|escape}" data-desc="{$job['description']|escape}" onclick="openDescModal(this)" aria-label="Edit description">
								<svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"></path><path d="M16.5 3.5a2.1 2.1 0 0 1 3 3L7 19l-4 1 1-4Z"></path></svg>
							</button>
						</div>
						{if !empty($job['description'])}<div class="job-desc">{$job['description']|escape}</div>{/if}
					</td>
					{if !empty($job['lastRun'])}
					<td>{if $job['lastRun']['exitCode'] == 0}<span class="badge badge-ok">OK</span>{else}<span class="badge badge-fail">Exit: {$job['lastRun']['exitCode']}</span>{/if}</td>
					<td>
						<div class="sparkline-box">
							<svg class="sparkline" width="240" height="28" stroke-width="2" {if $job['lastRun']['exitCode'] == 0}stroke="#16a34a" fill="rgba(22,163,74,0.12)"{else}stroke="#dc2626" fill="rgba(220,38,38,0.12)"{/if} data-history='{$job['history']|history}'></svg>
							<span class="spark-tooltip" hidden="true"></span>
						</div>
					</td>
					<td class="mono">{eval echo sprintf("%.3fs", $job['lastRun']['duration'])}</td>
					<td class="mono">{eval echo date("Y-m-d H:i:s", strtotime($job['lastRun']['stamp']))}</td>
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

<div id="descModal" class="modal-overlay" hidden>
	<div class="modal">
		<div class="modal__head">
			<span>Edit description</span>
			<button class="icon-btn" type="button" onclick="closeDescModal()" aria-label="Close">&times;</button>
		</div>
		<div class="modal__body">
			<div class="modal__label" id="descJobName"></div>
			<textarea id="descText" class="form-control" rows="4" placeholder="Describe what this job does…"></textarea>
			<div id="descError" class="alert alert-danger" hidden></div>
		</div>
		<div class="modal__foot">
			<button class="btn btn-secondary" type="button" onclick="closeDescModal()">Cancel</button>
			<button class="btn btn-primary" type="button" onclick="saveDesc()">Save</button>
		</div>
	</div>
</div>

<script>
var descName = null;

function openDescModal(btn) {
	descName = btn.getAttribute('data-name');
	document.getElementById('descJobName').textContent = descName;
	document.getElementById('descText').value = btn.getAttribute('data-desc') || '';
	var err = document.getElementById('descError');
	err.hidden = true;
	err.textContent = '';
	document.getElementById('descModal').hidden = false;
	document.getElementById('descText').focus();
}

function closeDescModal() {
	document.getElementById('descModal').hidden = true;
}

function showDescError(msg) {
	var err = document.getElementById('descError');
	err.textContent = msg;
	err.hidden = false;
}

function saveDesc() {
	if (!descName) { return; }
	var body = new FormData();
	body.append('name', descName);
	body.append('description', document.getElementById('descText').value);
	fetch('/description', { method: 'POST', body: body })
		.then(function (r) { return r.json(); })
		.then(function (data) {
			if (data && data.ok) { location.reload(); }
			else { showDescError((data && data.err) || 'Save failed'); }
		})
		.catch(function () { showDescError('Network error'); });
}

document.addEventListener('keydown', function (event) {
	if (event.key === 'Escape') { closeDescModal(); }
});
</script>

{load html_footer.tpl}
