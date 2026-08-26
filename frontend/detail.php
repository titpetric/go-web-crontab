<?php

// @route GET /{jobName}/{ID}
// @route GET /{group}/{jobName}/{ID}

include("bootstrap.php");

$jobName = route_job_name();
$ID = (int)$_PATH["ID"];

$row = $db->get("SELECT output FROM (SELECT ROW_NUMBER() OVER (ORDER BY stamp DESC) AS ID, output FROM logs WHERE name = ?) WHERE ID = ?", $jobName, $ID);

$output = "";
if ($row) {
	$output = $row["output"];

	// A job that logged structured output gets the per-line rendering, and one
	// that logged plain text is shown as the text it is. The shape is checked
	// before decoding rather than after: json_decode raises on input that is
	// not json, so calling it on a plain log line would fail the request
	// instead of returning the null the fallback is written for.
	$head = substr(ltrim($output), 0, 1);
	if ($head == "[" || $head == "{") {
		$decoded = json_decode($output);
		if ($decoded) {
			$output = $decoded;
		}
	}
}

$job = array("jobName" => $jobName, "ID" => $ID);
$title = $jobName . ", log " . $ID;

$tpl->load("job_detail.tpl");
$tpl->assign(array(
	"title" => $title,
	"output" => $output,
	"job" => $job,
));
$tpl->render();

$db->close();
