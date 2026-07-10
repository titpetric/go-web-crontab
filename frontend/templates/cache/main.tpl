<?php $_v=&$this->vars; $this->push();$this->load("html_header.tpl");$this->assign($_v);$this->render();$this->pop();?>


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
			<div class="stat-value"><?php echo $_v['summary']['total'];?>
</div>
			<div class="stat-hint">tracked in crontab</div>
		</div>
		<div class="stat-card">
			<div class="stat-label">Scheduled</div>
			<div class="stat-value"><?php echo $_v['summary']['scheduled'];?>
</div>
			<div class="stat-hint"><?php echo $_v['summary']['onetime'];?>
 one-time</div>
		</div>
		<div class="stat-card is-ok">
			<div class="stat-label">Healthy</div>
			<div class="stat-value"><?php echo $_v['summary']['healthy'];?>
</div>
			<div class="stat-hint">last run exit 0</div>
		</div>
		<div class="stat-card is-fail">
			<div class="stat-label">Failing</div>
			<div class="stat-value"><?php echo $_v['summary']['failing'];?>
</div>
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
			<?php if(!empty($_v['body']['jobs']))foreach($_v['body']['jobs'] as $_v['k'] => $_v['job']){?>
				<?php if($_v['job']['active'] == 1 && strpos($_v['job']['name'], "onetime/") !== false){?>

				<tr>
					<td class="job-name"><a href="/<?php echo $_v['job']['name'];?>
"><code><?php echo $_v['job']['name'];?>
</code></a></td>
					<?php if(!empty($_v['job']['lastRun'])){?>

					<td><?php if($_v['job']['lastRun']['exitCode'] == 0){?>
<span class="badge badge-ok">OK</span><?php }else{?>
<span class="badge badge-fail">Exit: <?php echo $_v['job']['lastRun']['exitCode'];?>
</span><?php }?></td>
					<td>
						<div class="sparkline-box">
							<svg class="job-<?php echo $_v['k'];?>
" width="240" height="28" stroke-width="2" <?php if($_v['job']['lastRun']['exitCode'] == 0){?>
stroke="#16a34a" fill="rgba(22,163,74,0.12)"<?php }else{?>
stroke="#dc2626" fill="rgba(220,38,38,0.12)"<?php }?>></svg>
							<span class="spark-tooltip" hidden="true"></span>
						</div>
						<script>Sparkline.draw(document.querySelector(".job-<?php echo $_v['k'];?>
"), <?php echo history($_v['job']['history']);?>
, sparkOptions)</script>
					</td>
					<td class="tar mono"><?php echo sprintf("%.3fs", $_v['job']['lastRun']['duration']);?>
</td>
					<td class="tar mono"><?php echo date("Y-m-d H:i:s", strtotime($_v['job']['lastRun']['stamp']));?>
</td>
					<?php }else{?>

					<td colspan="4" class="empty-cell">No run history yet</td>
					<?php }?>
				</tr>
				<?php }?>
			<?php }else{?>

				<tr><td colspan="5" class="empty-cell">No jobs found.</td></tr>
			<?php }?>
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
			<?php if(!empty($_v['body']['jobs']))foreach($_v['body']['jobs'] as $_v['k'] => $_v['job']){?>
				<?php if($_v['job']['active'] == 1 && strpos($_v['job']['name'], "onetime/") === false){?>

				<tr>
					<td class="job-name"><a href="/<?php echo $_v['job']['name'];?>
"><code><?php echo $_v['job']['name'];?>
</code></a></td>
					<?php if(!empty($_v['job']['lastRun'])){?>

					<td><?php if($_v['job']['lastRun']['exitCode'] == 0){?>
<span class="badge badge-ok">OK</span><?php }else{?>
<span class="badge badge-fail">Exit: <?php echo $_v['job']['lastRun']['exitCode'];?>
</span><?php }?></td>
					<td>
						<div class="sparkline-box">
							<svg class="job-<?php echo $_v['k'];?>
" width="240" height="28" stroke-width="2" <?php if($_v['job']['lastRun']['exitCode'] == 0){?>
stroke="#16a34a" fill="rgba(22,163,74,0.12)"<?php }else{?>
stroke="#dc2626" fill="rgba(220,38,38,0.12)"<?php }?>></svg>
							<span class="spark-tooltip" hidden="true"></span>
						</div>
						<script>Sparkline.draw(document.querySelector(".job-<?php echo $_v['k'];?>
"), <?php echo history($_v['job']['history']);?>
, sparkOptions)</script>
					</td>
					<td class="tar mono"><?php echo sprintf("%.3fs", $_v['job']['lastRun']['duration']);?>
</td>
					<td class="tar mono"><?php echo date("Y-m-d H:i:s", strtotime($_v['job']['lastRun']['stamp']));?>
</td>
					<?php }else{?>

					<td colspan="4" class="empty-cell">No run history yet</td>
					<?php }?>
				</tr>
				<?php }?>
			<?php }else{?>

				<tr><td colspan="5" class="empty-cell">No jobs found.</td></tr>
			<?php }?>
			</tbody>
			</table>
		</div>
	</section>

</div>

<?php $this->push();$this->load("html_footer.tpl");$this->assign($_v);$this->render();$this->pop();?>

