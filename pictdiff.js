#!/usr/bin/env node

var j = require("jimp");
var fs = require("fs");

function start()
{
	j.read(process.argv[2], (err, img1) => {
		if (err) {
			console.error("First file is not a picture or is not readable.");
			process.exit(1);
			return;
		}
		j.read(process.argv[3], (err, img2) => {
			if (err) {
				console.error("Second file is not a picture or is not readable.");
				process.exit(1);
				return;
			}
			if (img1.bitmap.width !== img2.bitmap.width ||
					img1.bitmap.height !== img2.bitmap.height) {
				console.error("Pictures to be compared must have the same size");
				process.exit(1);
				return;
			}
			var imgmap = new j(img1.bitmap.width, img1.bitmap.height, 0x306090FF,
			(err, image) => {
				compare(img1, img2, image);
			});
		});
	});
}

function mult_alpha(old, neu, channel)
{
	old = old[channel] * old[3] / 255.0
	neu = neu[channel] * neu[3] / 255.0
	return Math.floor(neu) - Math.floor(old)
}

function compare(img1, img2, imgmap)
{
	var totaldiff = 0;

	imgmap.scan(0, 0, img1.bitmap.width, img1.bitmap.height, (x, y, idx) => {
		var p1 = [
			img1.bitmap.data[ idx + 0 ],
			img1.bitmap.data[ idx + 1 ],
			img1.bitmap.data[ idx + 2 ],
			img1.bitmap.data[ idx + 3 ]
		];
		var p2 = [
			img2.bitmap.data[ idx + 0 ],
			img2.bitmap.data[ idx + 1 ],
			img2.bitmap.data[ idx + 2 ],
			img2.bitmap.data[ idx + 3 ]
		];

		var i;
		var diffpixel = [255, 255, 255];
		var absdiff = Math.abs(p2[3] - p1[3]);
		var diffs = [0, 0, 0];
		var totplus = 0;

		for (i = 0; i < 3; ++i) {
			diffs[i] = mult_alpha(p1, p2, i);
			absdiff += Math.abs(diffs[i]);
			totplus += Math.max(0, diffs[i]);
			diffpixel[i] += diffs[i];
		}

		for (i = 0; i < 3; ++i) {
			diffpixel[i] -= totplus;
			if (absdiff > 0 && absdiff < 5) {
				diffpixel[i] -= 2;
			}
			diffpixel[i] = Math.max(0, diffpixel[i]);
			imgmap.bitmap.data[ idx + i ] = diffpixel[i];
		}

		totaldiff += absdiff;
	});

	imgmap.write(process.argv[4], (err) => {
		if (err) {
			console.error("Could not write diff map.");
		}
		console.log(totaldiff);
	});
}

function usage()
{
	console.error();
	console.error("Usage: node pictdiff.js <picture A> <picture B> <diff map>");
	console.error();
	console.error("Example: node pictdiff.js a.png b.png diff.png");
	console.error();
}

if (process.argv.length < 5) {
	usage();
	return;
}
start();
