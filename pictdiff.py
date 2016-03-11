#!/usr/bin/env python

TINT = 1 # exxagerate differences in tone
MINUTE = 5 # less than "x" points of difference
INCREASE_MINUTE = 2 # boost of minute differences

import sys
from PIL import Image

img1 = Image.open(sys.argv[1])
img2 = Image.open(sys.argv[2])

i1 = img1.load()
i2 = img2.load()

if img1.size != img2.size:
	print "Images %s and %s have different sizes, cannot compare" \
		% (sys.argv[1], sys.argv[2])
	sys.exit(1)

imgmap = Image.new( 'RGB', img1.size, "white")
imap = imgmap.load()

row_averages = []
for y in range(img1.size[1]):
	for x in range(img1.size[0]):
        	p1 = i1[x, y]
        	p2 = i2[x, y]
		diffpixel = [255, 255, 255]

		# color differences
		diffs = [p2[0] - p1[0], p2[1] - p1[1], p2[2] - p1[2]]
		absdiff = reduce(lambda a, b: abs(a) + abs(b), diffs)
		diffsmag = [a * TINT for a in diffs]
		diffplus = [max(0, a) for a in diffs]
		totplus = reduce(lambda a, b: a + b, diffplus)
		diffminus = [min(0, a) for a in diffs]

		# apply negative differences (e.g. less red -> take red)
		diffpixel = [ a + b for a, b in zip(diffpixel, diffminus)]
		# subtract positive differences (e.g. more red -> take from non-red channels)
		diffpixel = [ a - totplus for a in diffpixel ]
		# ... put back what we took from red
		diffpixel = [ a + b for a, b in zip(diffpixel, diffplus)]
		
		if absdiff > 0 and absdiff < MINUTE:
			# Increase contrast of minute differences
			diffpixel = [a - INCREASE_MINUTE for a in diffpixel]
		diffpixel = [max(0, a) for a in diffpixel]
		
		imap[x, y] = tuple(diffpixel)

imgmap.save(sys.argv[3])
