#!/usr/bin/env python3

MINUTE = 5 # less than "x" points of difference
INCREASE_MINUTE = 2 # boost of minute differences

import sys
from PIL import Image

if len(sys.argv) < 4:
	print("Usage: %s oldpicture newpicture diffmap" % sys.argv[0], file=sys.stderr)
	sys.exit(2)

try:
	img1 = Image.open(sys.argv[1]).convert('RGBA')
except FileNotFoundError:
	print("First file could not be opened or it is not a picture.", file=sys.stderr)
	sys.exit(2)

try:
	img2 = Image.open(sys.argv[2]).convert('RGBA')
except FileNotFoundError:
	print("Second file could not be opened or it is not a picture.", file=sys.stderr)
	sys.exit(2)

i1 = img1.load()
i2 = img2.load()

if img1.size != img2.size:
	print("Pictures %s and %s have different sizes, cannot compare" \
		% (sys.argv[1], sys.argv[2]), file=sys.stderr)
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
		absdiff = abs(p2[3] - p1[3])
		diffs = [0, 0, 0]
		totplus = 0

		for i in range(0, 3):
			diffs[i] = pre_mult_alpha(p1, p2, i)
			absdiff += abs(diffs[i])
			diffpixel[i] += diffs[i]
			totplus += max(0, diffs[i])

		totaldiff += absdiff

		for i in range(0, 3):
			diffpixel[i] -= totplus
			if absdiff > 0 and absdiff < MINUTE:
				# Increase contrast of minute differences
				diffpixel[i] -= INCREASE_MINUTE
			diffpixel[i] = max(0, diffpixel[i])
		
		imap[x, y] = tuple(diffpixel)

try:
	imgmap.save(sys.argv[3])
except PermissionError:
	print("Could not write to diff picture file.", file=sys.stderr)
print(totaldiff)
