<?php

// @route POST /description

include("bootstrap.php");

header("Content-Type: application/json");

if (!isset($_POST["name"]) || !isset($_POST["description"])) {
	echo json_encode(array("ok" => false, "err" => "Missing name or description"));
	$db->close();
	exit;
}

$name = $_POST["name"];
$description = $_POST["description"];

$job = $db->get("SELECT name FROM jobs WHERE name = ?", $name);
if (!$job) {
	echo json_encode(array("ok" => false, "err" => "Job not found"));
	$db->close();
	exit;
}

$now = date("Y-m-d H:i:s");
$db->query("UPDATE jobs SET description = ?, updated_at = ? WHERE name = ?", $description, $now, $name);

echo json_encode(array("ok" => true, "name" => $name, "description" => $description));

$db->close();
