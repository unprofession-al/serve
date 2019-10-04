package main

type htmlData struct {
	Body  string
	Title string
}

var htmlScaffold = `<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>{{ .Title }}</title>
	<link href="https://fonts.googleapis.com/css?family=Merriweather&display=swap" rel="stylesheet">
	<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.4.1/jquery.min.js"></script>

	<style>
		html, body {
			font-family: 'Merriweather', serif;
			line-height: 1.6;
			margin: 0px;
			padding: 0px;
		}

		.container {
			max-width: 900px;
			padding: 60px;
			margin: 0 auto;
		}

		table {
		    font-size: 0.8em;
		    margin: 25px auto;
		    border-collapse: collapse;
		    border: 1px solid #eee;
		    border-bottom: 2px solid #c5e8e2;
		    box-shadow: 0px 0px 20px rgba(0, 0, 0, 0.1), 0px 10px 20px rgba(0, 0, 0, 0.05), 0px 20px 20px rgba(0, 0, 0, 0.05), 0px 30px 20px rgba(0, 0, 0, 0.05);
		}

		table tr:hover {
		    background: #f4f4f4;
		}

		table th, table td {
		    border: 1px solid #eee;
		    padding: 10px;
		    border-collapse: collapse;
		}

		table th {
			background-color: #c5e8e2;
		    text-transform: uppercase;
		}

		table td {
		    text-align: left;
		}

		table th.last {
		    border-right: none;
		}

		a {
			color: #50a395;
			text-decoration: none;
		}

		code {
			background-color: #c5e8e2;
			border-radius: 4px;
			padding-right: 4px;
			padding-left: 4px;
		}

		pre {
			background-color: #c5e8e2;
			border-radius: 4px;
			padding-left: 6px;
			padding-right: 6px;
            overflow-x: auto;
            white-space: pre-wrap;
            white-space: -moz-pre-wrap;
            white-space: -pre-wrap;
            white-space: -o-pre-wrap;
            word-wrap: break-word;
		}

		img {
			width: 100%;
		    box-shadow: 0px 0px 20px rgba(0, 0, 0, 0.1), 0px 10px 20px rgba(0, 0, 0, 0.05), 0px 20px 20px rgba(0, 0, 0, 0.05), 0px 30px 20px rgba(0, 0, 0, 0.05);
		}

		 pre>code {
			padding-left: 0;
			padding-right: 0;
		}

		#drawler {
			background-image: linear-gradient(to right, rgba(0,0,0,0.2), rgba(0,0,0,0));
			position: fixed;
			height: 100%;
			z-index: 20;
			width: 30px;
			transition: background-color .07s ease;
			font-size: 18px;
		}

		#drawler:hover {
			background-image: none;
			background-color: #FFF;
			width: 400px;
			box-shadow: 0 5px 10px 5px rgba(0,0,0,0.25);
		}

		#drawler:hover>.drawler_indicator {
			display: none;
		}

		#drawler>.drawler_content {
			position: fixed;
			width: 400px;
			padding: 40px;
			padding-top: 60px;
			display: block;
			color: black;
			font-size: .8em;
			margin: 20px;
			left: -400px;
			transition: left .07s linear;
		}

		#drawler:hover>.drawler_content {
			left: 0px;
			width: 300px;
		}

		#drawler:not(:hover)>.drawler_content {
			width: 0px;
		}

		#drawler>.drawler_content>.toc>nav>ul {
			list-style-type: none;
			padding-left: 0;
		}
	</style>
</head>
<body>
	<div id="drawler">
        <div class="drawler_content">
			<div class="toc"></div>
        </div>
    </div>

	<div class="container">
	{{ .Body }}
	</div>
</body>

<script>
	var ToC =
	  "<nav role='navigation' class='table-of-contents'>" +
		"<div><b>Table of Contents</b></div>" +
		"<ul>";

	var newLine, el, title, link;

	$(".container h2").each(function() {

	  el = $(this);
	  title = el.text();
	  link = "#" + el.attr("id");

	  newLine =
		"<li>" +
		  "<a href='" + link + "'>" +
			title +
		  "</a>" +
		"</li>";

	  ToC += newLine;

	});

	ToC +=
	   "</ul>" +
	  "</nav>";

	$(".toc").prepend(ToC);
</script>
</html>`
