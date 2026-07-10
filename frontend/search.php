<?php

// @route POST /search

include("bootstrap.php");

$pageNumber = 0;
if (isset($_POST["pageNumber"])) {
	$pageNumber = $_POST["pageNumber"] + 0;
}
$pageSize = 25;
$offset = $pageNumber * $pageSize;

if (!isset($_POST["jobName"]) || !isset($_POST["searchStr"])) {
	header("Content-Type: application/json");
	echo json_encode(array("err" => "Not correct paramaters"));
	$db->close();
	exit;
}

$jobName = $_POST["jobName"];
$searchStr = "%" . $_POST["searchStr"] . "%";
$link = "/" . $jobName;

$rows = $db->get_all("SELECT ID, stamp, duration, exitCode FROM (SELECT ROW_NUMBER() OVER (ORDER BY stamp DESC) AS ID, stamp, duration / 1000000000.0 AS duration, exit_code AS exitCode, output FROM logs WHERE name = ? AND output LIKE ?) ORDER BY ID LIMIT ? OFFSET ?", $jobName, $searchStr, $pageSize, $offset);

$logs = array();
foreach ($rows as $row) {
	$logs[] = array(
		"ID" => $row["ID"],
		"link" => $link,
		"exitCode" => $row["exitCode"],
		"duration" => $row["duration"],
		"stamp" => $row["stamp"],
	);
}

header("Content-Type: application/json");
echo json_encode(array("logs" => $logs));

$db->close();
