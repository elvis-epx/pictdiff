#!/usr/bin/env python

MINUTE = 5 # less than "x" points of difference
INCREASE_MINUTE = 2 # boost of minute differences

import sys
from PIL import Image

if len(sys.argv) < 4:
	print "Usage: %s oldpicture newpicture diffmap" % sys.argv[0]
	sys.exit(2)

img1 = Image.open(sys.argv[1]).convert('RGBA')
img2 = Image.open(sys.argv[2]).convert('RGBA')

i1 = img1.load()
i2 = img2.load()

if img1.size != img2.size:
	print "Images %s and %s have different sizes, cannot compare" \
		% (sys.argv[1], sys.argv[2])
	sys.exit(1)

imgmap = Image.new('RGB', img1.size, "white")
imap = imgmap.load()
totaldiff = 0

# Premultiply alpha, so color differences are mitigated by transparency
# Difference in alpha itself is accumulated separately in absdiff, so this
# will not mask any differenct
def pre_mult_alpha(old, new, channel):
	# simulate a black background
	old = old[channel] * (old[3] / 255.0)
	new = new[channel] * (new[3] / 255.0)
	return int(new) - int(old)

for y in range(img1.size[1]):
	for x in range(img1.size[0]):
        	p1 = i1[x, y]
        	p2 = i2[x, y]
		diffpixel = [255, 255, 255]

		# color differences, including alpha channel
		diffs = [
				pre_mult_alpha(p1, p2, 0), 
				pre_mult_alpha(p1, p2, 1),
				pre_mult_alpha(p1, p2, 2),
				p2[3] - p1[3]
			]
		absdiff = reduce(lambda a, b: abs(a) + abs(b), diffs)
		totaldiff += absdiff

		# these are for tinting, alpha left out of equation
		diffs = diffs[0:3]
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
print totaldiff
