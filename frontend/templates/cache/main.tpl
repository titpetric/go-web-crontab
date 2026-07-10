<?php $_v=&$this->vars; $this->push();$this->load("html_header.tpl");$this->assign($_v);$this->render();$this->pop();?>


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

		<?php if(!empty($_v['body']['jobs']))foreach($_v['body']['jobs'] as $_v['k'] => $_v['job']){?>
			<?php if($_v['job']['active'] == 1 && strpos($_v['job']['name'], "onetime/") !== false){?>

	<tr>
		<td><a href="/<?php echo $_v['job']['name'];?>
"><?php echo $_v['job']['name'];?>
</a></td>

				<?php if(!empty($_v['job']['lastRun'])){?>

			<td>

      <div style="position: relative">
        <svg class="job-<?php echo $_v['k'];?>
" width="300" height="20" stroke-width="2" <?php if($_v['job']['lastRun']['exitCode'] == 0){?>
stroke="blue" fill="rgba(0, 0, 255, .2)"<?php }else{?>
stroke="red" fill="rgba(255, 0, 0, .2)"<?php }?>></svg>
        <span class="tooltip" hidden="true"></span>
      </div>

				<script>
					sparkline.sparkline(document.querySelector(".job-<?php echo $_v['k'];?>
"), <?php echo history($_v['job']['history']);?>
, options)
				</script>
			</td>
			<td class="tar"><?php echo sprintf("%.3fs", $_v['job']['lastRun']['duration']);?>
</td>
			<td class="tar"><?php echo date("Y-m-d H:i:s", strtotime($_v['job']['lastRun']['stamp']));?>
</td>
				<?php }else{?>

			<td colspan="4"><i>No last run info</i></td>
				<?php }?>
	</tr>
			<?php }?>
		<?php }?>

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

		<?php if(!empty($_v['body']['jobs']))foreach($_v['body']['jobs'] as $_v['k'] => $_v['job']){?>
			<?php if($_v['job']['active'] == 1 && strpos($_v['job']['name'], "onetime/") === false){?>

	<tr>
		<td><a href="/<?php echo $_v['job']['name'];?>
"><?php echo $_v['job']['name'];?>
</a></td>

				<?php if(!empty($_v['job']['lastRun'])){?>

			<td>

      <div style="position: relative">
        <svg class="job-<?php echo $_v['k'];?>
" width="300" height="20" stroke-width="2" <?php if($_v['job']['lastRun']['exitCode'] == 0){?>
stroke="blue" fill="rgba(0, 0, 255, .2)"<?php }else{?>
stroke="red" fill="rgba(255, 0, 0, .2)"<?php }?>></svg>
        <span class="tooltip" hidden="true"></span>
      </div>

				<script>
					sparkline.sparkline(document.querySelector(".job-<?php echo $_v['k'];?>
"), <?php echo history($_v['job']['history']);?>
, options)
				</script>
			</td>
			<td class="tar"><?php echo sprintf("%.3fs", $_v['job']['lastRun']['duration']);?>
</td>
			<td class="tar"><?php echo date("Y-m-d H:i:s", strtotime($_v['job']['lastRun']['stamp']));?>
</td>
				<?php }else{?>

			<td colspan="4"><i>No last run info</i></td>
				<?php }?>
	</tr>
			<?php }?>
		<?php }?>

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

<?php $this->push();$this->load("html_footer.tpl");$this->assign($_v);$this->render();$this->pop();?>

