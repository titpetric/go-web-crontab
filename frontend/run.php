<?php

// @route GET /run/{jobName}
// @route GET /run/{group}/{jobName}

include("bootstrap.php");

$jobName = route_job_name();
$job = $db->get("SELECT * FROM jobs WHERE name = ?", $jobName);

$title = $jobName;
if ($job) {
	$output = sprintf("Running %s in the background, if not already running", $jobName);
} else {
	$output = sprintf("Job %s was not found", $jobName);
}

$tpl->load("job_detail.tpl");
$tpl->assign(array(
	"title" => $title,
	"output" => $output,
	"job" => $job,
));
$tpl->render();

$db->close();
