This collection of programs implement image comparison. Each program
prints a metric of difference, and generates an image for visual 
inspection of differences.

How to use:

./pictdiff.py img1.png img2.png diff.png

The diff image is white, with dark spots where differences are found.
The spots are darker as the differences are more intense. The spots
are also tinted to reflect the color changes. For example, if "img1"
is generally more red than "img2", then "diff" will be cyan. On the
other hand, if "img2" is more red, then "diff" is also red.
