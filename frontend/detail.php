<?php

// @route GET /{jobName}/{ID}
// @route GET /{group}/{jobName}/{ID}

include("bootstrap.php");

$jobName = route_job_name();
$ID = (int)$_PATH["ID"];

$row = $db->get("SELECT output FROM (SELECT ROW_NUMBER() OVER (ORDER BY stamp DESC) AS ID, output FROM logs WHERE name = ?) WHERE ID = ?", $jobName, $ID);

$output = "";
if ($row) {
	$output = json_decode($row["output"]);
	if (!$output) {
		$output = $row["output"];
	}
}

$job = array("jobName" => $jobName, "ID" => $ID);
$title = $jobName . ", log " . $ID;

$tpl->load("job_detail.tpl");
$tpl->assign(array("title" => $title, "output" => $output, "job" => $job));
$tpl->render();

$db->close();
